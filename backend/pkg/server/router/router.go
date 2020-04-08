package router

import (
	"net/http"

	"github.com/cnrancher/edge-ui/backend/pkg/auth"
	"github.com/gorilla/mux"
	"github.com/rancher/steve/pkg/responsewriter"
	restclient "k8s.io/client-go/rest"
)

func New(restConfig *restclient.Config) http.Handler {
	r := mux.NewRouter()
	r.UseEncodedPath()
	r.Use(responsewriter.ContentTypeOptions)

	r.Path("/v2/login").Handler(auth.NewLoginHandler(restConfig.Host))
	return r
}
