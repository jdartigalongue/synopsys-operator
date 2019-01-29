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

package fake

import (
	opssight_v1 "github.com/blackducksoftware/synopsys-operator/pkg/api/opssight/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeOpsSights implements OpsSightInterface
type FakeOpsSights struct {
	Fake *FakeSynopsysV1
	ns   string
}

var opssightsResource = schema.GroupVersionResource{Group: "synopsys", Version: "v1", Resource: "opssights"}

var opssightsKind = schema.GroupVersionKind{Group: "synopsys", Version: "v1", Kind: "OpsSight"}

// Get takes name of the opsSight, and returns the corresponding opsSight object, and an error if there is any.
func (c *FakeOpsSights) Get(name string, options v1.GetOptions) (result *opssight_v1.OpsSight, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(opssightsResource, c.ns, name), &opssight_v1.OpsSight{})

	if obj == nil {
		return nil, err
	}
	return obj.(*opssight_v1.OpsSight), err
}

// List takes label and field selectors, and returns the list of OpsSights that match those selectors.
func (c *FakeOpsSights) List(opts v1.ListOptions) (result *opssight_v1.OpsSightList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(opssightsResource, opssightsKind, c.ns, opts), &opssight_v1.OpsSightList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &opssight_v1.OpsSightList{ListMeta: obj.(*opssight_v1.OpsSightList).ListMeta}
	for _, item := range obj.(*opssight_v1.OpsSightList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested opsSights.
func (c *FakeOpsSights) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(opssightsResource, c.ns, opts))

}

// Create takes the representation of a opsSight and creates it.  Returns the server's representation of the opsSight, and an error, if there is any.
func (c *FakeOpsSights) Create(opsSight *opssight_v1.OpsSight) (result *opssight_v1.OpsSight, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(opssightsResource, c.ns, opsSight), &opssight_v1.OpsSight{})

	if obj == nil {
		return nil, err
	}
	return obj.(*opssight_v1.OpsSight), err
}

// Update takes the representation of a opsSight and updates it. Returns the server's representation of the opsSight, and an error, if there is any.
func (c *FakeOpsSights) Update(opsSight *opssight_v1.OpsSight) (result *opssight_v1.OpsSight, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(opssightsResource, c.ns, opsSight), &opssight_v1.OpsSight{})

	if obj == nil {
		return nil, err
	}
	return obj.(*opssight_v1.OpsSight), err
}

// Delete takes name of the opsSight and deletes it. Returns an error if one occurs.
func (c *FakeOpsSights) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(opssightsResource, c.ns, name), &opssight_v1.OpsSight{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeOpsSights) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(opssightsResource, c.ns, listOptions)

	_, err := c.Fake.Invokes(action, &opssight_v1.OpsSightList{})
	return err
}

// Patch applies the patch and returns the patched opsSight.
func (c *FakeOpsSights) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *opssight_v1.OpsSight, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(opssightsResource, c.ns, name, data, subresources...), &opssight_v1.OpsSight{})

	if obj == nil {
		return nil, err
	}
	return obj.(*opssight_v1.OpsSight), err
}
