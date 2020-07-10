package controllers

import (
	"context"

	"github.com/cnrancher/octopus-api-server/pkg/controllers/devicetemplaterevision"
	"github.com/cnrancher/octopus-api-server/pkg/controllers/settings"

	"github.com/cnrancher/octopus-api-server/pkg/apis/octopusapi.cattle.io/v1alpha1/schema"
	"github.com/cnrancher/octopus-api-server/pkg/controllers/authtoken"
	"github.com/cnrancher/octopus-api-server/pkg/controllers/catalog"
	"github.com/cnrancher/octopus-api-server/pkg/controllers/devicetemplate"
	octopusv1 "github.com/cnrancher/octopus-api-server/pkg/generated/controllers/octopusapi.cattle.io"
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
	factory, err := octopusv1.NewFactoryFromConfig(restConfig)
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

	catalog.Register(ctx, objectSetApply, factory.Octopusapi().V1alpha1().Catalog())
	devicetemplate.Register(ctx, objectSetApply, factory.Octopusapi().V1alpha1().DeviceTemplate())
	devicetemplaterevision.Register(ctx, objectSetApply, factory.Octopusapi().V1alpha1().DeviceTemplateRevision(), factory.Octopusapi().V1alpha1().DeviceTemplate())
	settings.Register(ctx, factory.Octopusapi().V1alpha1().Setting())
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
		schema.SetAndGetCRDName("Catalog"),
		schema.SetAndGetCRDName("DeviceTemplate"),
		schema.SetAndGetCRDName("DeviceTemplateRevision"))...)

	factory.BatchCreateCRDs(ctx, crd.NonNamespacedTypes(
		schema.SetAndGetCRDName("Setting"))...)

	return factory.BatchWait()
}
