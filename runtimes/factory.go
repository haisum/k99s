package runtimes

import (
	"errors"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
)

type Runtime interface {
	Deployment(image string) *appsv1.Deployment
	Service(port int32) *corev1.Service
	Ingress(domain string) *networkingv1.Ingress
}

type RuntimeType string

const (
	PHPRuntime = "php"
	GoRuntime  = "go"
)

func New(runtime RuntimeType, name, namespace string) (Runtime, error) {
	switch runtime {
	case PHPRuntime:
		return newPHPRuntime(name, namespace), nil
	case GoRuntime:
		return newGoRuntime(name, namespace), nil
	default:
		return nil, errors.New("runtime not supported")
	}
}
