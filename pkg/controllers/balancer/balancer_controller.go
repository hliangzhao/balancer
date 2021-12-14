package balancer

import (
	`context`
	balancerv1alpha1 `github.com/hliangzhao/balancer/pkg/apis/balancer/v1alpha1`
	`k8s.io/apimachinery/pkg/api/errors`
	`k8s.io/apimachinery/pkg/runtime`
	`sigs.k8s.io/controller-runtime/pkg/client`
	`sigs.k8s.io/controller-runtime/pkg/log`
	`sigs.k8s.io/controller-runtime/pkg/reconcile`
)

var logger = log.Log.WithName("controller_balancer")

// ReconcileBalancer reconciles a Balancer instance.
type ReconcileBalancer struct {
	// client reads obj from the cache
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads the status of the Balancer object and makes changes toward to Balancer.Spec.
func (r *ReconcileBalancer) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := logger.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling Balancer")

	// fetch the Balancer instance through the client
	balancer := &balancerv1alpha1.Balancer{}
	if err := r.client.Get(context.Background(), request.NamespacedName, balancer); err != nil {
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
