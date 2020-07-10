package suitetest

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"time"

	"github.com/cnrancher/octopus-api-server/pkg/apis/octopusapi.cattle.io/v1alpha1"
	"github.com/cnrancher/octopus-api-server/tests/integration/framework"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/apimachinery/pkg/watch"
)

func NewDeviceTemplateTester(k8s *framework.K8sCli) (*DeviceTemplateTester, error) {
	tester := &DeviceTemplateTester{k8s: k8s}
	var err error
	if tester.exampleDeviceTemplateYAML, err = ioutil.ReadFile("../deploy/example_template.yaml"); err != nil {
		return nil, err
	}
	if tester.exampleRevisionYAML, err = ioutil.ReadFile("../deploy/example_template_revision.yaml"); err != nil {
		return nil, err
	}

	var templateJSONData []byte
	if templateJSONData, err = yaml.ToJSON(tester.exampleDeviceTemplateYAML); err != nil {
		return nil, err
	}
	if err = json.Unmarshal(templateJSONData, &tester.deviceTemplate); err != nil {
		return nil, err
	}
	var revisionJSONData []byte
	if revisionJSONData, err = yaml.ToJSON(tester.exampleRevisionYAML); err != nil {
		return nil, err
	}
	if err = json.Unmarshal(revisionJSONData, &tester.revision); err != nil {
		return nil, err
	}

	tester.deviceTemplateGVR = schema.GroupVersionResource{
		Group:    tester.deviceTemplate.GroupVersionKind().Group,
		Version:  tester.deviceTemplate.GroupVersionKind().Version,
		Resource: "devicetemplates",
	}
	tester.revisionGVR = schema.GroupVersionResource{
		Group:    tester.revision.GroupVersionKind().Group,
		Version:  tester.revision.GroupVersionKind().Version,
		Resource: "devicetemplaterevisions",
	}

	return tester, nil
}

type DeviceTemplateTester struct {
	k8s                       *framework.K8sCli
	exampleRevisionYAML       []byte
	exampleDeviceTemplateYAML []byte

	revisionGVR       schema.GroupVersionResource
	deviceTemplateGVR schema.GroupVersionResource

	deviceTemplate *v1alpha1.DeviceTemplate
	revision       *v1alpha1.DeviceTemplateRevision
}

func (t *DeviceTemplateTester) CreateExampleDeviceTemplate() error {
	cli := t.k8s.Dyclient.Resource(t.deviceTemplateGVR)
	uns, err := framework.ConvertToUnstruct(t.deviceTemplate)
	if err != nil {
		return err
	}
	_, err = cli.Namespace("default").Create(context.Background(), uns, metav1.CreateOptions{})
	if err != nil {
		return err
	}
	_, err = cli.Namespace("default").Get(context.Background(), t.deviceTemplate.Name, metav1.GetOptions{})
	if err != nil {
		return err
	}
	return nil
}

func (t *DeviceTemplateTester) CreateExampleRevision() error {
	cli := t.k8s.Dyclient.Resource(t.revisionGVR)
	uns, err := framework.ConvertToUnstruct(t.revision)
	if err != nil {
		return err
	}
	_, err = cli.Namespace("default").Create(context.Background(), uns, metav1.CreateOptions{})
	if err != nil {
		return err
	}
	_, err = cli.Namespace("default").Get(context.Background(), t.revision.Name, metav1.GetOptions{})
	if err != nil {
		return err
	}
	return nil
}

func (t *DeviceTemplateTester) DeleteOwnerReferences() error {
	cli := t.k8s.Dyclient.Resource(t.deviceTemplateGVR)
	err := cli.Namespace(t.deviceTemplate.Namespace).Delete(context.Background(), t.deviceTemplate.Name, metav1.DeleteOptions{})
	if err != nil {
		return err
	}

	rcli := t.k8s.Dyclient.Resource(t.revisionGVR)
	revisionWatch, err := rcli.Namespace(t.deviceTemplate.Namespace).Watch(context.Background(), metav1.ListOptions{})
	if err != nil {
		return err
	}

	for {
		select {
		case result := <-revisionWatch.ResultChan():
			unStruct := result.Object.(*unstructured.Unstructured)
			data, _ := unStruct.MarshalJSON()
			var revision v1alpha1.DeviceTemplateRevision
			if err = json.Unmarshal(data, &revision); err != nil {
				return err
			}
			if result.Type == watch.Deleted && revision.Name == t.revision.Name {
				return nil
			}
		case <-time.After(1 * time.Minute):
			return errors.New("DeviceTemplateTester watch delete revision timeout")
		}
	}
}

func (t *DeviceTemplateTester) Clean() error {
	cli := t.k8s.Dyclient.Resource(t.deviceTemplateGVR)
	if _, err := cli.Namespace(t.deviceTemplate.Namespace).Get(context.Background(), t.deviceTemplate.Name, metav1.GetOptions{}); err == nil {
		if err = cli.Namespace(t.deviceTemplate.Namespace).Delete(context.Background(), t.deviceTemplate.Name, metav1.DeleteOptions{}); err != nil {
			return err
		}
	}
	cli = t.k8s.Dyclient.Resource(t.revisionGVR)
	if _, err := cli.Namespace(t.revision.Namespace).Get(context.Background(), t.revision.Name, metav1.GetOptions{}); err == nil {
		if err = cli.Namespace(t.revision.Namespace).Delete(context.Background(), t.revision.Name, metav1.DeleteOptions{}); err != nil {
			return err
		}
	}
	return nil
}
