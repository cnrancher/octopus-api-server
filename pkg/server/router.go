package server

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/cnrancher/edge-api-server/pkg/auth"
	"github.com/cnrancher/edge-api-server/pkg/extendapi"
)

func SetupLocalHandler(server *EdgeServer, a auth.Authenticator) http.Handler {
	r := mux.NewRouter()
	r.UseEncodedPath()

	authHandler := auth.NewAuthHandler(server.RestConfig.Host, server.ClientSet, server.Context)
	dataStorageHealthHandler := extendapi.NewDataStorgeHealthHandler(server.ClientSet)
	deviceTemplateHandler := extendapi.NewDeviceTemplateHandler(server.ClientSet, server.DynamicClient, a)

	r.Path("/v2-public/auth").Handler(authHandler)
	r.Path("/v2-public/health/datastorage").Handler(dataStorageHealthHandler)
	r.Path("/v2/devicetemplates").Handler(deviceTemplateHandler)

	return r
}
