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

package v1alpha1

import (
	"context"
	"time"

	v1alpha1 "github.com/cnrancher/edge-api-server/pkg/apis/edgeapi.cattle.io/v1alpha1"
	clientset "github.com/cnrancher/edge-api-server/pkg/generated/clientset/versioned/typed/edgeapi.cattle.io/v1alpha1"
	informers "github.com/cnrancher/edge-api-server/pkg/generated/informers/externalversions/edgeapi.cattle.io/v1alpha1"
	listers "github.com/cnrancher/edge-api-server/pkg/generated/listers/edgeapi.cattle.io/v1alpha1"
	"github.com/rancher/wrangler/pkg/generic"
	"k8s.io/apimachinery/pkg/api/equality"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"
)

type SettingHandler func(string, *v1alpha1.Setting) (*v1alpha1.Setting, error)

type SettingController interface {
	generic.ControllerMeta
	SettingClient

	OnChange(ctx context.Context, name string, sync SettingHandler)
	OnRemove(ctx context.Context, name string, sync SettingHandler)
	Enqueue(namespace, name string)
	EnqueueAfter(namespace, name string, duration time.Duration)

	Cache() SettingCache
}

type SettingClient interface {
	Create(*v1alpha1.Setting) (*v1alpha1.Setting, error)
	Update(*v1alpha1.Setting) (*v1alpha1.Setting, error)

	Delete(namespace, name string, options *metav1.DeleteOptions) error
	Get(namespace, name string, options metav1.GetOptions) (*v1alpha1.Setting, error)
	List(namespace string, opts metav1.ListOptions) (*v1alpha1.SettingList, error)
	Watch(namespace string, opts metav1.ListOptions) (watch.Interface, error)
	Patch(namespace, name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.Setting, err error)
}

type SettingCache interface {
	Get(namespace, name string) (*v1alpha1.Setting, error)
	List(namespace string, selector labels.Selector) ([]*v1alpha1.Setting, error)

	AddIndexer(indexName string, indexer SettingIndexer)
	GetByIndex(indexName, key string) ([]*v1alpha1.Setting, error)
}

type SettingIndexer func(obj *v1alpha1.Setting) ([]string, error)

type settingController struct {
	controllerManager *generic.ControllerManager
	clientGetter      clientset.SettingsGetter
	informer          informers.SettingInformer
	gvk               schema.GroupVersionKind
}

func NewSettingController(gvk schema.GroupVersionKind, controllerManager *generic.ControllerManager, clientGetter clientset.SettingsGetter, informer informers.SettingInformer) SettingController {
	return &settingController{
		controllerManager: controllerManager,
		clientGetter:      clientGetter,
		informer:          informer,
		gvk:               gvk,
	}
}

func FromSettingHandlerToHandler(sync SettingHandler) generic.Handler {
	return func(key string, obj runtime.Object) (ret runtime.Object, err error) {
		var v *v1alpha1.Setting
		if obj == nil {
			v, err = sync(key, nil)
		} else {
			v, err = sync(key, obj.(*v1alpha1.Setting))
		}
		if v == nil {
			return nil, err
		}
		return v, err
	}
}

func (c *settingController) Updater() generic.Updater {
	return func(obj runtime.Object) (runtime.Object, error) {
		newObj, err := c.Update(obj.(*v1alpha1.Setting))
		if newObj == nil {
			return nil, err
		}
		return newObj, err
	}
}

func UpdateSettingDeepCopyOnChange(client SettingClient, obj *v1alpha1.Setting, handler func(obj *v1alpha1.Setting) (*v1alpha1.Setting, error)) (*v1alpha1.Setting, error) {
	if obj == nil {
		return obj, nil
	}

	copyObj := obj.DeepCopy()
	newObj, err := handler(copyObj)
	if newObj != nil {
		copyObj = newObj
	}
	if obj.ResourceVersion == copyObj.ResourceVersion && !equality.Semantic.DeepEqual(obj, copyObj) {
		return client.Update(copyObj)
	}

	return copyObj, err
}

func (c *settingController) AddGenericHandler(ctx context.Context, name string, handler generic.Handler) {
	c.controllerManager.AddHandler(ctx, c.gvk, c.informer.Informer(), name, handler)
}

func (c *settingController) AddGenericRemoveHandler(ctx context.Context, name string, handler generic.Handler) {
	removeHandler := generic.NewRemoveHandler(name, c.Updater(), handler)
	c.controllerManager.AddHandler(ctx, c.gvk, c.informer.Informer(), name, removeHandler)
}

func (c *settingController) OnChange(ctx context.Context, name string, sync SettingHandler) {
	c.AddGenericHandler(ctx, name, FromSettingHandlerToHandler(sync))
}

func (c *settingController) OnRemove(ctx context.Context, name string, sync SettingHandler) {
	removeHandler := generic.NewRemoveHandler(name, c.Updater(), FromSettingHandlerToHandler(sync))
	c.AddGenericHandler(ctx, name, removeHandler)
}

func (c *settingController) Enqueue(namespace, name string) {
	c.controllerManager.Enqueue(c.gvk, c.informer.Informer(), namespace, name)
}

func (c *settingController) EnqueueAfter(namespace, name string, duration time.Duration) {
	c.controllerManager.EnqueueAfter(c.gvk, c.informer.Informer(), namespace, name, duration)
}

func (c *settingController) Informer() cache.SharedIndexInformer {
	return c.informer.Informer()
}

func (c *settingController) GroupVersionKind() schema.GroupVersionKind {
	return c.gvk
}

func (c *settingController) Cache() SettingCache {
	return &settingCache{
		lister:  c.informer.Lister(),
		indexer: c.informer.Informer().GetIndexer(),
	}
}

func (c *settingController) Create(obj *v1alpha1.Setting) (*v1alpha1.Setting, error) {
	return c.clientGetter.Settings(obj.Namespace).Create(context.TODO(), obj, metav1.CreateOptions{})
}

func (c *settingController) Update(obj *v1alpha1.Setting) (*v1alpha1.Setting, error) {
	return c.clientGetter.Settings(obj.Namespace).Update(context.TODO(), obj, metav1.UpdateOptions{})
}

func (c *settingController) Delete(namespace, name string, options *metav1.DeleteOptions) error {
	if options == nil {
		options = &metav1.DeleteOptions{}
	}
	return c.clientGetter.Settings(namespace).Delete(context.TODO(), name, *options)
}

func (c *settingController) Get(namespace, name string, options metav1.GetOptions) (*v1alpha1.Setting, error) {
	return c.clientGetter.Settings(namespace).Get(context.TODO(), name, options)
}

func (c *settingController) List(namespace string, opts metav1.ListOptions) (*v1alpha1.SettingList, error) {
	return c.clientGetter.Settings(namespace).List(context.TODO(), opts)
}

func (c *settingController) Watch(namespace string, opts metav1.ListOptions) (watch.Interface, error) {
	return c.clientGetter.Settings(namespace).Watch(context.TODO(), opts)
}

func (c *settingController) Patch(namespace, name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.Setting, err error) {
	return c.clientGetter.Settings(namespace).Patch(context.TODO(), name, pt, data, metav1.PatchOptions{}, subresources...)
}

type settingCache struct {
	lister  listers.SettingLister
	indexer cache.Indexer
}

func (c *settingCache) Get(namespace, name string) (*v1alpha1.Setting, error) {
	return c.lister.Settings(namespace).Get(name)
}

func (c *settingCache) List(namespace string, selector labels.Selector) ([]*v1alpha1.Setting, error) {
	return c.lister.Settings(namespace).List(selector)
}

func (c *settingCache) AddIndexer(indexName string, indexer SettingIndexer) {
	utilruntime.Must(c.indexer.AddIndexers(map[string]cache.IndexFunc{
		indexName: func(obj interface{}) (strings []string, e error) {
			return indexer(obj.(*v1alpha1.Setting))
		},
	}))
}

func (c *settingCache) GetByIndex(indexName, key string) (result []*v1alpha1.Setting, err error) {
	objs, err := c.indexer.ByIndex(indexName, key)
	if err != nil {
		return nil, err
	}
	result = make([]*v1alpha1.Setting, 0, len(objs))
	for _, obj := range objs {
		result = append(result, obj.(*v1alpha1.Setting))
	}
	return result, nil
}
