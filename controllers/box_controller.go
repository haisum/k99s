/*
Copyright 2021.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"fmt"
	"github.com/go-logr/logr"
	paasv1 "github.com/haisum/k99s/api/v1"
	"github.com/haisum/k99s/backends"
	"github.com/haisum/k99s/runtimes"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"math/rand"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"time"
)

// BoxReconciler reconciles a Box object
type BoxReconciler struct {
	client.Client
	Scheme *runtime.Scheme
	Log    logr.Logger
}

type NullableObject struct {
	IsSet  bool
	Object client.Object
}

var (
	alphanumericChars = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	backendImages     = map[backends.BackendType]string{
		"mysql": "mysql:5.7",
	}
	runtimeImages = map[runtimes.RuntimeType]string{
		"php": "k3d-registry.localhost:5000/php:1.0",
		"go":  "k3d-registry.localhost:5000/go:1.0",
	}
)

const (
	domain       = "k99s-paas.com"
	backendPort  = 3306
	frontendPort = 80
)

//+kubebuilder:rbac:groups=paas.example.com,resources=boxes,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=paas.example.com,resources=boxes/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=paas.example.com,resources=boxes/finalizers,verbs=update
//+kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups="",resources=configmaps,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups="",resources=services,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups="apps/v1",resources=deployments,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups="networking/v1",resources=ingresses,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Box object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.8.3/pkg/reconcile
func (r *BoxReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("box", req.NamespacedName)
	var box paasv1.Box
	if err := r.Get(ctx, req.NamespacedName, &box); err != nil {
		log.Error(err, "unable to fetch Box")
		// we'll ignore not-found errors, since they can't be fixed by an immediate
		// requeue (we'll need to wait for a new notification), and we can get them
		// on deleted requests.
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	if box.Status.Status == "" {
		box.Status.Status = "Inactive"
	}
	if err := r.createConfigMapIfDoesNotExist(&box, ctx, req); err != nil {
		box.Status.Error = err.Error()
		log.Error(err, "unable to create configmap for Box")
		return ctrl.Result{}, err
	}
	if err := r.createSecretIfDoesNotExist(&box, ctx, req); err != nil {
		box.Status.Error = err.Error()
		log.Error(err, "unable to create secret for Box")
		return ctrl.Result{}, err
	}

	// create backend Objects
	backend, err := backends.New(box.Spec.Backend, req.Name, req.Namespace, req.Name+"-credentials")
	if err != nil {
		box.Status.Error = err.Error()
		return ctrl.Result{}, err
	}

	// backend deployment
	backendDeployment := backend.Deployment(backendImages[box.Spec.Backend])
	if err := r.createOrUpdate(NullableObject{Object: backendDeployment}, &box, ctx); err != nil {
		box.Status.Error = err.Error()
		log.Error(err, "unable to create backend deployment for box", "deployment", backendDeployment)
		return ctrl.Result{}, err
	}

	// backend service
	backendService := backend.Service(backendPort)
	if err := r.createOrUpdate(NullableObject{Object: backendService}, &box, ctx); err != nil {
		box.Status.Error = err.Error()
		log.Error(err, "unable to create Backend service for Box", "service", backendService)
		return ctrl.Result{}, err
	}

	// create frontend Objects
	frontend, err := runtimes.New(box.Spec.Runtime, req.Name, req.Namespace)
	if err != nil {
		box.Status.Error = err.Error()
		return ctrl.Result{}, err
	}

	// frontend deployment
	deployment := frontend.Deployment(runtimeImages[box.Spec.Runtime])
	if err := r.createOrUpdate(NullableObject{Object: deployment}, &box, ctx); err != nil {
		box.Status.Error = err.Error()
		log.Error(err, "unable to create Deployment for Box", "deployment", deployment)
		return ctrl.Result{}, err
	}

	// frontend service
	service := frontend.Service(frontendPort)
	if err := r.createOrUpdate(NullableObject{Object: service}, &box, ctx); err != nil {
		box.Status.Error = err.Error()
		log.Error(err, "unable to create Service for Box", "service", service)
		return ctrl.Result{}, err
	}

	// frontend Ingress
	ingress := frontend.Ingress(domain)
	if err := r.createOrUpdate(NullableObject{Object: ingress}, &box, ctx); err != nil {
		box.Status.Error = err.Error()
		log.Error(err, "unable to create Ingress for Box", "ingress", ingress)
		return ctrl.Result{}, err
	}

	box.Status.Error = ""
	box.Status.Status = "Active"
	box.Status.URL = fmt.Sprintf("http://%s.%s", req.Name, domain)
	r.Client.Status().Update(context.Background(), &box)
	log.V(1).Info("reconciled Deployment for Box", "deployment", deployment)
	return ctrl.Result{
		RequeueAfter: time.Minute,
	}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *BoxReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&paasv1.Box{}).
		Owns(&corev1.Service{}).
		Owns(&appsv1.Deployment{}).
		Owns(&networkingv1.Ingress{}).
		Owns(&corev1.ConfigMap{}).
		Owns(&corev1.Secret{}).
		Complete(r)
}

func (r *BoxReconciler) createConfigMapIfDoesNotExist(box *paasv1.Box,
	ctx context.Context, req ctrl.Request) error {
	// create configmap
	configmap := corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-config", req.Name),
			Namespace: req.Namespace,
		},
		Data: map[string]string{
			"DB_HOST":     fmt.Sprintf("%s-backend", req.Name),
			"APP_URL":     fmt.Sprintf("%s.%s", req.Name, domain),
			"GIT_SUBPATH": box.Spec.GitSubPath,
			"GIT_URL":     box.Spec.GitURL,
		},
	}
	return r.createOrUpdate(NullableObject{Object: &configmap}, box, ctx)
}

func (r *BoxReconciler) createSecretIfDoesNotExist(box *paasv1.Box,
	ctx context.Context, req ctrl.Request) error {
	secret := corev1.Secret{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-credentials", req.Name),
			Namespace: req.Namespace,
		},
		StringData: map[string]string{
			"DB_PASSWORD":      randSeq(20),
			"DB_USER":          randSeq(10),
			"DB_BOOTSTRAP_SQL": box.Spec.BootstrapSQL,
			"DB_NAME":          req.Name,
		},
	}
	return r.createOrUpdate(NullableObject{Object: &secret}, box, ctx)
}

func (r *BoxReconciler) createOrUpdate(nullableObject NullableObject, box *paasv1.Box,
	ctx context.Context) error {
	namespacedName := types.NamespacedName{
		Namespace: nullableObject.Object.GetNamespace(),
		Name:      nullableObject.Object.GetName(),
	}
	nullableObject.IsSet = true
	if err := r.Get(ctx, namespacedName, nullableObject.Object.DeepCopyObject().(client.Object)); err != nil {
		if errors.IsNotFound(err) {
			nullableObject.IsSet = false
		} else {
			return err
		}
	}
	if nullableObject.IsSet {
		return nil
	}
	if err := ctrl.SetControllerReference(box, nullableObject.Object, r.Scheme); err != nil {
		return err
	}
	return r.Create(ctx, nullableObject.Object)
}

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = alphanumericChars[rand.Intn(len(alphanumericChars))]
	}
	return string(b)
}
