package devicetemplaterevisionapi

import (
	"context"
	"errors"
	"fmt"

	"github.com/cnrancher/octopus-api-server/pkg/apis/octopusapi.cattle.io/v1alpha1"
	"github.com/cnrancher/octopus-api-server/pkg/auth"
	controller "github.com/cnrancher/octopus-api-server/pkg/generated/controllers/octopusapi.cattle.io/v1alpha1"
	"github.com/cnrancher/octopus-api-server/pkg/util"
	"github.com/rancher/apiserver/pkg/types"
	"github.com/rancher/steve/pkg/accesscontrol"
	"github.com/rancher/wrangler/pkg/data/convert"
	"github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
)

type Store struct {
	types.Store
	asl                      accesscontrol.AccessSetLookup
	client                   dynamic.Interface
	ctx                      context.Context
	auth                     auth.Authenticator
	deviceTemplateController controller.DeviceTemplateController
}

const (
	templateDeviceTypeName    = "octopusapi.cattle.io/template-device-type"
	templateDeviceVersionName = "octopusapi.cattle.io/template-device-version"
	templateOwnerName         = "octopusapi.cattle.io/template-owner"
	templateRevisionReference = "octopusapi.cattle.io/template-revision-reference"
)

func (s *Store) Create(apiOp *types.APIRequest, schema *types.APISchema, data types.APIObject) (types.APIObject, error) {
	var deviceTemplateRevision v1alpha1.DeviceTemplateRevision
	err := convert.ToObj(data.Data(), &deviceTemplateRevision)
	if err != nil {
		logrus.Errorf("failed to convert device template revision data, error: %s", err.Error())
		return data, err
	}

	if err := validateTemplateRequest(&deviceTemplateRevision.Spec); err != nil {
		logrus.Errorf("invalid device template revision request, error: %s", err.Error())
		return data, err
	}

	if deviceTemplateRevision.Name == "" {
		deviceTemplateRevision.Name = fmt.Sprintf("%s-revision-%s", deviceTemplateRevision.Spec.DeviceTemplateName, deviceTemplateRevision.Spec.DisplayName)
	}

	deviceTemplate, err := s.deviceTemplateController.Get(deviceTemplateRevision.Namespace, deviceTemplateRevision.Spec.DeviceTemplateName, metav1.GetOptions{})
	if err != nil {
		logrus.Errorf("device template is not exist, error: %s", err.Error())
		return data, err
	}

	if err := s.validTemplateSpec(&deviceTemplateRevision, deviceTemplate); err != nil {
		logrus.Errorf("valid template spec error: %s", err.Error())
		return data, err
	}

	authed, user, err := s.auth.Authenticate(apiOp.Request)
	if !authed || err != nil {
		logrus.Error("Invalid user error:", err.Error())
		return data, err
	}

	deviceTemplateRevision.Labels = map[string]string{
		templateDeviceTypeName:    deviceTemplate.Spec.DeviceKind,
		templateDeviceVersionName: deviceTemplate.Spec.DeviceVersion,
		templateRevisionReference: deviceTemplateRevision.Spec.DeviceTemplateName,
		templateOwnerName:         user,
	}

	err = convert.ToObj(deviceTemplateRevision, &data.Object)
	if err != nil {
		logrus.Errorf("failed to convert device template revision data, error: %s", err.Error())
		return data, err
	}

	return s.Store.Create(apiOp, schema, data)
}

func (s *Store) Update(apiOp *types.APIRequest, schema *types.APISchema, data types.APIObject, id string) (types.APIObject, error) {
	var deviceTemplateRevision v1alpha1.DeviceTemplateRevision
	err := convert.ToObj(data.Data(), &deviceTemplateRevision)
	if err != nil {
		logrus.Errorf("failed to convert device template revision data, error: %s", err.Error())
		return data, err
	}

	if err := validateTemplateRequest(&deviceTemplateRevision.Spec); err != nil {
		logrus.Errorf("invalid device template revision request, error: %s", err.Error())
		return data, err
	}

	deviceTemplate, err := s.deviceTemplateController.Get(deviceTemplateRevision.Namespace, deviceTemplateRevision.Spec.DeviceTemplateName, metav1.GetOptions{})
	if err != nil {
		logrus.Errorf("device template is not exist, error: %s", err.Error())
		return data, err
	}

	if err := s.validTemplateSpec(&deviceTemplateRevision, deviceTemplate); err != nil {
		logrus.Errorf("valid template spec error: %s", err.Error())
		return data, err
	}

	err = convert.ToObj(deviceTemplateRevision, &data.Object)
	if err != nil {
		logrus.Errorf("failed to convert device template revision data, error: %s", err.Error())
		return data, err
	}

	return s.Store.Update(apiOp, schema, data, id)
}

func validateTemplateRequest(spec *v1alpha1.DeviceTemplateRevisionSpec) error {
	if spec.DisplayName == "" {
		return errors.New("displayName is required of DeviceTemplateRevision")
	}
	if spec.DeviceTemplateName == "" {
		return errors.New("deviceTemplateName is required of DeviceTemplateRevision")
	}
	if spec.DeviceTemplateAPIVersion == "" {
		return errors.New("deviceTemplateAPIVersion is required of DeviceTemplateRevision")
	}
	if spec.TemplateSpec == nil {
		return errors.New("templateSpec is required of DeviceTemplateRevision")
	}
	return nil
}

func (s *Store) validTemplateSpec(revision *v1alpha1.DeviceTemplateRevision, deviceTemplate *v1alpha1.DeviceTemplate) error {
	deviceGroup := deviceTemplate.Spec.DeviceGroup
	deviceVersion := deviceTemplate.Spec.DeviceVersion
	deviceResource := deviceTemplate.Spec.DeviceResource
	deviceKind := deviceTemplate.Spec.DeviceKind

	tempStr := util.GenerateRandomTempKey(7)
	device := unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": fmt.Sprintf("%s/%s", deviceGroup, deviceVersion),
			"kind":       deviceKind,
			"metadata": map[string]interface{}{
				"name":      fmt.Sprintf("devicetemplate-%s", tempStr),
				"namespace": revision.Namespace,
			},
			"spec": revision.Spec.TemplateSpec,
		},
	}

	opt := metav1.CreateOptions{DryRun: []string{metav1.DryRunAll}}

	resource := schema.GroupVersionResource{
		Group:    deviceGroup,
		Version:  deviceVersion,
		Resource: deviceResource,
	}

	crdClient := s.client.Resource(resource)
	if _, err := crdClient.Namespace(revision.Namespace).Create(s.ctx, &device, opt); err != nil {
		return err
	}

	return nil
}
