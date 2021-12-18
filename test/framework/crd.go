package framework

import (
	`context`
	`github.com/pkg/errors`
	apiextensionsv1beta1 `k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1`
	`k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset`
	apierrors `k8s.io/apimachinery/pkg/api/errors`
	metav1 `k8s.io/apimachinery/pkg/apis/meta/v1`
	`k8s.io/apimachinery/pkg/runtime`
	`k8s.io/apimachinery/pkg/util/wait`
	`k8s.io/apimachinery/pkg/util/yaml`
	`net/http`
	`time`
)

func CreateCRD(client clientset.Clientset, namespace string, crd *apiextensionsv1beta1.CustomResourceDefinition) error {
	crd.Namespace = namespace
	// TODO: it seems like this operation will update crd too
	crd, err := client.ApiextensionsV1beta1().CustomResourceDefinitions().Create(context.Background(), crd, metav1.CreateOptions{})
	if err != nil {
		return err
	}
	return nil
}

func WaitForCRDReady(listFunc func(opts metav1.ListOptions) (runtime.Object, error)) error {
	err := wait.Poll(3*time.Second, 10*time.Minute, func() (done bool, err error) {
		_, err = listFunc(metav1.ListOptions{})
		if err != nil {
			if se, ok := err.(*apierrors.StatusError); ok {
				if se.Status().Code == http.StatusNotFound {
					done = false
					err = nil
					return
				}
			}
			done = false
			err = errors.Wrap(err, "failed to list CRD")
		}
		done = true
		err = nil
		return
	})

	// TODO: if execute to here, it means timeout?
	return errors.Wrap(err, "timeout waiting for CRD")
}

func ParseCRDYaml(relativeFilePath string) (*apiextensionsv1beta1.CustomResourceDefinition, error) {
	manifest, err := GetFileDescriptor(relativeFilePath)
	if err != nil {
		return nil, err
	}
	crd := &apiextensionsv1beta1.CustomResourceDefinition{}
	if err = yaml.NewYAMLOrJSONDecoder(manifest, 100).Decode(crd); err != nil {
		return nil, err
	}
	return crd, nil
}