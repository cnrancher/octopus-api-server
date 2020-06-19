package framework

import (
	"context"

	edgev1Schema "github.com/cnrancher/edge-api-server/pkg/apis/edgeapi.cattle.io/v1alpha1/schema"
	"github.com/cnrancher/edge-api-server/pkg/controllers/devicetemplate"
	"github.com/cnrancher/edge-api-server/pkg/controllers/devicetemplaterevision"
	edgev1 "github.com/cnrancher/edge-api-server/pkg/generated/controllers/edgeapi.cattle.io"
	"github.com/rancher/steve/pkg/debug"
	"github.com/rancher/wrangler/pkg/apply"
	"github.com/rancher/wrangler/pkg/crd"
	"github.com/rancher/wrangler/pkg/start"
	discovery "k8s.io/client-go/discovery"
	"k8s.io/client-go/rest"
)

var (
	debugConfig debug.Config
)

func SetupDeviceTemplateController(ctx context.Context, cli *discovery.DiscoveryClient, restConfig *rest.Config) error {
	factory, err := edgev1.NewFactoryFromConfig(restConfig)
	if err != nil {
		return err
	}

	crdFactory, err := crd.NewFactoryFromClient(restConfig)
	if err != nil {
		return err
	}

	crdFactory.BatchCreateCRDs(ctx, crd.NamespacedTypes(edgev1Schema.SetAndGetCRDName("DeviceTemplate"))...)
	crdFactory.BatchCreateCRDs(ctx, crd.NamespacedTypes(edgev1Schema.SetAndGetCRDName("DeviceTemplateRevision"))...)
	if err = crdFactory.BatchWait(); err != nil {
		return err
	}

	objectSetApply := apply.New(cli, apply.NewClientFactory(restConfig))
	devicetemplate.Register(ctx, objectSetApply, factory.Edgeapi().V1alpha1().DeviceTemplate())
	devicetemplaterevision.Register(ctx, objectSetApply, factory.Edgeapi().V1alpha1().DeviceTemplateRevision(), factory.Edgeapi().V1alpha1().DeviceTemplate())

	if err := start.All(ctx, 1, factory); err != nil {
		return err
	}

	return nil
}
