package devicetemplate

import (
	"context"
	"time"

	"github.com/cnrancher/octopus-api-server/pkg/apis/octopusapi.cattle.io/v1alpha1"
	controllers "github.com/cnrancher/octopus-api-server/pkg/generated/controllers/octopusapi.cattle.io/v1alpha1"
	"github.com/rancher/wrangler/pkg/apply"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	name = "device-template-controller"
)

type Controller struct {
	context            context.Context
	templateController controllers.DeviceTemplateController
	apply              apply.Apply
}

func Register(ctx context.Context, apply apply.Apply, devicetempaltes controllers.DeviceTemplateController) {
	ctrl := &Controller{
		context:            ctx,
		templateController: devicetempaltes,
		apply:              apply,
	}
	devicetempaltes.OnChange(ctx, name, ctrl.OnChanged)
	devicetempaltes.OnRemove(ctx, name, ctrl.OnRemoved)
}

func (c *Controller) OnChanged(key string, obj *v1alpha1.DeviceTemplate) (*v1alpha1.DeviceTemplate, error) {
	if key == "" {
		return nil, nil
	}

	if obj == nil || obj.DeletionTimestamp != nil {
		return nil, nil
	}
	objCopy := obj.DeepCopy()
	objCopy.Status.UpdatedAt = metav1.Time{Time: time.Now()}
	return c.templateController.Update(objCopy)
}

func (c *Controller) OnRemoved(key string, obj *v1alpha1.DeviceTemplate) (*v1alpha1.DeviceTemplate, error) {
	if key == "" {
		return obj, nil
	}
	return obj, c.templateController.Delete(obj.Namespace, obj.Name, &metav1.DeleteOptions{})
}
