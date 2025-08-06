package k8s

import (
	"context"
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/cache"
	watchtools "k8s.io/client-go/tools/watch"

	"github.com/busyster996/dagflow/internal/storage"
	"github.com/busyster996/dagflow/internal/worker/runner/kubectl/types"
)

type sStatusViewer struct {
	ctx     context.Context
	lw      cache.ListerWatcher
	fn      func(obj runtime.Unstructured) (string, bool, error)
	storage storage.IStep
}

// Status StatusViewerFor returns a StatusViewer for the resource specified by kind.
// https://github.com/kubernetes/kubectl/blob/master/pkg/polymorphichelpers/rollout_status.go
func Status(ctx context.Context, storage storage.IStep, dynamicClient dynamic.Interface, resource *types.SResource) error {
	var s = &sStatusViewer{
		ctx:     ctx,
		storage: storage,
	}
	var gvr = schema.GroupVersionResource{
		Group:   "apps",
		Version: "v1",
	}
	switch resource.GetKind() {
	case types.Deployment:
		gvr.Resource = "deployments"
		s.fn = s.watchDeployment
	case types.DaemonSet:
		gvr.Resource = "daemonsets"
		s.fn = s.watchDaemonSet
	case types.StatefulSet:
		gvr.Resource = "statefulsets"
		s.fn = s.watchStatefulSet
	default:
		return fmt.Errorf("unsupported resource type %s", resource.GetKind())
	}
	fieldSelector := fields.OneTermEqualSelector("metadata.name", resource.GetName()).String()
	s.lw = &cache.ListWatch{
		ListWithContextFunc: func(ctx context.Context, options metav1.ListOptions) (runtime.Object, error) {
			options.FieldSelector = fieldSelector
			return dynamicClient.Resource(gvr).Namespace(resource.GetNamespace()).List(ctx, options)
		},
		WatchFuncWithContext: func(ctx context.Context, options metav1.ListOptions) (watch.Interface, error) {
			options.FieldSelector = fieldSelector
			return dynamicClient.Resource(gvr).Namespace(resource.GetNamespace()).Watch(ctx, options)
		},
	}
	s.status()
	return nil
}

func (s *sStatusViewer) status() {
	_, err := watchtools.UntilWithSync(s.ctx, s.lw, &unstructured.Unstructured{}, nil, func(e watch.Event) (bool, error) {
		switch t := e.Type; t {
		case watch.Added, watch.Modified:
			status, done, err := s.fn(e.Object.(runtime.Unstructured))
			if err != nil {
				return false, err
			}
			s.storage.Log().Write(status)
			if done {
				return true, nil
			}
			return false, nil
		case watch.Deleted:
			// We need to abort to avoid cases of recreation and not to silently watch the wrong (new) object
			return true, fmt.Errorf("object has been deleted")

		default:
			return true, fmt.Errorf("internal error: Unexpected event %#v", e)
		}
	})
	if err != nil {
		s.storage.Log().Writef("error: %s", err.Error())
		return
	}
}

func (s *sStatusViewer) getResourceCondition(status appsv1.DeploymentStatus, condType appsv1.DeploymentConditionType) *appsv1.DeploymentCondition {
	for i := range status.Conditions {
		c := status.Conditions[i]
		if c.Type == condType {
			return &c
		}
	}
	return nil
}

func (s *sStatusViewer) watchDeployment(obj runtime.Unstructured) (string, bool, error) {
	deployment := &appsv1.Deployment{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(obj.UnstructuredContent(), deployment)
	if err != nil {
		return "", false, fmt.Errorf("failed to convert %T to %T: %v", obj, deployment, err)
	}

	if deployment.Generation <= deployment.Status.ObservedGeneration {
		cond := s.getResourceCondition(deployment.Status, appsv1.DeploymentProgressing)
		if cond != nil && cond.Reason == types.TimedOutReason {
			return "", false, fmt.Errorf("deployment %s/%s exceeded its progress deadline",
				deployment.Namespace, deployment.Name)
		}
		if deployment.Spec.Replicas != nil && deployment.Status.UpdatedReplicas < *deployment.Spec.Replicas {
			return fmt.Sprintf("waiting for deployment %s/%s rollout to finish: %d out of %d new replicas have been updated...",
				deployment.Namespace, deployment.Name, deployment.Status.UpdatedReplicas, *deployment.Spec.Replicas), false, nil
		}
		if deployment.Status.Replicas > deployment.Status.UpdatedReplicas {
			return fmt.Sprintf("waiting for deployment %s/%s rollout to finish: %d old replicas are pending termination...",
				deployment.Namespace, deployment.Name, deployment.Status.Replicas-deployment.Status.UpdatedReplicas), false, nil
		}
		if deployment.Status.AvailableReplicas < deployment.Status.UpdatedReplicas {
			return fmt.Sprintf("waiting for deployment %s/%s rollout to finish: %d of %d updated replicas are available...",
				deployment.Namespace, deployment.Name, deployment.Status.UpdatedReplicas, deployment.Status.AvailableReplicas), false, nil
		}
		return fmt.Sprintf("deployment %s/%s successfully rolled out",
			deployment.Namespace, deployment.Name), true, nil
	}
	return fmt.Sprintf("waiting for deployment %s/%s spec update to be observed...",
		deployment.Namespace, deployment.Name), false, nil
}

func (s *sStatusViewer) watchDaemonSet(obj runtime.Unstructured) (string, bool, error) {
	daemon := &appsv1.DaemonSet{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(obj.UnstructuredContent(), daemon)
	if err != nil {
		return "", false, fmt.Errorf("failed to convert %T to %T: %v", obj, daemon, err)
	}
	if daemon.Spec.UpdateStrategy.Type != appsv1.RollingUpdateDaemonSetStrategyType {
		return "", true, fmt.Errorf("rollout status is only available for %s strategy type", appsv1.RollingUpdateStatefulSetStrategyType)
	}
	if daemon.Generation <= daemon.Status.ObservedGeneration {
		if daemon.Status.UpdatedNumberScheduled < daemon.Status.DesiredNumberScheduled {
			return fmt.Sprintf("waiting for daemon set %s/%s rollout to finish: %d out of %d new pods have been updated...",
				daemon.Namespace, daemon.Name, daemon.Status.UpdatedNumberScheduled, daemon.Status.DesiredNumberScheduled), false, nil
		}
		if daemon.Status.NumberAvailable < daemon.Status.DesiredNumberScheduled {
			return fmt.Sprintf("waiting for daemon set %s/%s rollout to finish: %d of %d updated pods are available...",
				daemon.Namespace, daemon.Name, daemon.Status.NumberAvailable, daemon.Status.DesiredNumberScheduled), false, nil
		}
		return fmt.Sprintf("daemon set %s/%s successfully rolled out",
			daemon.Namespace, daemon.Name), true, nil
	}
	return fmt.Sprintf("waiting for daemon set %s/%s spec update to be observed...",
		daemon.Namespace, daemon.Name), false, nil
}

func (s *sStatusViewer) watchStatefulSet(obj runtime.Unstructured) (string, bool, error) {
	sts := &appsv1.StatefulSet{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(obj.UnstructuredContent(), sts)
	if err != nil {
		return "", false, fmt.Errorf("failed to convert %T to %T: %v", obj, sts, err)
	}
	if sts.Spec.UpdateStrategy.Type != appsv1.RollingUpdateStatefulSetStrategyType {
		return "", true, fmt.Errorf("rollout status is only available for %s strategy type", appsv1.RollingUpdateStatefulSetStrategyType)
	}
	if sts.Status.ObservedGeneration == 0 || sts.Generation > sts.Status.ObservedGeneration {
		return fmt.Sprintf("waiting for stateful set %s/%s spec update to be observed...",
			sts.Namespace, sts.Name), false, nil
	}
	if sts.Spec.Replicas != nil && sts.Status.ReadyReplicas < *sts.Spec.Replicas {
		return fmt.Sprintf("waiting for stateful set %s/%s %d pods to be ready...",
			sts.Namespace, sts.Name, *sts.Spec.Replicas-sts.Status.ReadyReplicas), false, nil
	}
	if (sts.Spec.Replicas == nil || *sts.Spec.Replicas == 0) && sts.Status.UpdatedReplicas > 0 {
		return fmt.Sprintf("waiting for stateful set %s/%s %d pods to be delete...",
			sts.Namespace, sts.Name, sts.Status.UpdatedReplicas), false, nil
	}
	if sts.Spec.UpdateStrategy.Type == appsv1.RollingUpdateStatefulSetStrategyType && sts.Spec.UpdateStrategy.RollingUpdate != nil {
		if sts.Spec.Replicas != nil && sts.Spec.UpdateStrategy.RollingUpdate.Partition != nil {
			if sts.Status.UpdatedReplicas < (*sts.Spec.Replicas - *sts.Spec.UpdateStrategy.RollingUpdate.Partition) {
				return fmt.Sprintf("waiting for stateful set %s/%s  partitioned roll out to finish: %d out of %d new pods have been updated...",
					sts.Namespace, sts.Name, sts.Status.UpdatedReplicas, *sts.Spec.Replicas-*sts.Spec.UpdateStrategy.RollingUpdate.Partition), false, nil
			}
		}
		return fmt.Sprintf("stateful set %s/%s partitioned roll out complete: %d new pods have been updated...",
			sts.Namespace, sts.Name, sts.Status.UpdatedReplicas), true, nil
	}
	if sts.Status.UpdateRevision != sts.Status.CurrentRevision {
		return fmt.Sprintf("waiting for stateful set %s/%s rolling update to complete %d pods at revision %s...",
			sts.Namespace, sts.Name, sts.Status.UpdatedReplicas, sts.Status.UpdateRevision), false, nil
	}
	return fmt.Sprintf("stateful set %s/%s rolling update complete %d pods at revision %s...",
		sts.Namespace, sts.Name, sts.Status.CurrentReplicas, sts.Status.CurrentRevision), true, nil
}
