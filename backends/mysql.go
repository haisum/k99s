package backends

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

type MySQL struct {
	defaultBackend
}

func newMySQLBackend(name, namespace, secretName string) Backend {
	return MySQL{defaultBackend{
		Name:       name,
		Namespace:  namespace,
		SecretName: secretName,
	}}
}

func (m MySQL) Deployment(image string) *appsv1.Deployment {
	deployment := m.defaultBackend.Deployment(image)
	deployment.Spec.Template.Spec.Containers[0].Name = "mysql"
	deployment.Spec.Template.ObjectMeta.Labels["type"] = "mysql"
	deployment.Spec.Template.Spec.Containers[0].Env = []corev1.EnvVar{
		{
			Name: "MYSQL_ROOT_PASSWORD",
			ValueFrom: &corev1.EnvVarSource{SecretKeyRef: &corev1.SecretKeySelector{
				LocalObjectReference: corev1.LocalObjectReference{Name: m.SecretName},
				Key:                  "DB_PASSWORD",
			}},
		},
		{
			Name: "MYSQL_PASSWORD",
			ValueFrom: &corev1.EnvVarSource{SecretKeyRef: &corev1.SecretKeySelector{
				LocalObjectReference: corev1.LocalObjectReference{Name: m.SecretName},
				Key:                  "DB_PASSWORD",
			}},
		},
		{
			Name: "MYSQL_USER",
			ValueFrom: &corev1.EnvVarSource{SecretKeyRef: &corev1.SecretKeySelector{
				LocalObjectReference: corev1.LocalObjectReference{Name: m.SecretName},
				Key:                  "DB_USER",
			}},
		},
		{
			Name: "MYSQL_DATABASE",
			ValueFrom: &corev1.EnvVarSource{SecretKeyRef: &corev1.SecretKeySelector{
				LocalObjectReference: corev1.LocalObjectReference{Name: m.SecretName},
				Key:                  "DB_NAME",
			}},
		},
	}
	deployment.Spec.Template.Spec.Volumes = []corev1.Volume{
		{
			Name: "bootstrap-sql",
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName: m.SecretName,
				},
			},
		},
	}
	deployment.Spec.Template.Spec.Containers[0].VolumeMounts = []corev1.VolumeMount{
		{
			Name:      "bootstrap-sql",
			MountPath: "/docker-entrypoint-initdb.d/bootstrap.sql",
			SubPath:   "DB_BOOTSTRAP_SQL",
		},
	}
	return deployment
}
