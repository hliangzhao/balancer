package framework

import (
	"context"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"testing"
)

func CreateNamespace(client kubernetes.Interface, name string) (*corev1.Namespace, error) {
	namespace, err := client.CoreV1().Namespaces().Create(context.Background(), &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{Name: name},
	}, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}
	return namespace, nil
}

func DeleteNamespace(client kubernetes.Interface, name string) error {
	return client.CoreV1().Namespaces().Delete(context.Background(), name, metav1.DeleteOptions{})
}

func AddLabelsToNamespace(client kubernetes.Interface, name string, additionalLabels map[string]string) error {
	namespace, err := client.CoreV1().Namespaces().Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		return err
	}
	if namespace.Labels == nil {
		namespace.Labels = map[string]string{}
	}

	// add
	for k, v := range additionalLabels {
		namespace.Labels[k] = v
	}

	// update to api-server
	_, err = client.CoreV1().Namespaces().Update(context.Background(), namespace, metav1.UpdateOptions{})
	if err != nil {
		return err
	}
	return nil
}

func (ctx *TestContext) CreateNamespace(t *testing.T, client kubernetes.Interface) string {
	name := ctx.GetObjId()
	if _, err := CreateNamespace(client, name); err != nil {
		t.Fatal(err)
	}
	finalizer := func() error {
		return DeleteNamespace(client, name)
	}
	ctx.AddFinalizerFunc(finalizer)
	return name
}
