package framework

import `k8s.io/client-go/kubernetes`

func CreateOrUpdateClusterRoleBinding(client kubernetes.Interface, namespace, relativeFilePath string) (finalizeFunc, error) {
	return nil, nil
}