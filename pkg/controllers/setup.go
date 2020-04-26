package controllers

import (
	"context"

	edgev1Schema "github.com/cnrancher/edge-api-server/pkg/apis/edgeapi.cattle.io/v1alpha1/schema"
	foocontroller "github.com/cnrancher/edge-api-server/pkg/controllers/foo"
	foov1 "github.com/cnrancher/edge-api-server/pkg/generated/controllers/edgeapi.cattle.io"
	"github.com/rancher/wrangler/pkg/apply"
	"github.com/rancher/wrangler/pkg/crd"
	"github.com/rancher/wrangler/pkg/start"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/klog"
)

func Setup(ctx context.Context, restConfig *rest.Config, clientSet *kubernetes.Clientset, threadiness int) error {

	foos, err := foov1.NewFactoryFromConfig(restConfig)
	if err != nil {
		klog.Fatalf("Error building sample controllers: %s", err.Error())
	}

	if err = crds(ctx, restConfig); err != nil {
		klog.Fatalf("Error apply CRDs: %s", err.Error())
	}

	objectSetApply := apply.New(clientSet.DiscoveryClient, apply.NewClientFactory(restConfig))

	foocontroller.Register(ctx, objectSetApply, foos.Edgeapi().V1alpha1().Foo())

	if err := start.All(ctx, threadiness, foos); err != nil {
		klog.Fatalf("Error starting: %s", err.Error())
	}
	return nil
}

func crds(ctx context.Context, config *rest.Config) error {
	factory, err := crd.NewFactoryFromClient(config)
	if err != nil {
		return err
	}

	factory.BatchCreateCRDs(ctx, crd.NonNamespacedTypes(
		edgev1Schema.SetAndGetCRDName("Foo"))...)

	return factory.BatchWait()
}
