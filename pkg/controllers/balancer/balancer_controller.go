/*
Copyright 2021 hliangzhao.

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

package balancer

import (
	"context"
	balancerv1alpha1 "github.com/hliangzhao/balancer/pkg/apis/balancer/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("balancer-controller")

// ReconcilerBalancer reconciles a Balancer instance. Reconciler is the core of a controller.
type ReconcilerBalancer struct {
	// client reads obj from the cache
	client client.Client
	scheme *runtime.Scheme
}

// newReconciler creates the ReconcilerBalancer with input controller-manager.
func newReconciler(manager manager.Manager) reconcile.Reconciler {
	return &ReconcilerBalancer{
		client: manager.GetClient(),
		scheme: manager.GetScheme(),
	}
}

// addReconciler adds r to controller-manager.
func addReconciler(manager manager.Manager, r reconcile.Reconciler) error {
	// creates a balancer-controller registered in controller-manager
	c, err := controller.New("balancer-controller", manager, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// takes events provided by a Source and uses the EventHandler to enqueue reconcile.Requests in response to the events.
	if err = c.Watch(&source.Kind{Type: &balancerv1alpha1.Balancer{}}, &handler.EnqueueRequestForObject{}); err != nil {
		return err
	}
	// the changes of the configmap, pod, and svc which are created by balancer will also be enqueued
	if err = c.Watch(&source.Kind{Type: &corev1.ConfigMap{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &balancerv1alpha1.Balancer{}},
	); err != nil {
		return err
	}
	if err = c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &balancerv1alpha1.Balancer{}},
	); err != nil {
		return err
	}
	if err = c.Watch(&source.Kind{Type: &corev1.Service{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &balancerv1alpha1.Balancer{}},
	); err != nil {
		return err
	}

	return nil
}

// Add creates a newly registered balancer-controller to controller-manager.
func Add(manager manager.Manager) error {
	return addReconciler(manager, newReconciler(manager))
}

// Here we provide a static check that ReconcilerBalancer satisfies reconcile.Reconciler interface.
// The _ used as a name of the variable tells the compiler to effectively discard the RHS value,
// but to type-check it and evaluate it if it has any side effects, but the anonymous variable per
// se doesn't take any process space.
var _ reconcile.Reconciler = &ReconcilerBalancer{}

// Reconcile reads the status of the Balancer object and makes changes toward to Balancer.Spec.
// This func must be implemented to be a legal reconcile.Reconciler!
func (r *ReconcilerBalancer) Reconcile(context context.Context, request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling Balancer")

	// fetch the expected Balancer instance through the client
	balancer := &balancerv1alpha1.Balancer{}
	if err := r.client.Get(context, request.NamespacedName, balancer); err != nil {
		// balancer not exist
		if errors.IsNotFound(err) {
			// the namespaced name in request is not found, return empty result and requeue the request
			return reconcile.Result{}, nil
		}
	}

	// Founded. Update SVCs, deployments, etc. according to the expected Balancer.
	// If any error happens, the request would be requeue
	if err := r.syncFrontendService(balancer); err != nil {
		return reconcile.Result{}, err
	}
	if err := r.syncDeployment(balancer); err != nil {
		return reconcile.Result{}, nil
	}
	if err := r.syncBackendServices(balancer); err != nil {
		return reconcile.Result{}, nil
	}
	if err := r.syncBalancerStatus(balancer); err != nil {
		return reconcile.Result{}, nil
	}

	return reconcile.Result{}, nil
}
