package devicetemplateapi

import (
	"context"
	"errors"
	"fmt"

	"github.com/cnrancher/edge-api-server/pkg/apis/edgeapi.cattle.io/v1alpha1"
	"github.com/cnrancher/edge-api-server/pkg/auth"
	"github.com/cnrancher/edge-api-server/pkg/util"
	"github.com/rancher/steve/pkg/accesscontrol"
	"github.com/rancher/steve/pkg/schemaserver/types"
	"github.com/rancher/wrangler/pkg/data/convert"
	"github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
)

type Store struct {
	types.Store
	asl    accesscontrol.AccessSetLookup
	client dynamic.Interface
	ctx    context.Context
	auth   auth.Authenticator
}

const (
	templateDeviceTypeName = "edgeapi.cattle.io/device-template-type"
	templateOwnerName      = "edgeapi.cattle.io/device-template-owner"
)

func (s *Store) Create(apiOp *types.APIRequest, schema *types.APISchema, data types.APIObject) (types.APIObject, error) {
	var deviceTemplate v1alpha1.DeviceTemplate
	err := convert.ToObj(data.Data(), &deviceTemplate)
	if err != nil {
		logrus.Errorf("failed to convert device template data, error: %s", err.Error())
		return data, err
	}

	if err := ValidateTemplateRequest(deviceTemplate.Spec); err != nil {
		logrus.Errorf("invalid device template request, error: %s", err.Error())
		return data, err
	}

	if err := ValidTemplateSpec(s.ctx, &deviceTemplate, s.client); err != nil {
		return data, err
	}

	authed, user, err := s.auth.Authenticate(apiOp.Request)
	if !authed || err != nil {
		logrus.Error("Invalid user error:", err.Error())
		return data, err
	}

	deviceTemplate.Annotations = map[string]string{
		templateDeviceTypeName: deviceTemplate.Spec.DeviceKind,
		templateOwnerName:      user,
	}
	err = convert.ToObj(deviceTemplate, &data.Object)
	if err != nil {
		logrus.Errorf("failed to convert device template data, error: %s", err.Error())
		return data, err
	}
	return s.Store.Create(apiOp, schema, data)
}

func (s *Store) Update(apiOp *types.APIRequest, schema *types.APISchema, data types.APIObject, id string) (types.APIObject, error) {
	var deviceTemplate v1alpha1.DeviceTemplate
	err := convert.ToObj(data.Data(), &deviceTemplate)
	if err != nil {
		logrus.Errorf("failed to parse device template data, error: %s", err.Error())
		return data, err
	}

	if err := ValidateTemplateRequest(deviceTemplate.Spec); err != nil {
		logrus.Errorf("invalid device template request, error: %s", err.Error())
		return data, err
	}

	if err := ValidTemplateSpec(s.ctx, &deviceTemplate, s.client); err != nil {
		return data, err
	}

	return s.Store.Update(apiOp, schema, data, id)
}

func ValidateTemplateRequest(spec v1alpha1.DeviceTemplateSpec) error {
	if spec.DeviceKind == "" {
		return errors.New("deviceKind is required of DeviceTemplate")
	}
	if spec.DeviceVersion == "" {
		return errors.New("deviceVersion is required of DeviceTemplate")
	}
	if spec.DeviceGroup == "" {
		return errors.New("deviceGroup is required of DeviceTemplate")
	}
	if spec.DeviceResource == "" {
		return errors.New("deviceResource is required of DeviceTemplate")
	}
	if spec.TemplateSpec == nil {
		return errors.New("templateSpec is required of DeviceTemplate")
	}
	return nil
}

func ValidTemplateSpec(ctx context.Context, obj *v1alpha1.DeviceTemplate, client dynamic.Interface) error {

	tempStr := util.GenerateRandomTempKey(7)
	device := unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": fmt.Sprintf("%s/%s", obj.Spec.DeviceGroup, obj.Spec.DeviceVersion),
			"kind":       obj.Spec.DeviceKind,
			"metadata": map[string]interface{}{
				"name":      fmt.Sprintf("devicetemplate-%s", tempStr),
				"namespace": obj.Namespace,
			},
			"spec": obj.Spec.TemplateSpec,
		},
	}

	opt := metav1.CreateOptions{DryRun: []string{metav1.DryRunAll}}

	resource := schema.GroupVersionResource{
		Group:    obj.Spec.DeviceGroup,
		Version:  obj.Spec.DeviceVersion,
		Resource: obj.Spec.DeviceResource,
	}

	crdClient := client.Resource(resource)
	if _, err := crdClient.Namespace(obj.Namespace).Create(ctx, &device, opt); err != nil {
		return err
	}

	return nil
}
