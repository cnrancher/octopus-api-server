package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/cnrancher/octopus-api-server/pkg/auth"
	"github.com/cnrancher/octopus-api-server/pkg/controllers"
	"github.com/cnrancher/octopus-api-server/pkg/server/ui"
	"github.com/cnrancher/octopus-api-server/pkg/settings"
	"github.com/cnrancher/octopus-api-server/pkg/steve/pkg/catalogapi"
	"github.com/cnrancher/octopus-api-server/pkg/steve/pkg/devicetemplateapi"
	"github.com/cnrancher/octopus-api-server/pkg/steve/pkg/devicetemplaterevisionapi"

	"github.com/rancher/apiserver/pkg/writer"
	"github.com/rancher/steve/pkg/accesscontrol"
	"github.com/rancher/steve/pkg/server"
	"github.com/rancher/wrangler/pkg/ratelimit"
	"github.com/sirupsen/logrus"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type Config struct {
	Namespace       string
	Threadiness     int
	HTTPListenPort  int
	HTTPSListenPort int
}

type EdgeServer struct {
	*server.Controllers

	Config        Config
	RestConfig    *restclient.Config
	DynamicClient dynamic.Interface
	ClientSet     *kubernetes.Clientset
	Context       context.Context
	Handler       http.Handler
	ASL           accesscontrol.AccessSetLookup
}

func (s *EdgeServer) ListenAndServe(ctx context.Context) error {
	server, err := newSteveServer(ctx, s)
	if err != nil {
		return err
	}
	return server.ListenAndServe(ctx, s.Config.HTTPSListenPort, s.Config.HTTPListenPort, nil)
}

func New(ctx context.Context, clientConfig clientcmd.ClientConfig, cfg *Config) (*EdgeServer, error) {
	restConfig, err := clientConfig.ClientConfig()
	if err != nil {
		return nil, err
	}
	restConfig.RateLimiter = ratelimit.None

	if err := Wait(ctx, *restConfig); err != nil {
		return nil, err
	}

	clientSet, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return nil, fmt.Errorf("kubernetes clientset create error: %s", err.Error())
	}

	dynamicClient, err := dynamic.NewForConfig(restConfig)
	if err != nil {
		return nil, fmt.Errorf("kubernetes dynamic client create error:%s", err.Error())
	}

	steveControllers, err := server.NewController(restConfig, nil)
	if err != nil {
		return nil, err
	}

	asl := accesscontrol.NewAccessStore(ctx, true, steveControllers.RBAC)

	err = controllers.Setup(ctx, restConfig, clientSet, 5)
	if err != nil {
		return nil, err
	}

	return &EdgeServer{
		Controllers:   steveControllers,
		Config:        *cfg,
		Context:       ctx,
		ClientSet:     clientSet,
		DynamicClient: dynamicClient,
		RestConfig:    restConfig,
		ASL:           asl,
	}, nil
}

func newSteveServer(ctx context.Context, edgeServer *EdgeServer) (*server.Server, error) {
	a := auth.NewK3sAuthenticator(ctx, edgeServer.RestConfig.Host, edgeServer.ClientSet)
	handler := SetupLocalHandler(edgeServer)

	catalogAPIServer := &catalogapi.Server{}
	deviceTemplateAPIServer := &devicetemplateapi.Server{Authenticator: a}
	deviceTemplateRevisionAPIServer := &devicetemplaterevisionapi.Server{Authenticator: a}
	return &server.Server{
		Controllers:     edgeServer.Controllers,
		AccessSetLookup: edgeServer.ASL,
		RestConfig:      edgeServer.RestConfig,
		AuthMiddleware:  auth.ToAuthMiddleware(a),
		Next:            handler,
		DashboardURL: func() string {
			if settings.UIIndex.Get() == "local" {
				return settings.UIPath.Get()
			}
			return settings.UIIndex.Get()
		},
		StartHooks: []server.StartHook{
			catalogAPIServer.Setup,
			deviceTemplateAPIServer.Setup,
			deviceTemplateRevisionAPIServer.Setup,
		},
		HTMLResponseWriter: writer.HTMLResponseWriter{
			CSSURL:       ui.CSSURLGetter,
			JSURL:        ui.JSURLGetter,
			APIUIVersion: ui.APIUIVersionGetter,
		},
	}, nil
}

func Wait(ctx context.Context, config rest.Config) error {
	client, err := kubernetes.NewForConfig(&config)
	if err != nil {
		return err
	}

	for {
		_, err := client.Discovery().ServerVersion()
		if err == nil {
			break
		}
		logrus.Infof("Waiting for server to become available: %v", err)
		select {
		case <-ctx.Done():
			return fmt.Errorf("startup canceled")
		case <-time.After(2 * time.Second):
		}
	}

	return nil
}
