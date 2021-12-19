package framework

import (
	"context"
	appv1 "k8s.io/api/apps/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/kubernetes"
	"time"
)

func CreateDeployment(client kubernetes.Interface, namespace string, deployment *appv1.Deployment) error {
	deployment.Namespace = namespace
	_, err := client.AppsV1().Deployments(namespace).Create(context.Background(), deployment, metav1.CreateOptions{})
	if err != nil {
		return err
	}
	return nil
}

func DeleteDeployment(client kubernetes.Interface, namespace, name string) error {
	deployment, err := client.AppsV1().Deployments(namespace).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		return err
	}

	// TODO: why not Delete directly?
	zero := int32(0)
	deployment.Spec.Replicas = &zero
	deployment, err = client.AppsV1().Deployments(namespace).Update(context.Background(), deployment, metav1.UpdateOptions{})
	if err != nil {
		return err
	}
	return client.AppsV1().Deployments(namespace).Delete(context.Background(), deployment.Name, metav1.DeleteOptions{})
}

// WaitUntilDeploymentGone returns true only if the deployment is deleted.
func WaitUntilDeploymentGone(client kubernetes.Interface, namespace, name string, timeout time.Duration) error {
	return wait.Poll(time.Second, timeout, func() (bool, error) {
		_, err := client.AppsV1().Deployments(namespace).Get(context.Background(), name, metav1.GetOptions{})
		if err != nil {
			if apierrors.IsNotFound(err) {
				return true, nil
			}
			return false, err
		}
		return false, nil
	})
}

// WaitForDeploymentCreated returns true only if the deployment is created.
func (f *Framework) WaitForDeploymentCreated(namespace, name string, timeout time.Duration) error {
	return wait.Poll(time.Second, timeout, func() (bool, error) {
		_, err := f.KubeClient.AppsV1().Deployments(namespace).Get(context.Background(), name, metav1.GetOptions{})
		if err != nil {
			if apierrors.IsNotFound(err) {
				return false, nil
			}
			return false, err
		}
		return true, nil
	})
}

func ParseDeploymentYaml(relativeFilePath string) (*appv1.Deployment, error) {
	manifest, err := GetFileDescriptor(relativeFilePath)
	if err != nil {
		return nil, err
	}
	deployment := &appv1.Deployment{}
	if err = yaml.NewYAMLOrJSONDecoder(manifest, 100).Decode(deployment); err != nil {
		return nil, err
	}
	return deployment, nil
}
