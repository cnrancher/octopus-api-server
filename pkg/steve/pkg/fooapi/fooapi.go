package fooapi

import (
	"context"
	"net/http"

	"github.com/rancher/steve/pkg/server/store/proxy"

	"github.com/rancher/steve/pkg/schema"

	"github.com/cnrancher/edge-api-server/pkg/steve/customization"
	"github.com/rancher/steve/pkg/accesscontrol"
	"github.com/rancher/steve/pkg/auth"
	"github.com/rancher/steve/pkg/client"
	"github.com/rancher/steve/pkg/schemaserver/types"
	steveserver "github.com/rancher/steve/pkg/server"
	"github.com/rancher/wrangler/pkg/schemas"
)

type Server struct {
	ctx     context.Context
	asl     accesscontrol.AccessSetLookup
	auth    auth.Middleware
	cf      *client.Factory
	schemas *types.APISchemas
}

func (s *Server) Setup(ctx context.Context, server *steveserver.Server) error {
	s.ctx = ctx
	s.asl = server.AccessSetLookup
	s.auth = server.AuthMiddleware
	s.cf = server.ClientFactory
	s.schemas = server.BaseSchemas

	store := proxy.NewProxyStore(s.cf, s.asl)
	refreshHandler := customization.NewFooHandler()

	server.SchemaTemplates = append(server.SchemaTemplates, schema.Template{
		Store: Wrap(store),
		ID:    "edgeapi.cattle.io.foo",
		Formatter: func(request *types.APIRequest, resource *types.RawResource) {
			resource.AddAction(request, "refresh")
			resource.Schema.ActionHandlers = map[string]http.Handler{
				"refresh": refreshHandler,
			}
		},
		Customize: func(schema *types.APISchema) {
			schema.ResourceActions = map[string]schemas.Action{
				"refresh": {Output: "refresh"},
			}
			schema.CollectionActions = map[string]schemas.Action{
				"refresh": {Output: "refresh"},
			}
			schema.ActionHandlers = map[string]http.Handler{
				"refresh": refreshHandler,
			}
			schema.Schema.ResourceMethods = []string{"POST"}
		},
	})
	return nil
}
