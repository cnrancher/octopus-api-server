package server

import (
	"net/http"

	"github.com/cnrancher/octopus-api-server/pkg/auth"
	"github.com/cnrancher/octopus-api-server/pkg/extendapi"
	"github.com/cnrancher/octopus-api-server/pkg/server/ui"
	"github.com/gorilla/mux"
	"github.com/rancher/apiserver/pkg/middleware"
)

func SetupLocalHandler(server *EdgeServer) http.Handler {
	r := mux.NewRouter()
	r.UseEncodedPath()

	authHandler := auth.NewAuthHandler(server.Context, server.RestConfig.Host, server.ClientSet)
	dataStorageHealthHandler := extendapi.NewDataStorgeHealthHandler(server.ClientSet)

	r.Path("/").Handler(http.RedirectHandler("/dashboard", http.StatusTemporaryRedirect))
	r.Path("/v2-public/auth").Handler(authHandler)
	r.Path("/v2-public/health/datastorage").Handler(dataStorageHealthHandler)

	//API UI
	uiContent := middleware.NewMiddlewareChain(middleware.Gzip, middleware.DenyFrameOptions,
		middleware.CacheMiddleware("json", "js", "css")).Handler(ui.Content())

	r.PathPrefix("/api-ui").Handler(uiContent)
	r.NotFoundHandler = http.NotFoundHandler()
	return r
}
