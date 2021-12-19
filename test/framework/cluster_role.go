package framework

import (
	"context"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/kubernetes"
)

func CreateOrUpdateClusterRole(client kubernetes.Interface, relativeFilePath string) error {
	clusterRole, err := ParseClusterRoleYaml(relativeFilePath)
	if err != nil {
		return err
	}

	_, err = client.RbacV1().ClusterRoles().Get(context.Background(), clusterRole.Name, metav1.GetOptions{})
	if err != nil {
		// non-exist, create
		if _, err := client.RbacV1().ClusterRoles().Create(context.Background(),
			clusterRole, metav1.CreateOptions{}); err != nil {
			return err
		}
	} else {
		// exist, update
		if _, err := client.RbacV1().ClusterRoles().Update(context.Background(),
			clusterRole, metav1.UpdateOptions{}); err != nil {
			return err
		}
	}
	return nil
}

func DeleteClusterRole(client kubernetes.Interface, relativeFilePath string) error {
	clusterRole, err := ParseClusterRoleYaml(relativeFilePath)
	if err != nil {
		return err
	}
	return client.RbacV1().ClusterRoles().Delete(context.Background(), clusterRole.Name, metav1.DeleteOptions{})
}

// ParseClusterRoleYaml creates a ClusterRole resource with given yaml file path.
func ParseClusterRoleYaml(relativeFilePath string) (*rbacv1.ClusterRole, error) {
	manifest, err := GetFileDescriptor(relativeFilePath)
	if err != nil {
		return nil, err
	}
	clusterRole := &rbacv1.ClusterRole{}
	if err = yaml.NewYAMLOrJSONDecoder(manifest, 100).Decode(clusterRole); err != nil {
		return nil, err
	}
	return clusterRole, nil
}
