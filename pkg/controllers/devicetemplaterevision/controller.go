package devicetemplaterevision

import (
	"context"
	"fmt"
	"time"

	"github.com/cnrancher/octopus-api-server/pkg/apis/octopusapi.cattle.io/v1alpha1"
	controllers "github.com/cnrancher/octopus-api-server/pkg/generated/controllers/octopusapi.cattle.io/v1alpha1"
	"github.com/rancher/wrangler/pkg/apply"
	apierrs "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
)

const (
	name = "device-template-revision-controller"
)

const (
	templateRevisionReference = "octopusapi.cattle.io/template-revision-reference"
	deviceTemplateKindName    = "DeviceTemplate"
)

const (
	templateDeviceTypeName    = "octopusapi.cattle.io/template-device-type"
	templateDeviceVersionName = "octopusapi.cattle.io/template-device-version"
)

type Controller struct {
	context            context.Context
	templateController controllers.DeviceTemplateController
	templateLister     controllers.DeviceTemplateCache
	revisionController controllers.DeviceTemplateRevisionController
	revisionLister     controllers.DeviceTemplateRevisionCache
	apply              apply.Apply
}

func Register(ctx context.Context, apply apply.Apply, revisionController controllers.DeviceTemplateRevisionController, templateController controllers.DeviceTemplateController) {
	ctrl := &Controller{
		context:            ctx,
		templateController: templateController,
		templateLister:     templateController.Cache(),
		revisionController: revisionController,
		revisionLister:     revisionController.Cache(),
		apply:              apply,
	}
	revisionController.OnChange(ctx, name, ctrl.OnChanged)
	revisionController.OnRemove(ctx, name, ctrl.OnRemoved)
}

func (c *Controller) OnChanged(key string, obj *v1alpha1.DeviceTemplateRevision) (*v1alpha1.DeviceTemplateRevision, error) {
	if key == "" {
		return nil, nil
	}

	if obj == nil || obj.DeletionTimestamp != nil {
		return nil, nil
	}
	deviceTemplate, err := c.templateLister.Get(obj.Namespace, obj.Spec.DeviceTemplateName)
	if err != nil {
		return nil, err
	}

	if len(obj.Labels) <= 0 {
		obj.Labels = map[string]string{
			templateDeviceTypeName:    deviceTemplate.Spec.DeviceKind,
			templateDeviceVersionName: deviceTemplate.Spec.DeviceVersion,
			templateRevisionReference: obj.Spec.DeviceTemplateName,
		}
	}

	objCopy := obj.DeepCopy()
	objCopy.Status.UpdatedAt = metav1.Time{Time: time.Now()}
	objCopy.OwnerReferences = append(objCopy.OwnerReferences[:0], SetRevisionOwner(objCopy, deviceTemplate.UID))

	if err := c.SyncDeviceTemplateDefaultRevision(objCopy, deviceTemplate, false); err != nil {
		return nil, err
	}

	return c.revisionController.Update(objCopy)
}

func (c *Controller) OnRemoved(key string, obj *v1alpha1.DeviceTemplateRevision) (*v1alpha1.DeviceTemplateRevision, error) {
	if key == "" {
		return obj, nil
	}

	deviceTemplate, err := c.templateLister.Get(obj.Namespace, obj.Spec.DeviceTemplateName)
	if err != nil {
		if !apierrs.IsNotFound(err) {
			return nil, err
		}
		return obj, c.revisionController.Delete(obj.Namespace, obj.Name, &metav1.DeleteOptions{})
	}
	if err := c.SyncDeviceTemplateDefaultRevision(obj, deviceTemplate, true); err != nil {
		return nil, err
	}

	return obj, c.revisionController.Delete(obj.Namespace, obj.Name, &metav1.DeleteOptions{})
}

func (c *Controller) SyncDeviceTemplateDefaultRevision(obj *v1alpha1.DeviceTemplateRevision, deviceTemplate *v1alpha1.DeviceTemplate, onRemove bool) error {
	set := labels.Set(map[string]string{templateRevisionReference: obj.Spec.DeviceTemplateName})
	revisions, err := c.revisionLister.List(obj.Namespace, set.AsSelector())
	if err != nil {
		return err
	}
	revisionCount := len(revisions)
	if revisionCount == 1 {
		revisionName := fmt.Sprintf(`%s/%s`, obj.Namespace, revisions[0].Name)
		if onRemove {
			if deviceTemplate.Spec.DefaultRevisionName == revisionName {
				deviceTemplateCopy := deviceTemplate.DeepCopy()
				deviceTemplateCopy.Spec.DefaultRevisionName = ""
				if _, err := c.templateController.Update(deviceTemplateCopy); err != nil {
					return err
				}
			}
		} else {
			if deviceTemplate.Spec.DefaultRevisionName != revisionName {
				deviceTemplateCopy := deviceTemplate.DeepCopy()
				deviceTemplateCopy.Spec.DefaultRevisionName = fmt.Sprintf(`%s/%s`, obj.Namespace, obj.Name)
				if _, err := c.templateController.Update(deviceTemplateCopy); err != nil {
					return err
				}
			}
		}
	}

	if revisionCount == 0 && deviceTemplate.Spec.DefaultRevisionName != "" {
		deviceTemplateCopy := deviceTemplate.DeepCopy()
		deviceTemplateCopy.Spec.DefaultRevisionName = ""
		if _, err := c.templateController.Update(deviceTemplateCopy); err != nil {
			return err
		}
	}

	return nil
}

func SetRevisionOwner(obj *v1alpha1.DeviceTemplateRevision, uid types.UID) metav1.OwnerReference {
	isController := true
	return metav1.OwnerReference{
		APIVersion: obj.Spec.DeviceTemplateAPIVersion,
		Controller: &isController,
		Kind:       deviceTemplateKindName,
		Name:       obj.Spec.DeviceTemplateName,
		UID:        uid,
	}
}
