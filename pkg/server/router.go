package server

import (
	"net/http"

	"github.com/rancher/steve/pkg/responsewriter"

	"github.com/cnrancher/edge-api-server/pkg/server/ui"

	"github.com/cnrancher/edge-api-server/pkg/auth"
	"github.com/cnrancher/edge-api-server/pkg/extendapi"
	"github.com/gorilla/mux"
)

func SetupLocalHandler(server *EdgeServer) http.Handler {
	r := mux.NewRouter()
	r.UseEncodedPath()

	authHandler := auth.NewAuthHandler(server.Context, server.RestConfig.Host, server.ClientSet)
	dataStorageHealthHandler := extendapi.NewDataStorgeHealthHandler(server.ClientSet)

	r.Path("/v2-public/auth").Handler(authHandler)
	r.Path("/v2-public/health/datastorage").Handler(dataStorageHealthHandler)

	//API UI
	uiContent := responsewriter.NewMiddlewareChain(responsewriter.Gzip, responsewriter.DenyFrameOptions,
		responsewriter.CacheMiddleware("json", "js", "css")).Handler(ui.Content())

	r.PathPrefix("/api-ui").Handler(uiContent)
	r.NotFoundHandler = ui.UI(http.NotFoundHandler())

	return r
}
