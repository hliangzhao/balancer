package balancer

import (
	balancerv1alpha1 "github.com/hliangzhao/balancer/pkg/apis/balancer/v1alpha1"
)

func NewPodLabels(balancer *balancerv1alpha1.Balancer) map[string]string {
	return map[string]string{
		balancerv1alpha1.BalancerKey: balancer.Name,
	}
}

func NewServiceLabels(balancer *balancerv1alpha1.Balancer) map[string]string {
	return map[string]string{
		balancerv1alpha1.BalancerKey: balancer.Name,
	}
}
