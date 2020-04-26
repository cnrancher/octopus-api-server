package customization

import (
	"net/http"
	"strings"

	"github.com/cnrancher/edge-api-server/pkg/generated/controllers/edgeapi.cattle.io/v1alpha1"
)

type FooHandler struct {
	FooClient v1alpha1.FooClient
}

func NewFooHandler() *FooHandler {
	return &FooHandler{}
}

func (h *FooHandler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	action := strings.ToLower(req.URL.Query().Get("action"))
	writer.Write([]byte("make action call:" + action))
}
