package runtimes

import (
	"fmt"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

type defaultRuntime struct {
	Name      string
	Namespace string
}

func (r defaultRuntime) Deployment(image string) *appsv1.Deployment {
	var replicas int32 = 1
	return &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name:      r.Name,
			Namespace: r.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": r.Name,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name: "",
					Labels: map[string]string{
						"app": r.Name,
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Name:  "default",
						Image: image,
						EnvFrom: []corev1.EnvFromSource{
							{
								SecretRef: &corev1.SecretEnvSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: fmt.Sprintf("%s-credentials", r.Name),
									},
								},
							},
							{
								ConfigMapRef: &corev1.ConfigMapEnvSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: fmt.Sprintf("%s-config", r.Name),
									},
								},
							}},
					}},
				},
			},
		},
	}
}

func (r defaultRuntime) Service(port int32) *corev1.Service {
	return &corev1.Service{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name:      r.Name,
			Namespace: r.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{{
				Name: "http",
				Port: port,
				TargetPort: intstr.IntOrString{
					Type:   intstr.Int,
					IntVal: port,
				},
			}},
			Selector: map[string]string{
				"app": r.Name,
			},
		},
	}
}

func (r defaultRuntime) Ingress(domain string) *networkingv1.Ingress {
	pathTypePrefix := networkingv1.PathTypePrefix
	return &networkingv1.Ingress{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name:      r.Name,
			Namespace: r.Namespace,
		},
		Spec: networkingv1.IngressSpec{
			Rules: []networkingv1.IngressRule{{
				Host: fmt.Sprintf("%s.%s", r.Name, domain),
				IngressRuleValue: networkingv1.IngressRuleValue{
					HTTP: &networkingv1.HTTPIngressRuleValue{
						Paths: []networkingv1.HTTPIngressPath{{
							Path:     "/",
							PathType: &pathTypePrefix,
							Backend: networkingv1.IngressBackend{
								Service: &networkingv1.IngressServiceBackend{
									Name: r.Name,
									Port: networkingv1.ServiceBackendPort{
										Name: "http",
									},
								},
							},
						}},
					},
				},
			}},
		},
	}
}
