package balancer

import (
	balancerv1alpha1 `github.com/hliangzhao/balancer/pkg/apis/balancer/v1alpha1`
	appv1 `k8s.io/api/apps/v1`
)

func (r *ReconcileBalancer) syncDeployment(balancer *balancerv1alpha1.Balancer) error {
	// TODO
	return nil
}

func NewDeployment(balancer balancerv1alpha1.Balancer) (*appv1.Deployment, error) {
	// TODO
	return &appv1.Deployment{}, nil
}
