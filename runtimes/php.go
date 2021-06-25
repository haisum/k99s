package runtimes

import (
	appsv1 "k8s.io/api/apps/v1"
)

type phpRuntime struct {
	defaultRuntime
}

func (p phpRuntime) Deployment(image string) *appsv1.Deployment {
	deployment := p.defaultRuntime.Deployment(image)
	deployment.Spec.Template.Spec.Containers[0].Name = "php"
	deployment.Spec.Template.ObjectMeta.Labels["type"] = "php"
	return deployment
}

func newPHPRuntime(name, namespace string) Runtime {
	return phpRuntime{defaultRuntime{
		Name:      name,
		Namespace: namespace,
	}}
}
