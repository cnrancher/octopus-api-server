package framework

import (
	"context"

	"github.com/cnrancher/octopus-api-server/pkg/apis/octopusapi.cattle.io/v1alpha1/schema"
	"github.com/cnrancher/octopus-api-server/pkg/controllers/devicetemplate"
	"github.com/cnrancher/octopus-api-server/pkg/controllers/devicetemplaterevision"
	octopusv1 "github.com/cnrancher/octopus-api-server/pkg/generated/controllers/octopusapi.cattle.io"
	"github.com/rancher/steve/pkg/debug"
	"github.com/rancher/wrangler/pkg/apply"
	"github.com/rancher/wrangler/pkg/crd"
	"github.com/rancher/wrangler/pkg/start"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/rest"
)

var (
	debugConfig debug.Config
)

func SetupCustomAPIController(ctx context.Context, cli *discovery.DiscoveryClient, restConfig *rest.Config) error {
	factory, err := octopusv1.NewFactoryFromConfig(restConfig)
	if err != nil {
		return err
	}

	crdFactory, err := crd.NewFactoryFromClient(restConfig)
	if err != nil {
		return err
	}

	crdFactory.BatchCreateCRDs(ctx, crd.NamespacedTypes(schema.SetAndGetCRDName("DeviceTemplate"))...)
	crdFactory.BatchCreateCRDs(ctx, crd.NamespacedTypes(schema.SetAndGetCRDName("DeviceTemplateRevision"))...)
	if err = crdFactory.BatchWait(); err != nil {
		return err
	}

	objectSetApply := apply.New(cli, apply.NewClientFactory(restConfig))
	devicetemplate.Register(ctx, objectSetApply, factory.Octopusapi().V1alpha1().DeviceTemplate())
	devicetemplaterevision.Register(ctx, objectSetApply, factory.Octopusapi().V1alpha1().DeviceTemplateRevision(), factory.Octopusapi().V1alpha1().DeviceTemplate())

	if err := start.All(ctx, 1, factory); err != nil {
		return err
	}

	return nil
}
