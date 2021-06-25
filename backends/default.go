package backends

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

type defaultBackend struct {
	Name       string
	Namespace  string
	SecretName string
}

func (backend defaultBackend) Deployment(image string) *appsv1.Deployment {
	var replicas int32 = 1
	return &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name:      backend.Name,
			Namespace: backend.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": backend.Name,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name: backend.Name,
					Labels: map[string]string{
						"app": backend.Name,
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Name:  "default",
						Image: image,
					}},
				},
			},
		},
	}
}

func (backend defaultBackend) Service(port int32) *corev1.Service {
	return &corev1.Service{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name:      backend.Name,
			Namespace: backend.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{{
				Name: "tcp",
				Port: port,
				TargetPort: intstr.IntOrString{
					Type:   intstr.Int,
					IntVal: port,
				},
			}},
			Selector: map[string]string{
				"app": backend.Name,
			},
		},
	}
}
