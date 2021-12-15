package balancer

import (
	`context`
	`fmt`
	balancerv1alpha1 `github.com/hliangzhao/balancer/pkg/apis/balancer/v1alpha1`
	corev1 `k8s.io/api/core/v1`
	`k8s.io/apimachinery/pkg/api/errors`
	v1 `k8s.io/apimachinery/pkg/apis/meta/v1`
	`k8s.io/apimachinery/pkg/types`
	`reflect`
	`sigs.k8s.io/controller-runtime/pkg/client`
	`sigs.k8s.io/controller-runtime/pkg/controller/controllerutil`
	`sync`
)

// syncBalancerStatus sync the status of balancer.
func (r *ReconcileBalancer) syncBalancerStatus(balancer *balancerv1alpha1.Balancer) error {
	// get current services
	var svcList corev1.ServiceList
	if err := r.client.List(context.Background(), &svcList, client.MatchingLabels(NewServiceLabels(balancer))); err != nil {
		return err
	}

	_, servicesToDelete, activeServices := groupServers(balancer, svcList.Items)

	expectedStatus := balancerv1alpha1.BalancerStatus{
		ActiveBackendsNum:   int32(len(activeServices)),
		ObsoleteBackendsNum: int32(len(servicesToDelete)),
	}
	// nothing to do, return directly
	if reflect.DeepEqual(balancer.Status, expectedStatus) {
		return nil
	}

	// status updating is required
	balancer.Status = expectedStatus
	return r.client.Status().Update(context.Background(), balancer)
}

// syncServers creates services to be created and deletes services to be deleted.
func (r *ReconcileBalancer) syncServers(balancer *balancerv1alpha1.Balancer) error {
	// get current services
	var svcList corev1.ServiceList
	if err := r.client.List(context.Background(), &svcList, client.MatchingLabels(NewServiceLabels(balancer))); err != nil {
		return err
	}

	servicesToCreate, servicesToDelete, _ := groupServers(balancer, svcList.Items)

	wg := sync.WaitGroup{}

	// start coroutines to delete services-to-be-deleted
	deleteErrCh := make(chan error, len(servicesToDelete))
	wg.Add(len(servicesToDelete))
	for _, svcToDelete := range servicesToDelete {
		go func(svc *corev1.Service) {
			defer wg.Done()
			if err := r.client.Delete(context.Background(), svc); err != nil {
				deleteErrCh <- err
			}
		}(&svcToDelete)
	}
	wg.Wait()

	// start coroutines to create services-to-be-created
	createErrCh := make(chan error, len(servicesToCreate))
	wg.Add(len(servicesToCreate))
	for _, svcToCreate := range servicesToCreate {
		go func(svc *corev1.Service) {
			defer wg.Done()

			// set balancer as the owner of svc
			if err := controllerutil.SetOwnerReference(balancer, svc, r.scheme); err != nil {
				createErrCh <- err
				return
			}

			// create or update
			found := &corev1.Service{}
			err := r.client.Get(context.Background(), types.NamespacedName{Namespace: svc.Namespace, Name: svc.Name}, found)
			if err != nil && errors.IsNotFound(err) {
				if err = r.client.Create(context.Background(), svc); err != nil {
					createErrCh <- err
				}
				return
			} else if err != nil {
				createErrCh <- err
				return
			}

			found.Spec.Ports = svc.Spec.Ports
			found.Spec.Selector = svc.Spec.Selector
			if err = r.client.Update(context.Background(), svc); err != nil {
				createErrCh <- err
				return
			}
		}(&svcToCreate)
	}
	wg.Wait()

	// handle error if happened
	select {
	case err := <-deleteErrCh:
		return err
	case err := <-createErrCh:
		return err
	default:
		return r.syncBalancerStatus(balancer)
	}
}

// groupServers gets services to created from current balancer and divided current services into active and to-be-deleted.
// TODO: servicesToCreate is not necessary?
func groupServers(balancer *balancerv1alpha1.Balancer, services []corev1.Service) (servicesToCreate []corev1.Service,
	servicesToDelete []corev1.Service, activeServices []corev1.Service) {

	var balancerPorts []corev1.ServicePort
	for _, port := range balancer.Spec.Ports {
		balancerPorts = append(balancerPorts, corev1.ServicePort{
			Name:       port.Name,
			Protocol:   corev1.Protocol(port.Protocol),
			Port:       int32(port.Port),
			TargetPort: port.TargetPort,
		})
	}

	// create service for each backend in Balancer
	for _, backend := range balancer.Spec.Backends {
		selector := map[string]string{}
		for k, v := range balancer.Spec.Selector {
			selector[k] = v
		}
		for k, v := range backend.Selector {
			selector[k] = v
		}
		servicesToCreate = append(servicesToCreate, corev1.Service{
			ObjectMeta: v1.ObjectMeta{
				Name:      fmt.Sprintf("%s-%s-backend", balancer.Name, backend.Name),
				Namespace: balancer.Namespace,
				Labels:    NewServiceLabels(balancer),
			},
			Spec: corev1.ServiceSpec{
				Selector: selector,
				Type:     corev1.ServiceTypeClusterIP,
				Ports:    balancerPorts,
			},
		})
	}

	for _, svc := range services {
		// judge svc is active or need-to-be deleted
		existActiveSvc := false
		for _, svcToCreate := range servicesToCreate {
			if svcToCreate.Name == svc.Name && svcToCreate.Namespace == svc.Namespace {
				activeServices = append(activeServices, svc)
				existActiveSvc = true
				break
			}
		}
		if !existActiveSvc {
			servicesToDelete = append(servicesToDelete, svc)
		}
	}
	return
}
