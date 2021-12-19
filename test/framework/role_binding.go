package framework

import (
	"context"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/kubernetes"
)

func CreateRoleBinding(client kubernetes.Interface, namespace, relativeFilePath string) (FinalizeFunc, error) {
	finalizer := func() error {
		return DeleteRoleBinding(client, namespace, relativeFilePath)
	}
	roleBinding, err := ParseRoleBindingYaml(relativeFilePath)
	if err != nil {
		return finalizer, err
	}
	_, err = client.RbacV1().RoleBindings(namespace).Create(context.Background(), roleBinding, metav1.CreateOptions{})
	return finalizer, err
}

func DeleteRoleBinding(client kubernetes.Interface, namespace string, relativeFilePath string) error {
	roleBinding, err := ParseRoleBindingYaml(relativeFilePath)
	if err != nil {
		return err
	}
	return client.RbacV1().RoleBindings(namespace).Delete(context.Background(), roleBinding.Name, metav1.DeleteOptions{})
}

func ParseRoleBindingYaml(relativeFilePath string) (*rbacv1.RoleBinding, error) {
	manifest, err := GetFileDescriptor(relativeFilePath)
	if err != nil {
		return nil, err
	}

	roleBinding := &rbacv1.RoleBinding{}
	if err = yaml.NewYAMLOrJSONDecoder(manifest, 100).Decode(roleBinding); err != nil {
		return nil, err
	}
	return roleBinding, nil
}
