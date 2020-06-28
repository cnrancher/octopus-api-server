package controllers

import (
	"context"

	"github.com/cnrancher/edge-api-server/pkg/controllers/devicetemplaterevision"
	"github.com/cnrancher/edge-api-server/pkg/controllers/settings"

	edgev1Schema "github.com/cnrancher/edge-api-server/pkg/apis/edgeapi.cattle.io/v1alpha1/schema"
	authtoken "github.com/cnrancher/edge-api-server/pkg/controllers/authtoken"
	"github.com/cnrancher/edge-api-server/pkg/controllers/catalog"
	"github.com/cnrancher/edge-api-server/pkg/controllers/devicetemplate"
	edgev1 "github.com/cnrancher/edge-api-server/pkg/generated/controllers/edgeapi.cattle.io"
	corev1 "github.com/rancher/wrangler-api/pkg/generated/controllers/core"
	"github.com/rancher/wrangler/pkg/apply"
	"github.com/rancher/wrangler/pkg/crd"
	"github.com/rancher/wrangler/pkg/start"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/klog"
)

func Setup(ctx context.Context, restConfig *rest.Config, clientSet *kubernetes.Clientset,
	threadiness int) error {
	factory, err := edgev1.NewFactoryFromConfig(restConfig)
	if err != nil {
		klog.Fatalf("Error building sample controllers: %s", err.Error())
	}

	if err = crds(ctx, restConfig); err != nil {
		klog.Fatalf("Error apply CRDs: %s", err.Error())
	}

	cores, err := corev1.NewFactoryFromConfig(restConfig)
	if err != nil {
		klog.Fatalf("Error building kube-system core controllers: %s", err.Error())
	}

	objectSetApply := apply.New(clientSet.DiscoveryClient, apply.NewClientFactory(restConfig))

	catalog.Register(ctx, objectSetApply, factory.Edgeapi().V1alpha1().Catalog())
	devicetemplate.Register(ctx, objectSetApply, factory.Edgeapi().V1alpha1().DeviceTemplate())
	devicetemplaterevision.Register(ctx, objectSetApply, factory.Edgeapi().V1alpha1().DeviceTemplateRevision(), factory.Edgeapi().V1alpha1().DeviceTemplate())
	settings.Register(ctx, factory.Edgeapi().V1alpha1().Setting())
	authtoken.Register(ctx, objectSetApply, cores.Core().V1().Secret())

	if err := start.All(ctx, threadiness, factory, cores); err != nil {
		klog.Fatalf("Error starting: %s", err.Error())
	}
	return nil
}

func crds(ctx context.Context, config *rest.Config) error {
	factory, err := crd.NewFactoryFromClient(config)
	if err != nil {
		return err
	}

	factory.BatchCreateCRDs(ctx, crd.NamespacedTypes(
		edgev1Schema.SetAndGetCRDName("Catalog"),
		edgev1Schema.SetAndGetCRDName("DeviceTemplate"),
		edgev1Schema.SetAndGetCRDName("DeviceTemplateRevision"))...)

	factory.BatchCreateCRDs(ctx, crd.NonNamespacedTypes(
		edgev1Schema.SetAndGetCRDName("Setting"))...)

	return factory.BatchWait()
}
