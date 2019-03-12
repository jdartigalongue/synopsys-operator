/*
Copyright The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by client-gen. DO NOT EDIT.

package v1

import (
	v1 "github.com/blackducksoftware/synopsys-operator/pkg/api/rgp/v1"
	scheme "github.com/blackducksoftware/synopsys-operator/pkg/rgp/client/clientset/versioned/scheme"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// RgpsGetter has a method to return a RgpInterface.
// A group's client should implement this interface.
type RgpsGetter interface {
	Rgps(namespace string) RgpInterface
}

// RgpInterface has methods to work with Rgp resources.
type RgpInterface interface {
	Create(*v1.Rgp) (*v1.Rgp, error)
	Update(*v1.Rgp) (*v1.Rgp, error)
	Delete(name string, options *metav1.DeleteOptions) error
	DeleteCollection(options *metav1.DeleteOptions, listOptions metav1.ListOptions) error
	Get(name string, options metav1.GetOptions) (*v1.Rgp, error)
	List(opts metav1.ListOptions) (*v1.RgpList, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.Rgp, err error)
	RgpExpansion
}

// rgps implements RgpInterface
type rgps struct {
	client rest.Interface
	ns     string
}

// newRgps returns a Rgps
func newRgps(c *SynopsysV1Client, namespace string) *rgps {
	return &rgps{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the rgp, and returns the corresponding rgp object, and an error if there is any.
func (c *rgps) Get(name string, options metav1.GetOptions) (result *v1.Rgp, err error) {
	result = &v1.Rgp{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("rgps").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of Rgps that match those selectors.
func (c *rgps) List(opts metav1.ListOptions) (result *v1.RgpList, err error) {
	result = &v1.RgpList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("rgps").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested rgps.
func (c *rgps) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("rgps").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}

// Create takes the representation of a rgp and creates it.  Returns the server's representation of the rgp, and an error, if there is any.
func (c *rgps) Create(rgp *v1.Rgp) (result *v1.Rgp, err error) {
	result = &v1.Rgp{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("rgps").
		Body(rgp).
		Do().
		Into(result)
	return
}

// Update takes the representation of a rgp and updates it. Returns the server's representation of the rgp, and an error, if there is any.
func (c *rgps) Update(rgp *v1.Rgp) (result *v1.Rgp, err error) {
	result = &v1.Rgp{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("rgps").
		Name(rgp.Name).
		Body(rgp).
		Do().
		Into(result)
	return
}

// Delete takes name of the rgp and deletes it. Returns an error if one occurs.
func (c *rgps) Delete(name string, options *metav1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("rgps").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *rgps) DeleteCollection(options *metav1.DeleteOptions, listOptions metav1.ListOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("rgps").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched rgp.
func (c *rgps) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.Rgp, err error) {
	result = &v1.Rgp{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("rgps").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
