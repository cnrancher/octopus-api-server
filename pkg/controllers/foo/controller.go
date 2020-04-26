package foo

import (
	"context"

	foov1 "github.com/cnrancher/edge-api-server/pkg/apis/edgeapi.cattle.io/v1alpha1"
	foocontroller "github.com/cnrancher/edge-api-server/pkg/generated/controllers/edgeapi.cattle.io/v1alpha1"
	"github.com/rancher/wrangler/pkg/apply"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	name = "foo-controller"
)

func Register(ctx context.Context, apply apply.Apply, foos foocontroller.FooController) {
	apply = apply.WithSetID(name).WithCacheTypes(foos)
	controller := &handler{
		foos:     foos,
		fooCache: foos.Cache(),
		apply:    apply,
	}

	foos.OnChange(ctx, name, controller.onChanged)
	foos.OnRemove(ctx, name, controller.onRemove)
}

type handler struct {
	foos     foocontroller.FooClient
	fooCache foocontroller.FooCache
	apply    apply.Apply
}

func (h *handler) onChanged(key string, foo *foov1.Foo) (*foov1.Foo, error) {
	// foo will be nil if key is deleted from cache
	if foo == nil {
		return nil, nil
	}
	fooCopy := foo.DeepCopy()
	return h.foos.Update(fooCopy)
}

func (h *handler) onRemove(key string, foo *foov1.Foo) (*foov1.Foo, error) {
	if key == "" {
		return foo, nil
	}
	return foo, h.foos.Delete(foo.Namespace, key, &metav1.DeleteOptions{})
}
