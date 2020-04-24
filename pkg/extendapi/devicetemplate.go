package extendapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/cnrancher/edge-api-server/pkg/auth"
	"github.com/cnrancher/edge-api-server/pkg/util"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"

	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const (
	templateNameField       = "edgeapi.cattle.io/device-template-name"
	templateDeviceTypeField = "edgeapi.cattle.io/device-template-type"
	templateCreaterField    = "edgeapi.cattle.io/device-template-creator"
)

type DeviceTemplateRequest struct {
	DeviceType   string                `json:"deviceType"`
	TemplateName string                `json:"templateName"`
	TemplateSpec *runtime.RawExtension `json:"spec"`
	Namespace    string                `json:"namespace"`
	Version      string                `json:"version"`
	DeviceGroup  string                `json:"deviceGroup"`
	ResourceName string                `json:"resourceName"`
}

type DeviceTemplateHandler struct {
	clientset *kubernetes.Clientset
	dyclient  dynamic.Interface
}

func NewDeviceTemplateHandler(client *kubernetes.Clientset, dyclient dynamic.Interface) *DeviceTemplateHandler {
	return &DeviceTemplateHandler{
		clientset: client,
		dyclient:  dyclient,
	}
}

func (h *DeviceTemplateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logrus.Error("DeviceTemplateHandler request body error:", err.Error())
		h.WriteResponse(w, err.Error())
		return
	}

	user, err := h.GetAuthName(r)
	if err != nil {
		logrus.Error("DeviceTemplateHandler request GetAuthName error:", err.Error())
		h.WriteResponse(w, err.Error())
		return
	}

	var request DeviceTemplateRequest

	if err := json.Unmarshal(body, &request); err != nil {
		logrus.Error("DeviceTemplateHandler request body validate json error:", err.Error())
		h.WriteResponse(w, err.Error())
		return
	}

	if err := h.validateRequest(request); err != nil {
		logrus.Error("DeviceTemplateHandler request body validate error:", err.Error())
		h.WriteResponse(w, err.Error())
		return
	}

	if err := h.validTemplate(r.Context(), request); err != nil {
		logrus.Error("DeviceTemplateHandler validTemplate error:", err.Error())
		h.WriteResponse(w, err.Error())
		return
	}

	if err := h.createSecret(r.Context(), request, user); err != nil {
		logrus.Error("DeviceTemplateHandler createSecret error:", err.Error())
		h.WriteResponse(w, err.Error())
		return
	}

	h.WriteResponse(w, "success")
}

func (h *DeviceTemplateHandler) validateRequest(req DeviceTemplateRequest) error {

	if req.DeviceType == "" {
		return errors.New("request body have no field deviceType")
	}

	if req.TemplateName == "" {
		return errors.New("request body have no field templateName")
	}

	if req.Namespace == "" {
		return errors.New("request body have no field namespace")
	}

	if req.Version == "" {
		return errors.New("request body have no field version")
	}

	if req.DeviceGroup == "" {
		return errors.New("request body have no field deviceGroup")
	}

	if req.ResourceName == "" {
		return errors.New("request body have no field resourceName")
	}

	if req.TemplateSpec == nil {
		return errors.New("request body have no field spec")
	}

	return nil
}

func (h *DeviceTemplateHandler) GetAuthName(r *http.Request) (string, error) {
	tokenAuthValue := auth.GetTokenAuthFromRequest(r)
	if tokenAuthValue == "" {
		return "", errors.New("must authenticate")
	}

	tokenName, tokenKey := auth.SplitTokenParts(tokenAuthValue)

	if tokenName == "" || tokenKey == "" {
		return "", errors.New("must authenticate")
	}

	return tokenName, nil
}

func (h *DeviceTemplateHandler) WriteResponse(w http.ResponseWriter, msg string) {
	io.WriteString(w, fmt.Sprintf(`{"msg":%s}`, msg))
}

func (h *DeviceTemplateHandler) createSecret(ctx context.Context, req DeviceTemplateRequest, user string) error {

	js, err := req.TemplateSpec.MarshalJSON()
	if err != nil {
		return err
	}

	stringData := map[string]string{"template": string(js)}

	tempStr := util.GenerateRandomTempKey(7)

	name := fmt.Sprintf("%s-template-%s", strings.ToLower(req.DeviceType), tempStr)

	secret := &apiv1.Secret{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Secret",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Namespace:   req.Namespace,
			Annotations: map[string]string{templateCreaterField: user, templateDeviceTypeField: req.DeviceType, templateNameField: req.TemplateName},
		},
		StringData: stringData,
	}

	if _, err := h.clientset.CoreV1().Secrets(req.Namespace).Create(ctx, secret, metav1.CreateOptions{}); err != nil {
		return err
	}

	return nil
}

func (h *DeviceTemplateHandler) validTemplate(ctx context.Context, req DeviceTemplateRequest) error {

	tempStr := util.GenerateRandomTempKey(7)

	name := fmt.Sprintf("testdevicelink-%s", tempStr)

	obj := unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": fmt.Sprintf("%s/%s", req.DeviceGroup, req.Version),
			"kind":       req.DeviceType,
			"metadata": map[string]interface{}{
				"name":      name,
				"namespace": req.Namespace,
			},
			"spec": req.TemplateSpec,
		},
	}

	opt := metav1.CreateOptions{DryRun: []string{metav1.DryRunAll}}

	resource := schema.GroupVersionResource{
		Group:    req.DeviceGroup,
		Version:  req.Version,
		Resource: req.ResourceName,
	}

	crdClient := h.dyclient.Resource(resource)
	_, err := crdClient.Namespace(req.Namespace).Create(ctx, &obj, opt)

	if err != nil {
		return err
	}

	return nil
}
