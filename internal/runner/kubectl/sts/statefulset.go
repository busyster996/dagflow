package sts

import (
	"context"
	"fmt"
	"time"

	autoscalingv1 "k8s.io/api/autoscaling/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubetypes "k8s.io/apimachinery/pkg/types"
	appv1 "k8s.io/client-go/kubernetes/typed/apps/v1"
	kuberetry "k8s.io/client-go/util/retry"

	"github.com/busyster996/dagflow/internal/runner/kubectl/types"
	"github.com/busyster996/dagflow/internal/storage"
)

type SStatefulSet struct {
	context.Context
	*types.SResource
	Client  appv1.StatefulSetInterface
	Storage storage.IStep
}

func (s *SStatefulSet) Restart() error {
	return kuberetry.RetryOnConflict(kuberetry.DefaultRetry, func() error {
		path := fmt.Sprintf(`{"spec":{"template":{"metadata":{"annotations":{"kubectl.kubernetes.io/restartedAt":"%s"}}}}}`, time.Now().Format(time.RFC3339))
		_, err := s.Client.Patch(s.Context, s.GetName(), kubetypes.StrategicMergePatchType, []byte(path), metav1.PatchOptions{})
		if err != nil {
			return err
		}
		return nil
	})
}

func (s *SStatefulSet) Scale(replicas int32) error {
	return kuberetry.RetryOnConflict(kuberetry.DefaultRetry, func() error {
		scale := &autoscalingv1.Scale{
			Spec: autoscalingv1.ScaleSpec{
				Replicas: replicas,
			},
		}
		scale.Namespace = s.GetNamespace()
		scale.Name = s.GetName()
		_, err := s.Client.UpdateScale(s.Context, s.GetName(), scale, metav1.UpdateOptions{})
		if err != nil {
			return err
		}
		return nil
	})
}

func (s *SStatefulSet) Update() error {
	return kuberetry.RetryOnConflict(kuberetry.DefaultRetry, func() error {
		result, err := s.Client.Get(s.Context, s.GetName(), metav1.GetOptions{})
		if err != nil {
			if !errors.IsNotFound(err) {
				s.Storage.Log().Writef("get %s failed:%v", s.GetKind(), err)
			}
			return err
		}
		s.UpdateContainersImage(result.Spec.Template.Spec.Containers, func(str string) {
			s.Storage.Log().Write("update container image;", str)
		})
		s.UpdateEnvVariables(result.Spec.Template.Spec.Containers, func(str string) {
			s.Storage.Log().Write("update container env;", str)
		})
		if s.IgnoreInitContainer == nil || !*s.IgnoreInitContainer {
			s.UpdateContainersImage(result.Spec.Template.Spec.InitContainers, func(str string) {
				s.Storage.Log().Write("update init container image;", str)
			})
			s.UpdateEnvVariables(result.Spec.Template.Spec.InitContainers, func(str string) {
				s.Storage.Log().Write("update init container env;", str)
			})
		}

		_, err = s.Client.Update(s.Context, result, metav1.UpdateOptions{})
		if err != nil {
			return err
		}
		return nil
	})
}

func (s *SStatefulSet) Println() error {
	return kuberetry.RetryOnConflict(kuberetry.DefaultRetry, func() error {
		result, err := s.Client.Get(s.Context, s.GetName(), metav1.GetOptions{})
		if err != nil {
			if !errors.IsNotFound(err) {
				s.Storage.Log().Writef("get %s failed:%v", s.GetKind(), err)
			}
			return err
		}
		for _, container := range result.Spec.Template.Spec.InitContainers {
			s.Storage.Log().Writef("%s/%s init container: %s", result.Namespace, result.Name, container.Image)
		}
		for _, container := range result.Spec.Template.Spec.Containers {
			s.Storage.Log().Writef("%s/%s container: %s", result.Namespace, result.Name, container.Image)
		}
		return nil
	})
}
