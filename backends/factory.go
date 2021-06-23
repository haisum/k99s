package backends

import (
	"errors"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

type BackendType string

const (
	MySQLBackend BackendType = "mysql"
)

type Backend interface {
	Deployment(image string) *appsv1.Deployment
	Service(port int32) *corev1.Service
}

func New(backendType BackendType, name, namespace string) (Backend, error) {
	switch backendType {
	case MySQLBackend:
		return newMySQLBackend(name, namespace), nil
	default:
		return nil, errors.New("backend not implemented yet")
	}
}
