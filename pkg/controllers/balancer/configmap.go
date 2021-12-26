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
	"fmt"
	exposerv1alpha1 "github.com/hliangzhao/balancer/pkg/apis/balancer/v1alpha1"
	"github.com/hliangzhao/balancer/pkg/controllers/balancer/nginx"
	"hash/fnv"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	randutil "k8s.io/apimachinery/pkg/util/rand"
	hashutil "k8s.io/kubernetes/pkg/util/hash"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// NewConfigMap creates a new configmap for the input Balancer instance.
func NewConfigMap(balancer *exposerv1alpha1.Balancer) (*corev1.ConfigMap, error) {
	return &corev1.ConfigMap{
		ObjectMeta: v1.ObjectMeta{
			Name:      ConfigMapName(balancer),
			Namespace: balancer.Namespace,
		},
		Data: map[string]string{
			"nginx.conf": nginx.NewConfig(balancer),
		},
	}, nil
}

// syncConfigMap sync the configmap that created by the deployment of Balancer.
func (r *ReconcilerBalancer) syncConfigMap(balancer *exposerv1alpha1.Balancer) (*corev1.ConfigMap, error) {
	cm, err := NewConfigMap(balancer)
	if err != nil {
		return nil, err
	}

	// set balancer as the controller owner-reference of cm
	if err := controllerutil.SetControllerReference(balancer, cm, r.scheme); err != nil {
		return nil, err
	}

	foundCm := &corev1.ConfigMap{}
	err = r.client.Get(context.Background(), types.NamespacedName{Namespace: cm.Namespace, Name: cm.Name}, foundCm)
	if err != nil && errors.IsNotFound(err) {
		// corresponding cm not foundCm in the cluster, create it with the newest cm
		if err = r.client.Create(context.Background(), cm); err != nil {
			return nil, err
		}
		log.Info("Sync ConfigMap", cm.Name, "created")
		return cm, nil
	} else if err != nil {
		return nil, err
	}

	// corresponding cm foundCm, update it with the newest cm
	foundCm.Data = cm.Data
	if err = r.client.Update(context.Background(), foundCm); err != nil {
		return nil, err
	}
	log.Info("Sync ConfigMap", foundCm.Name, "updated")
	return cm, nil
}

func ConfigMapName(balancer *exposerv1alpha1.Balancer) string {
	return balancer.Name + "-proxy-configmap"
}

func ConfigMapHash(cm *corev1.ConfigMap) string {
	hasher := fnv.New32a()
	hashutil.DeepHashObject(hasher, cm)
	return randutil.SafeEncodeString(fmt.Sprint(hasher.Sum32()))
}
