package catalogapi

import (
	"context"

	v1 "github.com/cnrancher/octopus-api-server/pkg/generated/controllers/octopusapi.cattle.io"
	"github.com/rancher/apiserver/pkg/types"
	"github.com/rancher/steve/pkg/accesscontrol"
	"github.com/rancher/steve/pkg/auth"
	"github.com/rancher/steve/pkg/client"
	"github.com/rancher/steve/pkg/resources/common"
	"github.com/rancher/steve/pkg/schema"
	steveserver "github.com/rancher/steve/pkg/server"
	"github.com/rancher/steve/pkg/stores/proxy"
	"github.com/rancher/wrangler/pkg/schemas"
	"github.com/sirupsen/logrus"
)

type Server struct {
	ctx  context.Context
	asl  accesscontrol.AccessSetLookup
	auth auth.Middleware
	cf   *client.Factory
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
		Store:      proxyStore,
		asl:        s.asl,
		controller: controllers.Octopusapi().V1alpha1().Catalog(),
	}
	server.SchemaTemplates = append(server.SchemaTemplates, schema.Template{
		Store: store,
		ID:    "octopusapi.cattle.io.catalog",
		Formatter: func(request *types.APIRequest, resource *types.RawResource) {
			common.Formatter(request, resource)
			resource.AddAction(request, "refresh")
		},
		Customize: func(schema *types.APISchema) {
			schema.ResourceActions = map[string]schemas.Action{
				"refresh": {Output: "refresh"},
			}
			schema.CollectionActions = map[string]schemas.Action{
				"refresh": {Output: "refresh"},
			}
			schema.Schema.ResourceMethods = []string{"POST"}
		},
	})
	return nil
}
