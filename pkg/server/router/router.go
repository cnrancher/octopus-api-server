package router

import (
	"net/http"

	"github.com/cnrancher/edge-api-server/pkg/auth"
	"github.com/cnrancher/edge-api-server/pkg/extendapi"
	"github.com/cnrancher/edge-api-server/pkg/server"
	"github.com/gorilla/mux"
	"github.com/rancher/steve/pkg/responsewriter"
)

func New(svr *server.EdgeServer) http.Handler {
	r := mux.NewRouter()
	r.UseEncodedPath()
	r.Use(responsewriter.ContentTypeOptions)

	r.Path("/v2-public/auth").Handler(auth.NewAuthHandler(svr.RestConfig.Host, svr.Client, svr.Context))
	r.Path("/v2-public/health/datastorage").Handler(extendapi.NewDataStorgeHealthHandler(svr.Client))
	return r
}
