package ds

import (
	"context"
	"fmt"
	"time"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubetypes "k8s.io/apimachinery/pkg/types"
	appv1 "k8s.io/client-go/kubernetes/typed/apps/v1"
	kuberetry "k8s.io/client-go/util/retry"

	"github.com/busyster996/dagflow/internal/storage"
	"github.com/busyster996/dagflow/internal/worker/runner/kubectl/types"
)

type SDaemonSet struct {
	context.Context
	Client appv1.DaemonSetInterface
	*types.SResource
	Storage storage.IStep
}

func (d *SDaemonSet) Restart() error {
	return kuberetry.RetryOnConflict(kuberetry.DefaultRetry, func() error {
		path := fmt.Sprintf(`{"spec":{"template":{"metadata":{"annotations":{"kubectl.kubernetes.io/restartedAt":"%s"}}}}}`, time.Now().Format(time.RFC3339))
		_, err := d.Client.Patch(d.Context, d.SResource.GetName(), kubetypes.StrategicMergePatchType, []byte(path), metav1.PatchOptions{})
		if err != nil {
			return err
		}
		return nil
	})
}

func (d *SDaemonSet) Scale(_ int32) error { return nil }

func (d *SDaemonSet) Update() error {
	return kuberetry.RetryOnConflict(kuberetry.DefaultRetry, func() error {
		result, err := d.Client.Get(d.Context, d.GetName(), metav1.GetOptions{})
		if err != nil {
			if !errors.IsNotFound(err) {
				d.Storage.Log().Writef("get %s failed:%v", d.GetKind(), err)
			}
			return err
		}
		d.UpdateContainersImage(result.Spec.Template.Spec.Containers, func(str string) {
			d.Storage.Log().Write("update container image;", str)
		})
		d.UpdateEnvVariables(result.Spec.Template.Spec.Containers, func(str string) {
			d.Storage.Log().Write("update container env;", str)
		})

		if d.IgnoreInitContainer == nil || !*d.IgnoreInitContainer {
			d.UpdateContainersImage(result.Spec.Template.Spec.InitContainers, func(str string) {
				d.Storage.Log().Write("update init container image;", str)
			})
			d.UpdateEnvVariables(result.Spec.Template.Spec.InitContainers, func(str string) {
				d.Storage.Log().Write("update init container env;", str)
			})
		}
		_, err = d.Client.Update(d.Context, result, metav1.UpdateOptions{})
		if err != nil {
			return err
		}
		return nil
	})
}

func (d *SDaemonSet) Println() error {
	return kuberetry.RetryOnConflict(kuberetry.DefaultRetry, func() error {
		result, err := d.Client.Get(d.Context, d.GetName(), metav1.GetOptions{})
		if err != nil {
			if !errors.IsNotFound(err) {
				d.Storage.Log().Writef("get %s failed:%v", d.GetName(), err)
			}
			return err
		}
		for _, container := range result.Spec.Template.Spec.InitContainers {
			d.Storage.Log().Writef("%s/%s init container: %s", result.Namespace, result.Name, container.Image)
		}
		for _, container := range result.Spec.Template.Spec.Containers {
			d.Storage.Log().Writef("%s/%s container: %s", result.Namespace, result.Name, container.Image)
		}
		return nil
	})
}
