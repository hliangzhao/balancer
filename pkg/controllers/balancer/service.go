package balancer

import (
	`context`
	balancerv1alpha1 `github.com/hliangzhao/balancer/pkg/apis/balancer/v1alpha1`
	corev1 `k8s.io/api/core/v1`
	`k8s.io/apimachinery/pkg/api/errors`
	metav1 `k8s.io/apimachinery/pkg/apis/meta/v1`
	`k8s.io/apimachinery/pkg/types`
	`sigs.k8s.io/controller-runtime/pkg/controller/controllerutil`
)

// syncService sync the service that created by balancer.
// If the service of balancer is not found, create it with the newest service;
// Otherwise, update it with the newest service.
func (r *ReconcileBalancer) syncService(balancer *balancerv1alpha1.Balancer) error {
	svc, err := NewService(balancer)
	if err != nil {
		return err
	}

	// set balancer as the owner of svc
	if err := controllerutil.SetOwnerReference(balancer, svc, r.scheme); err != nil {
		return err
	}

	found := &corev1.Service{}
	err = r.client.Get(context.Background(), types.NamespacedName{Namespace: svc.Namespace, Name: svc.Name}, found)
	if err != nil && errors.IsNotFound(err) {
		// corresponding service not found in the cluster, create it with the newest svc
		if err = r.client.Create(context.Background(), svc); err != nil {
			return err
		}
		return nil
	} else if err != nil {
		return err
	}

	// corresponding service found, update it with the newest svc
	if err = r.client.Update(context.Background(), svc); err != nil {
		return err
	}
	return nil
}

// NewService creates a new Service which exposes all the ports exposed in the Balancer instance.
func NewService(balancer *balancerv1alpha1.Balancer) (*corev1.Service, error) {
	var balancerPorts []corev1.ServicePort
	for _, port := range balancer.Spec.Ports {
		balancerPorts = append(balancerPorts, corev1.ServicePort{
			Name:     port.Name,
			Protocol: corev1.Protocol(port.Protocol),
			Port:     int32(port.Port),
			// TODO: why not assign targetPort?
			// TargetPort:  port.TargetPort,
		})
	}

	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      balancer.Name,
			Namespace: balancer.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Selector: NewPodLabels(balancer),
			Type:     corev1.ServiceTypeClusterIP,
			Ports:    balancerPorts,
		},
	}, nil
}
