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
	"math/rand"

	"github.com/go-logr/logr"
	paasv1 "github.com/haisum/k99s/api/v1"
	"github.com/haisum/k99s/backends"
	runtimes "github.com/haisum/k99s/runtimes"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// BoxReconciler reconciles a Box object
type BoxReconciler struct {
	client.Client
	Scheme *runtime.Scheme
	Log    logr.Logger
}

const (
	domain          = "k99s-pass.com"
	phpRuntimeImage = "abc.com"
	goRuntimeImage  = "go.com"
	appPort         = 8080
	backendPort     = 3306
)

//+kubebuilder:rbac:groups=paas.example.com,resources=boxes,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=paas.example.com,resources=boxes/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=paas.example.com,resources=boxes/finalizers,verbs=update

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
	// create configmap
	configmap := &corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-config", req.Name),
			Namespace: req.Namespace,
		},
		Data: map[string]string{
			"DB_HOST": fmt.Sprintf("%s-backend", req.Name),
			"URL":     fmt.Sprintf("%s.%s", req.Name, domain),
		},
	}
	if err := ctrl.SetControllerReference(&box, configmap, r.Scheme); err != nil {
		return ctrl.Result{}, err
	}
	if err := r.Create(ctx, configmap); err != nil {
		log.Error(err, "unable to create configmap for Box", "configmap", configmap)
		return ctrl.Result{}, err
	}
	// create secret
	secret := &corev1.Secret{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-credentials", req.Name),
			Namespace: req.Namespace,
		},
		StringData: map[string]string{
			"DB_PASSWORD": randSeq(20),
			"DB_USERNAME": randSeq(10),
		},
	}
	if err := ctrl.SetControllerReference(&box, secret, r.Scheme); err != nil {
		return ctrl.Result{}, err
	}
	if err := r.Create(ctx, secret); err != nil {
		log.Error(err, "unable to create secret for Box")
		return ctrl.Result{}, err
	}

	// create backend resources
	backend, err := backends.New(box.Spec.Backend, req.Name, req.Namespace)
	if err != nil {
		return ctrl.Result{}, err
	}

	// backend deployment
	backendDeployment := backend.Deployment("")
	if err := ctrl.SetControllerReference(&box, backendDeployment, r.Scheme); err != nil {
		return ctrl.Result{}, err
	}
	if err := r.Create(ctx, backendDeployment); err != nil {
		log.Error(err, "unable to create Backend deployment for Box", "deployment", backendDeployment)
		return ctrl.Result{}, err
	}

	// backend service
	backendService := backend.Service(3306)
	if err := ctrl.SetControllerReference(&box, backendService, r.Scheme); err != nil {
		return ctrl.Result{}, err
	}
	if err := r.Create(ctx, backendService); err != nil {
		log.Error(err, "unable to create Backend service for Box", "service", backendService)
		return ctrl.Result{}, err
	}

	// create runtime Objects
	runtime, err := runtimes.New(box.Spec.Runtime, req.Name, req.Namespace)
	if err != nil {
		return ctrl.Result{}, err
	}

	// runtime deployment
	runtimeDeployment := runtime.Deployment("")
	if err := ctrl.SetControllerReference(&box, runtimeDeployment, r.Scheme); err != nil {
		return ctrl.Result{}, err
	}
	if err := r.Create(ctx, runtimeDeployment); err != nil {
		log.Error(err, "unable to create Deployment for Box", "deployment", runtimeDeployment)
		return ctrl.Result{}, err
	}

	// runtime service
	runTimeService := runtime.Service(8080)
	if err := ctrl.SetControllerReference(&box, runTimeService, r.Scheme); err != nil {
		return ctrl.Result{}, err
	}
	if err := r.Create(ctx, runTimeService); err != nil {
		log.Error(err, "unable to create Service for Box", "service", runTimeService)
		return ctrl.Result{}, err
	}

	// runtime Ingress
	runTimeIngress := runtime.Ingress("k99s-pass.com")
	if err := ctrl.SetControllerReference(&box, runTimeIngress, r.Scheme); err != nil {
		return ctrl.Result{}, err
	}
	if err := r.Create(ctx, runTimeIngress); err != nil {
		log.Error(err, "unable to create Ingress for Box", "ingress", runTimeIngress)
		return ctrl.Result{}, err
	}

	log.V(1).Info("created Deployment for Box", "deployment", runtimeDeployment)
	return ctrl.Result{}, nil
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// SetupWithManager sets up the controller with the Manager.
func (r *BoxReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&paasv1.Box{}).
		Complete(r)
}
