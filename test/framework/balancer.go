package framework

import (
	"context"
	"fmt"
	hliangzhaov1alpha1 "github.com/hliangzhao/balancer/pkg/apis/balancer/v1alpha1"
	controllerbalancer "github.com/hliangzhao/balancer/pkg/controllers/balancer"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"time"
)

func MakeBasicBalancer(namespace, name string, versions []string, weights []int32) *hliangzhaov1alpha1.Balancer {
	var backends []hliangzhaov1alpha1.BackendSpec
	{
	}
	for idx := range versions {
		backends = append(backends, hliangzhaov1alpha1.BackendSpec{
			Name:     versions[idx],
			Weight:   weights[idx],
			Selector: map[string]string{"version": versions[idx]},
		})
	}
	return &hliangzhaov1alpha1.Balancer{
		ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: namespace},
		Spec: hliangzhaov1alpha1.BalancerSpec{
			Backends: backends,
			Selector: map[string]string{"app": name},
			Ports: []hliangzhaov1alpha1.BalancerPort{
				{
					Name:     "http",
					Protocol: hliangzhaov1alpha1.TCP,
					Port:     80,
				},
			},
		},
	}
}

func (f *Framework) CreateBalancer(namespace string,
	balancer *hliangzhaov1alpha1.Balancer) (*hliangzhaov1alpha1.Balancer, error) {
	result, err := f.ExposerClientV1alpha1.Balancers(namespace).Create(context.Background(),
		balancer, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (f *Framework) WaitForBalancerReady(balancer *hliangzhaov1alpha1.Balancer, timeout time.Duration) error {
	var pollErr error
	err := wait.Poll(2*time.Second, timeout, func() (bool, error) {
		actualBalancer, pollErr := f.ExposerClientV1alpha1.Balancers(balancer.Namespace).Get(context.Background(),
			balancer.Name, metav1.GetOptions{})
		if pollErr != nil {
			return false, nil
		}

		// balancer is ready only when the following are ready
		deploymentName := controllerbalancer.DeploymentName(balancer)
		if err := f.WaitForDeploymentCreated(balancer.Namespace, deploymentName, timeout); err != nil {
			return false, err
		}
		if err := f.WaitForServiceCreated(balancer.Namespace, balancer.Name, timeout); err != nil {
			return false, err
		}
		for _, backend := range balancer.Spec.Backends {
			if err := f.WaitForServiceCreated(
				balancer.Namespace,
				fmt.Sprintf("%s-%s-backend", balancer.Name, backend.Name),
				timeout,
			); err != nil {
				return false, err
			}
		}
		if actualBalancer.Status.ActiveBackendsNum != int32(len(balancer.Spec.Backends)) ||
			actualBalancer.Status.ObsoleteBackendsNum != 0 {
			return false, nil
		}
		// TODO: check the balancer expected backends with current backends

		return true, nil
	})
	return errors.Wrapf(pollErr, "waiting for Balancer %s/%s: %v", balancer.Namespace, balancer.Name, err)
}

func (f *Framework) CreateBalancerAndWaitUntilReady(namespace string,
	balancer *hliangzhaov1alpha1.Balancer) (*hliangzhaov1alpha1.Balancer, error) {

	result, err := f.CreateBalancer(namespace, balancer)
	if err != nil {
		return nil, err
	}
	if err = f.WaitForBalancerReady(result, 15*time.Second); err != nil {
		return nil, fmt.Errorf("waiting for Balancer instances timed out %s: %v", balancer.Name, err)
	}
	return result, nil
}

func (f *Framework) UpdateBalancer(namespace string,
	balancer *hliangzhaov1alpha1.Balancer) (*hliangzhaov1alpha1.Balancer, error) {

	result, err := f.ExposerClientV1alpha1.Balancers(namespace).Update(context.Background(), balancer, metav1.UpdateOptions{})
	if err != nil {
		return nil, fmt.Errorf("updating Balancer instance failed (%s): %v", balancer.Name, err)
	}
	return result, nil
}

func (f *Framework) UpdateBalancerAndWaitUntilReady(namespace string,
	balancer *hliangzhaov1alpha1.Balancer) (*hliangzhaov1alpha1.Balancer, error) {

	result, err := f.UpdateBalancer(namespace, balancer)
	if err != nil {
		return nil, err
	}
	if err := f.WaitForBalancerReady(result, 15*time.Second); err != nil {
		return nil, fmt.Errorf("waiting for Balancer instance timed out (%s): %v", balancer.Name, err)
	}
	return result, nil
}
