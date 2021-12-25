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
	exposerv1alpha1 "github.com/hliangzhao/balancer/pkg/apis/balancer/v1alpha1"
)

func NewPodLabels(balancer *exposerv1alpha1.Balancer) map[string]string {
	return map[string]string{
		exposerv1alpha1.BalancerKey: balancer.Name,
	}
}

func NewServiceLabels(balancer *exposerv1alpha1.Balancer) map[string]string {
	return map[string]string{
		exposerv1alpha1.BalancerKey: balancer.Name,
	}
}
