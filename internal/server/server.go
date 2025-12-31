package server

import (
	"context"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/gin-gonic/gin"
	"github.com/pires/go-proxyproto"
	"github.com/pkg/errors"
	"github.com/segmentio/ksuid"
	"github.com/spf13/viper"
	"github.com/tus/tusd/v2/pkg/filelocker"
	"github.com/tus/tusd/v2/pkg/filestore"
	tusd "github.com/tus/tusd/v2/pkg/handler"

	"github.com/busyster996/dagflow/internal/server/router"
	"github.com/busyster996/dagflow/internal/server/tus/redislocker"
	"github.com/busyster996/dagflow/internal/server/tus/redisstore"
	"github.com/busyster996/dagflow/internal/utility"
	"github.com/busyster996/dagflow/pkg/listeners"
	"github.com/busyster996/dagflow/pkg/logx"
	"github.com/busyster996/dagflow/pkg/sockets"
)

var server *sServer

type sServer struct {
	ctx       context.Context
	cancel    context.CancelFunc
	listeners []net.Listener
	http      *http.Server
	tusd      *tusd.Handler
	wg        *sync.WaitGroup
}

func Close() error {
	if server == nil {
		return errors.New("server is not running")
	}
	logx.Infoln("stop service")
	server.close()
	server.wg.Wait()

	server.cancel()
	logx.Infoln("service stopped")
	time.Sleep(300 * time.Millisecond)
	return nil
}

func Start(ctx context.Context) error {
	if server != nil {
		return errors.New("server is running")
	}
	server = &sServer{
		wg: new(sync.WaitGroup),
	}
	server.ctx, server.cancel = context.WithCancel(ctx)
	return server.startServer()
}

func (p *sServer) startServer() error {
	gin.DisableConsoleColor()
	gin.SetMode(gin.ReleaseMode)
	handler, err := router.New()
	if err != nil {
		logx.Errorln(err)
		return err
	}

	if err = server.newTusSvr("/api/v1/files/"); err != nil {
		logx.Errorln(err)
		return err
	}
	handler.Any("/api/v1/files", gin.WrapH(http.StripPrefix("/api/v1/files", p.tusd)))
	handler.Any("/api/v1/files/*any", gin.WrapH(http.StripPrefix("/api/v1/files/", p.tusd)))

	p.http = &http.Server{
		Handler:           http.StripPrefix(viper.GetString("relative_path"), handler),
		ReadHeaderTimeout: 120 * time.Second,
		IdleTimeout:       120 * time.Second,
		MaxHeaderBytes:    15 << 20, // 15MB
		BaseContext: func(_ net.Listener) context.Context {
			return p.ctx
		},
	}

	_ = retry.Do(
		func() error {
			if err = p.loadListeners([]string{
				viper.GetString("addr"),
				sockets.DefaultPipePath(utility.ServiceName),
			}); err != nil {
				logx.Errorln(err)
				return err
			}
			return nil
		},
		retry.Attempts(0),
		retry.DelayType(func(n uint, err error, config *retry.Config) time.Duration {
			_max := time.Duration(n)
			if _max > 8 {
				_max = 8
			}
			duration := time.Second * _max * _max
			return duration
		}),
	)
	for _, ln := range p.listeners {
		p.wg.Add(1)
		go func(ln net.Listener) {
			defer p.wg.Done()
			if err := p.http.Serve(ln); err != nil && !errors.Is(err, http.ErrServerClosed) {
				logx.Errorln(err)
			}
		}(ln)
	}

	return nil
}

func (p *sServer) loadListeners(hosts []string) error {
	for _, ln := range p.listeners {
		_ = ln.Close()
	}
	p.listeners = []net.Listener{}
	for _, host := range hosts {
		proto, addr, ok := strings.Cut(host, "://")
		if !ok {
			logx.Warnf("bad format %s, expected PROTO://ADDR", host)
			proto = "tcp"
			addr = host
		}
		ln, err := listeners.New(p.ctx, proto, addr, nil)
		if err != nil {
			logx.Errorln(err)
			return err
		}
		logx.Infof("Listener created on %s (%s)", proto, addr)
		p.listeners = append(p.listeners, &proxyproto.Listener{Listener: ln})
	}
	return nil
}

func (p *sServer) close() {
	logx.Info("shutdown server")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	// close http server
	if p.http != nil {
		_ = p.http.Shutdown(ctx)
	}

	// close net listener
	if p.listeners != nil {
		for _, ln := range p.listeners {
			_ = ln.Close()
		}
	}
}

func (p *sServer) newTusSvr(basePath string) error {
	tmpdir := filepath.Join(viper.GetString("workspace_dir"), ".tusd")
	_ = os.MkdirAll(tmpdir, os.ModePerm)
	composer := tusd.NewStoreComposer()
	if viper.GetString("redis_uri") != "" {
		_store, err := redisstore.New(tmpdir, viper.GetString("redis_uri"))
		if err != nil {
			logx.Errorln(err)
			return err
		}
		_locker, err := redislocker.New(viper.GetString("redis_uri"))
		if err != nil {
			logx.Errorln(err)
			return err
		}
		_store.UseIn(composer)
		_locker.UseIn(composer)
	} else {
		_store := filestore.New(tmpdir)
		_locker := filelocker.New(tmpdir)
		_store.UseIn(composer)
		_locker.UseIn(composer)
	}
	var err error
	p.tusd, err = tusd.NewHandler(tusd.Config{
		BasePath:      basePath,
		StoreComposer: composer,
		Cors: &tusd.CorsConfig{
			Disable: true,
		},
		PreUploadCreateCallback: func(hook tusd.HookEvent) (tusd.HTTPResponse, tusd.FileInfoChanges, error) {
			id := ksuid.New().String()
			taskID, ok := hook.Upload.MetaData["task_id"]
			if !ok {
				return tusd.HTTPResponse{
					StatusCode: http.StatusBadRequest,
					Body:       "task_id is required",
				}, tusd.FileInfoChanges{}, errors.New("task_id is required")
			}
			return tusd.HTTPResponse{}, tusd.FileInfoChanges{
				ID: filepath.Join(taskID, id),
			}, nil
		},
		PreFinishResponseCallback: func(hook tusd.HookEvent) (tusd.HTTPResponse, error) {
			if hook.Upload.IsFinal {
				filename := hook.Upload.MetaData["filename"]
				if filename == "" {
					filename = filepath.Base(hook.Upload.ID)
				}

				src := filepath.Join(tmpdir, hook.Upload.ID)
				dst := filepath.Join(viper.GetString("workspace_dir"), filepath.Dir(hook.Upload.ID), filename)
				if err = utility.CopyFile(src, dst); err != nil {
					return tusd.HTTPResponse{
						StatusCode: http.StatusInternalServerError,
						Body:       "failed to copy file",
					}, err
				}
			}
			return tusd.HTTPResponse{
				Header: map[string]string{
					"ID":   filepath.Base(hook.Upload.ID),
					"Path": filepath.Dir(hook.Upload.ID),
				},
			}, nil
		},
	})
	if err != nil {
		logx.Errorln(err)
		return err
	}
	return nil
}
