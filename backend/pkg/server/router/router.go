package router

import (
	"net/http"

	"github.com/cnrancher/edge-ui/backend/pkg/auth"
	"github.com/cnrancher/edge-ui/backend/pkg/extendapi"
	"github.com/cnrancher/edge-ui/backend/pkg/server"
	"github.com/gorilla/mux"
	"github.com/rancher/steve/pkg/responsewriter"
)

func New(svr *server.EdgeServer) http.Handler {
	r := mux.NewRouter()
	r.UseEncodedPath()
	r.Use(responsewriter.ContentTypeOptions)

	r.Path("/v2/login").Handler(auth.NewLoginHandler(svr.RestConfig.Host))
	r.Path("/v2/health/datastorage").Handler(extendapi.NewDataStorgeHealthHandler(svr.Client))
	return r
}
