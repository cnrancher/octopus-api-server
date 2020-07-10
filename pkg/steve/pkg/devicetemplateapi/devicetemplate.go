package devicetemplateapi

import (
	"context"

	apiAuth "github.com/cnrancher/octopus-api-server/pkg/auth"
	v1 "github.com/cnrancher/octopus-api-server/pkg/generated/controllers/octopusapi.cattle.io"
	"github.com/rancher/steve/pkg/accesscontrol"
	"github.com/rancher/steve/pkg/auth"
	"github.com/rancher/steve/pkg/client"
	"github.com/rancher/steve/pkg/schema"
	steveserver "github.com/rancher/steve/pkg/server"
	"github.com/rancher/steve/pkg/stores/proxy"
	"github.com/sirupsen/logrus"
)

type Server struct {
	ctx           context.Context
	asl           accesscontrol.AccessSetLookup
	auth          auth.Middleware
	cf            *client.Factory
	Authenticator apiAuth.Authenticator
}

func (s *Server) Setup(ctx context.Context, server *steveserver.Server) error {
	s.ctx = ctx
	s.asl = server.AccessSetLookup
	s.auth = server.AuthMiddleware
	s.cf = server.ClientFactory

	controllers, err := v1.NewFactoryFromConfig(server.RestConfig)
	if err != nil {
		logrus.Fatalf("Error building controllers: %s", err.Error())
	}

	proxyStore := proxy.NewProxyStore(s.cf, s.asl)
	store := &Store{
		Store:              proxyStore,
		asl:                s.asl,
		ctx:                s.ctx,
		auth:               s.Authenticator,
		revisionController: controllers.Octopusapi().V1alpha1().DeviceTemplateRevision(),
	}
	server.SchemaTemplates = append(server.SchemaTemplates, schema.Template{
		Store: store,
		ID:    "octopusapi.cattle.io.devicetemplate",
	})

	return nil
}
