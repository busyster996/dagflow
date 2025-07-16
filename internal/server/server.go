package server

import (
	"context"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/gin-gonic/gin"
	"github.com/pires/go-proxyproto"
	"github.com/pkg/errors"
	"gorm.io/gorm"

	"github.com/busyster996/dagflow/internal/server/router"
	"github.com/busyster996/dagflow/internal/utils"
	"github.com/busyster996/dagflow/pkg/listeners"
	"github.com/busyster996/dagflow/pkg/logx"
)

var server *sServer

type sServer struct {
	ctx       context.Context
	cancel    context.CancelFunc
	listeners []net.Listener
	http      *http.Server
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

func Start(ctx context.Context, db *gorm.DB, addr, relativePath, workspace string) error {
	if server != nil {
		return errors.New("server is running")
	}
	server = &sServer{
		wg: new(sync.WaitGroup),
	}
	server.ctx, server.cancel = context.WithCancel(ctx)
	return server.startServer(db, addr, workspace, relativePath)
}

func (p *sServer) startServer(db *gorm.DB, addr, workspace string, relativePath string) error {
	gin.DisableConsoleColor()
	gin.SetMode(gin.ReleaseMode)
	handler, err := router.New(db, workspace)
	if err != nil {
		logx.Errorln(err)
		return err
	}

	p.http = &http.Server{
		Handler:           http.StripPrefix(relativePath, handler),
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
				addr,
				utils.PipeName(),
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
		ln, err := listeners.Init(proto, addr, nil)
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
