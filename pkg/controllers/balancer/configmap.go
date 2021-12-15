package balancer

import (
	`context`
	`fmt`
	balancerv1alpha1 `github.com/hliangzhao/balancer/pkg/apis/balancer/v1alpha1`
	`github.com/hliangzhao/balancer/pkg/controllers/balancer/nginx`
	`hash/fnv`
	corev1 `k8s.io/api/core/v1`
	`k8s.io/apimachinery/pkg/api/errors`
	v1 `k8s.io/apimachinery/pkg/apis/meta/v1`
	`k8s.io/apimachinery/pkg/types`
	randutil `k8s.io/apimachinery/pkg/util/rand`
	hashutil `k8s.io/kubernetes/pkg/util/hash`
	`sigs.k8s.io/controller-runtime/pkg/controller/controllerutil`
)

// NewConfigMap creates a new configmap for the input Balancer instance.
func NewConfigMap(balancer *balancerv1alpha1.Balancer) (*corev1.ConfigMap, error) {
	return &corev1.ConfigMap{
		ObjectMeta: v1.ObjectMeta{
			Name: ConfigMapName(balancer),
			Namespace: balancer.Namespace,
		},
		Data: map[string]string{
			"nginx.conf": nginx.NewConfig(balancer),
		},
	},nil
}

// syncConfigMap sync the configmap that created by the deployment of Balancer.
func (r *ReconcileBalancer) syncConfigMap(balancer *balancerv1alpha1.Balancer) (*corev1.ConfigMap, error) {
	cm, err := NewConfigMap(balancer)
	if err != nil {
		return nil, err
	}

	// set balancer as the owner of cm
	if err := controllerutil.SetOwnerReference(balancer, cm, r.scheme); err != nil {
		return nil, err
	}

	found := &corev1.ConfigMap{}
	err = r.client.Get(context.Background(), types.NamespacedName{Namespace: cm.Namespace, Name: cm.Name}, found)
	if err != nil && errors.IsNotFound(err) {
		// corresponding cm not found in the cluster, create it with the newest cm
		if err = r.client.Create(context.Background(), cm); err != nil {
			return nil, err
		}
		return cm, nil
	} else if err != nil {
		return nil, err
	}

	// corresponding cm found, update it with the newest cm
	found.Data = cm.Data
	if err = r.client.Update(context.Background(), found); err != nil {
		return nil, err
	}
	return cm, nil
}

func ConfigMapName(balancer *balancerv1alpha1.Balancer) string {
	return balancer.Name + "-proxy-configmap"
}

func ConfigMapHash(cm *corev1.ConfigMap) string {
	hasher := fnv.New32a()
	hashutil.DeepHashObject(hasher, cm)
	return randutil.SafeEncodeString(fmt.Sprint(hasher.Sum32()))
}
