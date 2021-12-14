package balancer

import (
	balancerv1alpha1 `github.com/hliangzhao/balancer/pkg/apis/balancer/v1alpha1`
	corev1 `k8s.io/api/core/v1`
)

func (r *ReconcileBalancer) syncService(balancer *balancerv1alpha1.Balancer) error {
	// TODO
	return nil
}

func NewService(balancer balancerv1alpha1.Balancer) (*corev1.Service, error) {
	// TODO
	return &corev1.Service{}, nil
}
