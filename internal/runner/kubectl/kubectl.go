package k8s

import (
	"context"
	"errors"
	"fmt"
	"os"
	"runtime/debug"
	"sync"
	"time"

	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/busyster996/dagflow/internal/common"
	"github.com/busyster996/dagflow/internal/runner/kubectl/types"
	"github.com/busyster996/dagflow/internal/storage"
	"github.com/busyster996/dagflow/pkg/logx"
)

type sKubectl struct {
	kubeConf      *rest.Config
	client        *kubernetes.Clientset
	dynamicClient *dynamic.DynamicClient
	storage       storage.IStep
	subCommand    string
	workspace     string

	Config              string             `json:"kube_config"`
	Namespace           string             `json:"namespace"`
	ImageTag            string             `json:"image_tag"`
	IgnoreInitContainer *bool              `json:"ignoreInitContainer"`
	Resources           []*types.SResource `json:"resources"`
}

func (k *sKubectl) init() (err error) {
	defer func() {
		r := recover()
		if r != nil {
			logx.Errorln(string(debug.Stack()), r)
			err = fmt.Errorf("panic: %s", r)
		}
	}()
	content, err := k.storage.Content()
	if err != nil {
		return err
	}
	if err = yaml.Unmarshal([]byte(content), k); err != nil {
		return err
	}

	if k.Config == "" {
		k.Config, err = k.storage.Env().Get("KUBECONFIG")
		if err != nil {
			return err
		}
	}

	file, err := os.ReadFile(k.Config)
	if err != nil {
		file = []byte(k.Config)
	}
	k.kubeConf, err = clientcmd.RESTConfigFromKubeConfig(file)
	if err != nil {
		return err
	}
	k.kubeConf.Burst = 1000
	k.kubeConf.QPS = 500
	k.kubeConf.Timeout = 5 * time.Minute

	k.client, err = kubernetes.NewForConfig(k.kubeConf)
	if err != nil {
		return err
	}

	k.dynamicClient, err = dynamic.NewForConfig(k.kubeConf)
	if err != nil {
		return err
	}
	for kk, res := range k.Resources {
		if res.Name == "" {
			return errors.New("name is empty")
		}
		k.Resources[kk].Namespace, err = k.getResourceValue(res.Namespace, k.Namespace, "NAMESPACE")
		if err != nil {
			return err
		}
		k.Resources[kk].ImageTag, err = k.getResourceValue(res.ImageTag, k.ImageTag, "IMAGE_TAG")
		if err != nil {
			return err
		}
		k.Resources[kk].Kind, err = k.getResourceValue(res.Kind, types.Deployment, "KIND")
		if err != nil {
			return err
		}
		if res.IgnoreInitContainer == nil {
			k.Resources[kk].IgnoreInitContainer = k.IgnoreInitContainer
		}
	}
	return
}

func (k *sKubectl) Run(ctx context.Context) (exit common.ExecCode, err error) {
	defer func() {
		r := recover()
		if r != nil {
			logx.Errorln(string(debug.Stack()), r)
			err = fmt.Errorf("panic: %s", r)
			exit = common.ExecCodeSystemErr
		}
	}()
	timeout, err := k.storage.Timeout()
	if err != nil {
		return common.ExecCodeSystemErr, err
	}
	if timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, timeout)
		defer cancel()
	}

	if err = k.init(); err != nil {
		return common.ExecCodeSystemErr, err
	}
	var wg sync.WaitGroup
	var errCh = make(chan error, len(k.Resources))
	var done = make(chan struct{})
	go func() {
		defer close(done)
		for {
			select {
			case _err, ok := <-errCh:
				if !ok {
					return
				}
				err = errors.Join(err, _err)
			case <-ctx.Done():
				if ctx.Err() != nil {
					k.storage.Log().Write(ctx.Err().Error())
					err = ctx.Err()
				}
				return
			}
		}
	}()
	for _, resource := range k.Resources {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if _err := k.run(ctx, resource); _err != nil {
				k.storage.Log().Writef("%s/%s Error: %s", resource.GetNamespace(), resource.GetName(), _err.Error())
				errCh <- _err
			}
		}()
	}
	wg.Wait()
	close(errCh)
	<-done
	if err != nil {
		return common.ExecCodeSystemErr, err
	}
	return
}

func (k *sKubectl) getResourceValue(lValue, gValue, env string) (string, error) {
	if lValue != "" {
		return lValue, nil
	}
	if gValue != "" {
		return gValue, nil
	}
	value, err := k.storage.Env().Get(env)
	if err != nil {
		return "", err
	}
	if value == "" {
		value, err = k.storage.GlobalEnv().Get(env)
	}
	return value, err
}

func (k *sKubectl) run(ctx context.Context, resource *types.SResource) (err error) {
	var rs = ResourceFor(ctx, k.storage, k.client, resource)
	switch k.subCommand {
	case "restart":
		if err = rs.Restart(); err != nil {
			return err
		}
	case "update":
		if err = rs.Update(); err != nil {
			return err
		}
	case "scale":
		if err = rs.Scale(resource.GetReplicas()); err != nil {
			return err
		}
	case "status":
	default:
		return fmt.Errorf("unknown command: %s", k.subCommand)
	}
	if err = Status(ctx, k.storage, k.dynamicClient, resource); err != nil {
		return err
	}
	return rs.Println()
}

func (k *sKubectl) Clear() error {
	return nil
}
