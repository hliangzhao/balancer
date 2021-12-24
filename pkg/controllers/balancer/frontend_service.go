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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// syncFrontendService sync the front-end service that created by balancer.
func (r *ReconcilerBalancer) syncFrontendService(balancer *balancerv1alpha1.Balancer) error {
	svc, err := NewFrontendService(balancer)
	if err != nil {
		return err
	}

	// set balancer as the controller owner-reference of svc
	if err := controllerutil.SetControllerReference(balancer, svc, r.scheme); err != nil {
		return err
	}

	foundSvc := &corev1.Service{}
	err = r.client.Get(context.Background(), types.NamespacedName{Namespace: svc.Namespace, Name: svc.Name}, foundSvc)
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
	foundSvc.Spec.Ports = svc.Spec.Ports
	foundSvc.Spec.Selector = svc.Spec.Selector
	if err = r.client.Update(context.Background(), foundSvc); err != nil {
		return err
	}
	return nil
}

// NewFrontendService creates a new front-end Service for handling all requests incoming.
// All the incoming requests will be forwarded to backend services by the nginx instance.
func NewFrontendService(balancer *balancerv1alpha1.Balancer) (*corev1.Service, error) {
	var balancerPorts []corev1.ServicePort
	for _, port := range balancer.Spec.Ports {
		balancerPorts = append(balancerPorts, corev1.ServicePort{
			// TODO: the mapping from Port to TargetPort has bug! After successfully built, fix this
			Name:     port.Name,
			Protocol: corev1.Protocol(port.Protocol),
			Port:     int32(port.Port),
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
