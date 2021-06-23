package backends

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

type defaultBackend struct {
	Name      string
	Namespace string
}

func (backend defaultBackend) Deployment(image string) *appsv1.Deployment {

}

func (backend defaultBackend) Service(port int32) *corev1.Service {

}
