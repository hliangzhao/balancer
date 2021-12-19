package framework

import (
	hliangzhaoclientv1alpha1 "github.com/hliangzhao/balancer/pkg/client/clientset/versioned/typed/balancer/v1alpha1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"net/http"
	"testing"
	"time"
)

// Framework is a test framework
type Framework struct {
	// client for operate k8s core resources
	KubeClient kubernetes.Interface
	// client for operate hliangzhao.io resources
	HliangzhaoClientV1alpha1 hliangzhaoclientv1alpha1.HliangzhaoV1alpha1Interface
	// client for operate extension resources
	ApiextensionsClientV1 apiextensionsv1.Interface
	HttpClient            *http.Client
	MasterHost            string
	DefaultTimeout        time.Duration
}

func NewFramework(kubeconfig string) (*Framework, error) {
	// masterUrl and kubeconfigPath, at least one is provided
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, err
	}
	// create clientset from the given kubeconfig
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	httpClient := client.CoreV1().RESTClient().(*rest.RESTClient).Client
	hliangzhaoClient, err := hliangzhaoclientv1alpha1.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	apiextensionsClient, err := apiextensionsv1.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return &Framework{
		KubeClient:               client,
		HliangzhaoClientV1alpha1: hliangzhaoClient,
		ApiextensionsClientV1:    apiextensionsClient,
		HttpClient:               httpClient,
		MasterHost:               config.Host,
		DefaultTimeout:           time.Minute,
	}, nil
}

func (f *Framework) CreateBalancerOperator(namespace string, operatorImage string, namespacesToWatch []string) error {
	// TODO
	return nil
}

func (ctx *TestContext) SetupBalancerRBAC(t *testing.T, namespace string, kubeClient kubernetes.Interface) {
	// TODO
}
