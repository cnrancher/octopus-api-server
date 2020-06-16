package devicetemplaterevisionapi

import (
	"context"

	"github.com/rancher/steve/pkg/schema"
	"github.com/rancher/steve/pkg/server/store/proxy"
	"github.com/sirupsen/logrus"

	apiAuth "github.com/cnrancher/edge-api-server/pkg/auth"
	v1 "github.com/cnrancher/edge-api-server/pkg/generated/controllers/edgeapi.cattle.io"
	"github.com/rancher/steve/pkg/accesscontrol"
	"github.com/rancher/steve/pkg/auth"
	"github.com/rancher/steve/pkg/client"
	"github.com/rancher/steve/pkg/schemaserver/types"
	steveserver "github.com/rancher/steve/pkg/server"
)

type Server struct {
	ctx           context.Context
	asl           accesscontrol.AccessSetLookup
	auth          auth.Middleware
	cf            *client.Factory
	schemas       *types.APISchemas
	Authenticator apiAuth.Authenticator
}

func (s *Server) Setup(ctx context.Context, server *steveserver.Server) error {
	s.ctx = ctx
	s.asl = server.AccessSetLookup
	s.auth = server.AuthMiddleware
	s.cf = server.ClientFactory
	s.schemas = server.BaseSchemas

	store := proxy.NewProxyStore(s.cf, s.asl)
	controllers, err := v1.NewFactoryFromConfig(server.RestConfig)
	if err != nil {
		logrus.Fatalf("Error building controllers: %s", err.Error())
	}
	store = &Store{
		Store:                    store,
		asl:                      s.asl,
		ctx:                      s.ctx,
		client:                   server.ClientFactory.DynamicClient(),
		auth:                     s.Authenticator,
		deviceTemplateController: controllers.Edgeapi().V1alpha1().DeviceTemplate(),
	}
	server.SchemaTemplates = append(server.SchemaTemplates, schema.Template{
		Store: store,
		ID:    "edgeapi.cattle.io.devicetemplaterevision",
	})
	return nil
}
