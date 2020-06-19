/*
Copyright 2020 Rancher Labs, Inc.

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

// Code generated by main. DO NOT EDIT.

package fake

import (
	"context"

	v1alpha1 "github.com/cnrancher/edge-api-server/pkg/apis/edgeapi.cattle.io/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeSettings implements SettingInterface
type FakeSettings struct {
	Fake *FakeEdgeapiV1alpha1
	ns   string
}

var settingsResource = schema.GroupVersionResource{Group: "edgeapi.cattle.io", Version: "v1alpha1", Resource: "settings"}

var settingsKind = schema.GroupVersionKind{Group: "edgeapi.cattle.io", Version: "v1alpha1", Kind: "Setting"}

// Get takes name of the setting, and returns the corresponding setting object, and an error if there is any.
func (c *FakeSettings) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.Setting, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(settingsResource, c.ns, name), &v1alpha1.Setting{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Setting), err
}

// List takes label and field selectors, and returns the list of Settings that match those selectors.
func (c *FakeSettings) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.SettingList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(settingsResource, settingsKind, c.ns, opts), &v1alpha1.SettingList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.SettingList{ListMeta: obj.(*v1alpha1.SettingList).ListMeta}
	for _, item := range obj.(*v1alpha1.SettingList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested settings.
func (c *FakeSettings) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(settingsResource, c.ns, opts))

}

// Create takes the representation of a setting and creates it.  Returns the server's representation of the setting, and an error, if there is any.
func (c *FakeSettings) Create(ctx context.Context, setting *v1alpha1.Setting, opts v1.CreateOptions) (result *v1alpha1.Setting, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(settingsResource, c.ns, setting), &v1alpha1.Setting{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Setting), err
}

// Update takes the representation of a setting and updates it. Returns the server's representation of the setting, and an error, if there is any.
func (c *FakeSettings) Update(ctx context.Context, setting *v1alpha1.Setting, opts v1.UpdateOptions) (result *v1alpha1.Setting, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(settingsResource, c.ns, setting), &v1alpha1.Setting{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Setting), err
}

// Delete takes name of the setting and deletes it. Returns an error if one occurs.
func (c *FakeSettings) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(settingsResource, c.ns, name), &v1alpha1.Setting{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeSettings) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(settingsResource, c.ns, listOpts)

	_, err := c.Fake.Invokes(action, &v1alpha1.SettingList{})
	return err
}

// Patch applies the patch and returns the patched setting.
func (c *FakeSettings) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.Setting, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(settingsResource, c.ns, name, pt, data, subresources...), &v1alpha1.Setting{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Setting), err
}
