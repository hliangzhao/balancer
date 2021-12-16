package balancer

import (
	`context`
	balancerv1alpha1 `github.com/hliangzhao/balancer/pkg/apis/balancer/v1alpha1`
	corev1 `k8s.io/api/core/v1`
	`k8s.io/apimachinery/pkg/api/errors`
	`k8s.io/apimachinery/pkg/runtime`
	`sigs.k8s.io/controller-runtime/pkg/client`
	`sigs.k8s.io/controller-runtime/pkg/controller`
	`sigs.k8s.io/controller-runtime/pkg/handler`
	`sigs.k8s.io/controller-runtime/pkg/log`
	`sigs.k8s.io/controller-runtime/pkg/manager`
	`sigs.k8s.io/controller-runtime/pkg/reconcile`
	`sigs.k8s.io/controller-runtime/pkg/source`
)

var logger = log.Log.WithName("controller_balancer")

// ReconcileBalancer reconciles a Balancer instance. Reconciler is the core of a controller.
type ReconcileBalancer struct {
	// client reads obj from the cache
	client client.Client
	scheme *runtime.Scheme
}

// newReconciler creates the ReconcileBalancer with input controller-manager.
func newReconciler(manager manager.Manager) reconcile.Reconciler {
	return &ReconcileBalancer{
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

	// TODO: why not watch deployment?
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

// provide a static check that ReconcileBalancer satisfies reconcile.Reconciler interface.
// The _ used as a name of the variable tells the compiler to effectively discard the RHS value,
// but to type-check it and evaluate it if it has any side effects, but the anonymous variable per
// se doesn't take any process space.
var _ reconcile.Reconciler = &ReconcileBalancer{}

// Reconcile reads the status of the Balancer object and makes changes toward to Balancer.Spec.
// This func must be implemented to be a legal reconcile.Reconciler!
func (r *ReconcileBalancer) Reconcile(context context.Context, request reconcile.Request) (reconcile.Result, error) {
	reqLogger := logger.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling Balancer")

	// fetch the Balancer instance through the client
	balancer := &balancerv1alpha1.Balancer{}
	if err := r.client.Get(context, request.NamespacedName, balancer); err != nil {
		// balancer not exist
		if errors.IsNotFound(err) {
			// the namespaced name in request is not found, return empty result and requeue the request
			return reconcile.Result{}, nil
		}
	}

	// Founded. Update it.
	// If any error happens, the request would be requeue
	if err := r.syncBalancerStatus(balancer); err != nil {
		return reconcile.Result{}, nil
	}
	if err := r.syncServers(balancer); err != nil {
		return reconcile.Result{}, nil
	}
	if err := r.syncDeployment(balancer); err != nil {
		return reconcile.Result{}, nil
	}
	if err := r.syncService(balancer); err != nil {
		return reconcile.Result{}, err
	}

	return reconcile.Result{}, nil
}
