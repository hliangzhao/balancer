package v1alpha1

const (
	// ConfigMapHashKey is the key of the annotation which is used by the Balancer.
	// the Balancer is actually a Nginx instance, and the configmap is acted as the `nginx.conf`.
	ConfigMapHashKey = "balancer.hliangzhao.io/configmap-hash"
)
