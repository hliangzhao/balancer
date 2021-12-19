package framework

import (
	"context"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/kubernetes"
)

func CreateServiceAccount(client kubernetes.Interface, namespace, relativeFilePath string) (FinalizeFunc, error) {
	finalizer := func() error {
		return DeleteServiceAccount(client, namespace, relativeFilePath)
	}

	serviceAccount, err := ParseServiceAccountYaml(relativeFilePath)
	if err != nil {
		return finalizer, err
	}
	serviceAccount.Namespace = namespace
	_, err = client.CoreV1().ServiceAccounts(namespace).Create(context.Background(),
		serviceAccount, metav1.CreateOptions{})
	if err != nil {
		return finalizer, err
	}
	return finalizer, nil
}

func DeleteServiceAccount(client kubernetes.Interface, namespace, relativeFilePath string) error {
	serviceAccount, err := ParseServiceAccountYaml(relativeFilePath)
	if err != nil {
		return err
	}
	return client.CoreV1().ServiceAccounts(namespace).Delete(context.Background(),
		serviceAccount.Name, metav1.DeleteOptions{})
}

func ParseServiceAccountYaml(relativeFilePath string) (*corev1.ServiceAccount, error) {
	manifest, err := GetFileDescriptor(relativeFilePath)
	if err != nil {
		return nil, err
	}
	serviceAccount := &corev1.ServiceAccount{}
	if err = yaml.NewYAMLOrJSONDecoder(manifest, 100).Decode(serviceAccount); err != nil {
		return nil, err
	}
	return serviceAccount, nil
}
