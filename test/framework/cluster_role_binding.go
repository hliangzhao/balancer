package framework

import (
	"context"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/kubernetes"
)

func CreateOrUpdateClusterRoleBinding(client kubernetes.Interface, namespace, relativeFilePath string) (FinalizeFunc, error) {
	finalizer := func() error {
		return DeleteClusterRoleBinding(client, relativeFilePath)
	}
	clusterRoleBinding, err := ParseClusterRoleBindingYaml(relativeFilePath)
	if err != nil {
		return finalizer, err
	}
	clusterRoleBinding.Subjects[0].Namespace = namespace

	_, err = client.RbacV1().ClusterRoleBindings().Get(context.Background(), clusterRoleBinding.Name, metav1.GetOptions{})
	if err != nil {
		// non-exist, create
		if _, err = client.RbacV1().ClusterRoleBindings().Create(context.Background(),
			clusterRoleBinding, metav1.CreateOptions{}); err != nil {
			return finalizer, err
		}
	} else {
		// exist, update
		if _, err = client.RbacV1().ClusterRoleBindings().Update(context.Background(),
			clusterRoleBinding, metav1.UpdateOptions{}); err != nil {
			return finalizer, err
		}
	}
	// TODO: why not nil?
	return finalizer, err
}

func DeleteClusterRoleBinding(client kubernetes.Interface, relativeFilePath string) error {
	clusterRoleBinding, err := ParseClusterRoleYaml(relativeFilePath)
	if err != nil {
		return err
	}
	return client.RbacV1().ClusterRoles().Delete(context.Background(), clusterRoleBinding.Name, metav1.DeleteOptions{})
}

// ParseClusterRoleBindingYaml creates a ClusterRoleBinding resource with given yaml file path.
func ParseClusterRoleBindingYaml(relativeFilePath string) (*rbacv1.ClusterRoleBinding, error) {
	manifest, err := GetFileDescriptor(relativeFilePath)
	if err != nil {
		return nil, err
	}
	clusterRoleBinding := &rbacv1.ClusterRoleBinding{}
	if err = yaml.NewYAMLOrJSONDecoder(manifest, 100).Decode(clusterRoleBinding); err != nil {
		return nil, err
	}
	return clusterRoleBinding, nil
}
