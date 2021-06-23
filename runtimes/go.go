package runtimes

import (
	appsv1 "k8s.io/api/apps/v1"
)

type goRuntime struct {
	defaultRuntime
}

func (p goRuntime) Deployment(image string) *appsv1.Deployment {
	deployment := p.defaultRuntime.Deployment(image)
	deployment.Spec.Template.Spec.Containers[0].Name = "go"
	return deployment
}

func newGoRuntime(name, namespace string) Runtime {
	return goRuntime{defaultRuntime{Name: name, Namespace: namespace}}
}
