package suitetest

import (
	"context"
	"testing"
	"time"

	"github.com/cnrancher/octopus-api-server/tests/integration/cluster"
	"github.com/cnrancher/octopus-api-server/tests/integration/framework"
	"github.com/stretchr/testify/suite"
)

type EdgeTestSuite struct {
	suite.Suite
	k8s    *framework.K8sCli
	cancel context.CancelFunc

	deviceTemplateTester *DeviceTemplateTester
}

func (s *EdgeTestSuite) SetupSuite() {
	s.T().Log("SetupSuite")
	var err error
	err = cluster.StepK3dCluster()
	s.NoError(err, "SetupTest StepK3dCluster error")
	if err != nil {
		cluster.CleanK3dCluster()
		s.FailNow("step k3d cluster error", err)
	}
	s.k8s, err = framework.NewK8sCli()
	s.NoError(err, "SetupTest create k8s client error")
	if err != nil {
		s.FailNow("new k8s client error", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	s.cancel = cancel
	err = framework.SetupCustomAPIController(ctx, s.k8s.Clientset.DiscoveryClient, s.k8s.Cfg)
	s.NoError(err, "SetupCustomAPIController error")

	s.deviceTemplateTester, err = NewDeviceTemplateTester(s.k8s)
	s.NoError(err, "SetupTest NewDeviceTemplateTester error")
}

func (s *EdgeTestSuite) TestDeviceTemplate() {
	err := s.deviceTemplateTester.CreateExampleDeviceTemplate()
	s.NoError(err, "TestDeviceTemplate Test CreateExampleDeviceTemplate error")

	err = s.deviceTemplateTester.CreateExampleRevision()
	s.NoError(err, "TestDeviceTemplate Test CreateExampleRevision error")

	err = s.deviceTemplateTester.DeleteOwnerReferences()
	s.NoError(err, "TestDeviceTemplate Test DeleteOwnerReferences error")
}

func (s *EdgeTestSuite) TearDownSuite() {
	s.T().Log("TearDownSuite")

	err := s.deviceTemplateTester.Clean()
	s.NoError(err, "DeviceTemplateTester Clean error")

	time.Sleep(5 * time.Second)
	err = cluster.CleanK3dCluster()
	s.NoError(err)
}

func (s *EdgeTestSuite) AfterTest(_, name string) {
	s.T().Log("After test ", name)
}

func TestEdgetSuite(t *testing.T) {
	suite.Run(t, new(EdgeTestSuite))
}
