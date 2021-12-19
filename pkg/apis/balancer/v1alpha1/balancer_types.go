package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

type Protocol string
type Port int32

const (
	TCP Protocol = "TCP"
	UDP Protocol = "UDP"
)

// BalancerPort contains the endpoints and exposed ports.
type BalancerPort struct {
	// +required
	Name string `json:"name,omitempty"` // name of the endpoints

	// +optional
	Protocol Protocol `json:"protocol,omitempty"` // default is TCP

	// the port that will be exposed by the balancer
	Port Port `json:"port"`

	// the port that used by the container
	// +optional
	TargetPort intstr.IntOrString `json:"targetPort,omitempty"`
}

// BackendSpec defines the desired status of endpoints of Balancer.
type BackendSpec struct {
	// +kubebuilder:validation:MinLength=1
	Name string `json:"name"`

	// +kubebuilder:validation:Minimum=1
	Weight int32 `json:"weight"`

	Selector map[string]string `json:"selector,omitempty"`
}

// BalancerSpec defines the desired status of Balancer.
// +k8s:openapi-gen=true
type BalancerSpec struct {
	// +kubebuilder:validation:MinItems=1
	Backends []BackendSpec `json:"backends"`

	Selector map[string]string `json:"selector,omitempty"`

	Ports []BalancerPort `json:"ports"`
}

// BalancerStatus defines the actual status of Balancer.
// TODO: could add more fields (new feature)
// +k8s:openapi-gen=true
type BalancerStatus struct {
	// +optional
	ActiveBackendsNum int32 `json:"activeBackendsNum,omitempty"`

	// +optional
	ObsoleteBackendsNum int32 `json:"obsoleteBackendsNum,omitempty"`
}

// Balancer is the Schema for the balancer api.
// Example:
// ==============================
// apiVersion: hliangzhao.io/v1alpha1
// kind: Balancer
// metadata:
//  name: example-balancer
// spec:
//  ports:
//    # all the ports will be exposed with a Service
//    - name: http
//      protocol: TCP
//      port: 80
//      targetPort: 5678
//  selector:
//    app: test
//  backends:
//    # each backend is a Pod that can handle the input load
//    # a service will also be created for each backend
//    - name: v1
//      weight: 90
//      selector:
//        version: v1
//    - name: v2
//      weight: 9
//      selector:
//        version: v2
// ==============================
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +genclient
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
type Balancer struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   BalancerSpec   `json:"spec,omitempty"`
	Status BalancerStatus `json:"status,omitempty"`
}

// BalancerList contains a list of Balancer.
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type BalancerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []Balancer `json:"items"`
}
