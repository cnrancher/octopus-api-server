package devicetemplate

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"testing"
	"time"

	"github.com/cnrancher/edge-api-server/pkg/apis/edgeapi.cattle.io/v1alpha1"
	"github.com/cnrancher/edge-api-server/tests/integration/framework"
	"github.com/stretchr/testify/suite"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/yaml"
)

type DeviceTemplateTestSuite struct {
	suite.Suite
	k8s *framework.K8sCli

	exampleRevisionYAML       []byte
	exampleDeviceTemplateYAML []byte

	revisionGVR       schema.GroupVersionResource
	deviceTemplateGVR schema.GroupVersionResource

	deviceTemplate *v1alpha1.DeviceTemplate
	revision       *v1alpha1.DeviceTemplateRevision
	cancel         context.CancelFunc
}

func (s *DeviceTemplateTestSuite) SetupSuite() {
	s.T().Log("SetupSuite")
	var err error
	s.k8s, err = framework.NewK8sCli()
	s.NoError(err, "SetupTest create k8s client error")

	ctx, cancel := context.WithCancel(context.Background())
	s.cancel = cancel

	err = framework.SetupDeviceTemplateController(ctx, s.k8s.Clientset.DiscoveryClient, s.k8s.Cfg)
	s.NoError(err, "SetupDeviceTemplateController error")

	templateYAMLPath := filepath.Join("..", "..", "..", "pkg", "controllers", "devicetemplate", "deploy", "example_template.yaml")
	s.exampleDeviceTemplateYAML, err = ioutil.ReadFile(templateYAMLPath)
	s.NoError(err, "SetupTest read examply device template yaml file error")

	revisionYAMLPath := filepath.Join("..", "..", "..", "pkg", "controllers", "devicetemplaterevision", "deploy", "example_template_revision.yaml")
	s.exampleRevisionYAML, err = ioutil.ReadFile(revisionYAMLPath)
	s.NoError(err, "SetupTest read examply revision yaml file error")

	templateJSONData, err := yaml.ToJSON(s.exampleDeviceTemplateYAML)
	s.NoError(err, "SetupTest example devcei template yaml to json error")
	err = json.Unmarshal(templateJSONData, &s.deviceTemplate)
	s.NoError(err, "SetupTest Unmarshal device template error")

	revisionJSONData, err := yaml.ToJSON(s.exampleRevisionYAML)
	s.NoError(err, "SetupTest example revision yaml to json error")
	err = json.Unmarshal(revisionJSONData, &s.revision)
	s.NoError(err, "SetupTest Unmarshal revision error")

	s.deviceTemplateGVR = schema.GroupVersionResource{
		Group:    s.deviceTemplate.GroupVersionKind().Group,
		Version:  s.deviceTemplate.GroupVersionKind().Version,
		Resource: "devicetemplates",
	}
	s.revisionGVR = schema.GroupVersionResource{
		Group:    s.revision.GroupVersionKind().Group,
		Version:  s.revision.GroupVersionKind().Version,
		Resource: "devicetemplaterevisions",
	}

	cli := s.k8s.Dyclient.Resource(s.deviceTemplateGVR)
	if _, err = cli.Namespace(s.deviceTemplate.Namespace).Get(context.Background(), s.deviceTemplate.Name, metav1.GetOptions{}); err == nil {
		err = cli.Namespace(s.deviceTemplate.Namespace).Delete(context.Background(), s.deviceTemplate.Name, metav1.DeleteOptions{})
		s.NoError(err)
	}

	cli = s.k8s.Dyclient.Resource(s.revisionGVR)
	if _, err = cli.Namespace(s.revision.Namespace).Get(context.Background(), s.revision.Name, metav1.GetOptions{}); err == nil {
		err = cli.Namespace(s.revision.Namespace).Delete(context.Background(), s.revision.Name, metav1.DeleteOptions{})
		s.NoError(err)
	}
}

func (s *DeviceTemplateTestSuite) TestCreateExampleDeviceTemplate() {
	cli := s.k8s.Dyclient.Resource(s.deviceTemplateGVR)
	uns, err := framework.ConvertToUnstruct(s.deviceTemplate)
	s.NoError(err, "conver to unstruct error")
	_, err = cli.Namespace("default").Create(context.Background(), uns, metav1.CreateOptions{})
	s.NoError(err, "dynamic create device template error")
}

func (s *DeviceTemplateTestSuite) TestCreateExampleRevision() {
	cli := s.k8s.Dyclient.Resource(s.revisionGVR)
	uns, err := framework.ConvertToUnstruct(s.revision)
	s.NoError(err, "conver to unstruct error")
	_, err = cli.Namespace("default").Create(context.Background(), uns, metav1.CreateOptions{})
	s.NoError(err, "dynamic create revision error")
}

// func (s *DeviceTemplateTestSuite) TestDeleteOwnerReferences() {
// 	cli := s.k8s.Dyclient.Resource(s.deviceTemplateGVR)
// 	err := cli.Namespace(s.deviceTemplate.Namespace).Delete(context.Background(), s.deviceTemplate.Name, metav1.DeleteOptions{})
// 	s.NoError(err, "delete device template error")
// 	time.Sleep(3 * time.Second)
// 	cli = s.k8s.Dyclient.Resource(s.revisionGVR)
// 	_, err = cli.Namespace(s.revision.Namespace).Get(context.Background(), s.revision.Name, metav1.GetOptions{})
// 	s.EqualError(err, `devicetemplaterevisions.edgeapi.cattle.io "my-template-revision" not found`)
// }

func (s *DeviceTemplateTestSuite) AfterTest(suiteName, testName string) {
	s.T().Log(suiteName, " ", testName, " ", "after")
}

func (s *DeviceTemplateTestSuite) TearDownTest() {

}

func (s *DeviceTemplateTestSuite) TearDownSuite() {
	s.T().Log("TearDownSuite")
	cli := s.k8s.Dyclient.Resource(s.deviceTemplateGVR)
	cli.Namespace(s.deviceTemplate.Namespace).Delete(context.Background(), s.deviceTemplate.Name, metav1.DeleteOptions{})
	cli = s.k8s.Dyclient.Resource(s.revisionGVR)
	cli.Namespace(s.revision.Namespace).Delete(context.Background(), s.revision.Name, metav1.DeleteOptions{})

	time.Sleep(5 * time.Second)
}

func TestDeviceTemplateSuite(t *testing.T) {
	suite.Run(t, new(DeviceTemplateTestSuite))
}
