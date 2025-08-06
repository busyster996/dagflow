package types

import (
	"fmt"
	"strings"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type SResource struct {
	Kind                string `json:"kind"`
	Namespace           string `json:"namespace"`
	Name                string `json:"name"`
	Replicas            *int32 `json:"replicas"`
	ImageTag            string `json:"imageTag"`
	IgnoreInitContainer *bool  `json:"ignoreInitContainer"`
	Env                 []Env  `json:"env"`
}

type Env struct {
	Containers []string `json:"containers"`
	Operator   Operator `json:"operator"`
	corev1.EnvVar
}

func (r *SResource) Check() error {
	if r.Namespace == "" {
		r.Namespace = metav1.NamespaceDefault
	}
	if r.Name == "" {
		return fmt.Errorf("name is required")
	}
	if r.Kind == "" {
		r.Kind = Deployment
	}
	switch r.Kind {
	case Deployment, DaemonSet, StatefulSet:
	default:
		return fmt.Errorf("kind must be one of %s, %s, %s", Deployment, DaemonSet, StatefulSet)
	}
	return nil
}

func (r *SResource) GetKind() string {
	return r.Kind
}

func (r *SResource) GetNamespace() string {
	return r.Namespace
}

func (r *SResource) GetName() string {
	return r.Name
}

func (r *SResource) GetReplicas() int32 {
	if r.Replicas == nil {
		return 0
	}
	return *r.Replicas
}

func (r *SResource) GetImageTag() string {
	return r.ImageTag
}

func (r *SResource) UpdateContainersImage(containers []corev1.Container, logger func(str string)) {
	if r.ImageTag == "" {
		return
	}
	for i, container := range containers {
		imageParts := strings.Split(container.Image, ":")
		if len(imageParts) != 2 {
			continue
		}
		oldTag := imageParts[1]
		newImage := fmt.Sprintf("%s:%s", imageParts[0], r.ImageTag)
		logger(fmt.Sprintf("%s %s %s %s -> %s", r.GetKind(), r.GetName(), container.Name, oldTag, r.ImageTag))
		containers[i].Image = newImage
	}
}

func (r *SResource) UpdateEnvVariables(containers []corev1.Container, logger func(str string)) {
	if r.Env == nil {
		return
	}
	for i, ctr := range containers {
		for _, env := range r.GetEnvs() {
			if !env.checkCtrName(ctr.Name) {
				continue
			}
			switch env.Operator {
			case OPERATOR_ADD:
				logger(fmt.Sprintf("Adding %s env var to %s", env.EnvVar.String(), ctr.Name))
				containers[i].Env = env.addOrUpdate(containers[i].Env)
			case OPERATOR_DELETE:
				logger(fmt.Sprintf("Deleting env var from %s", ctr.Name))
				containers[i].Env = env.remove(containers[i].Env)
			default:
				logger(fmt.Sprintf("Unsupported operator %s", env.Operator))
			}
		}
	}
}

func (r *SResource) GetEnvs() []Env {
	return r.Env
}

// 检查容器名称
func (e Env) checkCtrName(name string) bool {
	if len(e.Containers) == 0 {
		return true
	}
	for _, ctr := range e.Containers {
		if ctr == name {
			return true
		}
	}
	return false
}

func (e Env) addOrUpdate(envs []corev1.EnvVar) []corev1.EnvVar {
	var exist bool
	// 检查是否已存在, 如果存在则更新, 不存在则添加
	for k, env := range envs {
		if env.Name == e.Name {
			envs[k] = e.EnvVar
			exist = true
		}
	}
	if !exist {
		envs = append(envs, e.EnvVar)
	}

	return envs
}

func (e Env) remove(envs []corev1.EnvVar) []corev1.EnvVar {
	for k, env := range envs {
		if env.Name == e.Name {
			envs = append(envs[:k], envs[k+1:]...)
		}
	}
	return envs
}
