// Code generated by client-gen. DO NOT EDIT.

package v1alpha1

import (
	"context"
	"time"

	v1alpha1 "github.com/hliangzhao/balancer/pkg/apis/balancer/v1alpha1"
	scheme "github.com/hliangzhao/balancer/pkg/client/clientset/versioned/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// BalancersGetter has a method to return a BalancerInterface.
// A group's client should implement this interface.
type BalancersGetter interface {
	Balancers(namespace string) BalancerInterface
}

// BalancerInterface has methods to work with Balancer resources.
type BalancerInterface interface {
	Create(ctx context.Context, balancer *v1alpha1.Balancer, opts v1.CreateOptions) (*v1alpha1.Balancer, error)
	Update(ctx context.Context, balancer *v1alpha1.Balancer, opts v1.UpdateOptions) (*v1alpha1.Balancer, error)
	UpdateStatus(ctx context.Context, balancer *v1alpha1.Balancer, opts v1.UpdateOptions) (*v1alpha1.Balancer, error)
	Delete(ctx context.Context, name string, opts v1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error
	Get(ctx context.Context, name string, opts v1.GetOptions) (*v1alpha1.Balancer, error)
	List(ctx context.Context, opts v1.ListOptions) (*v1alpha1.BalancerList, error)
	Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.Balancer, err error)
	BalancerExpansion
}

// balancers implements BalancerInterface
type balancers struct {
	client rest.Interface
	ns     string
}

// newBalancers returns a Balancers
func newBalancers(c *HliangzhaoV1alpha1Client, namespace string) *balancers {
	return &balancers{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the balancer, and returns the corresponding balancer object, and an error if there is any.
func (c *balancers) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.Balancer, err error) {
	result = &v1alpha1.Balancer{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("balancers").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of Balancers that match those selectors.
func (c *balancers) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.BalancerList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1alpha1.BalancerList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("balancers").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested balancers.
func (c *balancers) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("balancers").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a balancer and creates it.  Returns the server's representation of the balancer, and an error, if there is any.
func (c *balancers) Create(ctx context.Context, balancer *v1alpha1.Balancer, opts v1.CreateOptions) (result *v1alpha1.Balancer, err error) {
	result = &v1alpha1.Balancer{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("balancers").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(balancer).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a balancer and updates it. Returns the server's representation of the balancer, and an error, if there is any.
func (c *balancers) Update(ctx context.Context, balancer *v1alpha1.Balancer, opts v1.UpdateOptions) (result *v1alpha1.Balancer, err error) {
	result = &v1alpha1.Balancer{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("balancers").
		Name(balancer.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(balancer).
		Do(ctx).
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *balancers) UpdateStatus(ctx context.Context, balancer *v1alpha1.Balancer, opts v1.UpdateOptions) (result *v1alpha1.Balancer, err error) {
	result = &v1alpha1.Balancer{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("balancers").
		Name(balancer.Name).
		SubResource("status").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(balancer).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the balancer and deletes it. Returns an error if one occurs.
func (c *balancers) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("balancers").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *balancers) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("balancers").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched balancer.
func (c *balancers) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.Balancer, err error) {
	result = &v1alpha1.Balancer{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("balancers").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}
