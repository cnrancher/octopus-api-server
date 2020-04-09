package extendapi

import (
	"context"
	"net/http"

	"k8s.io/client-go/kubernetes"
)

type DataStorageHealthHandler struct {
	clientset *kubernetes.Clientset
}

func NewDataStorgeHealthHandler(client *kubernetes.Clientset) *DataStorageHealthHandler {
	return &DataStorageHealthHandler{
		clientset: client,
	}
}

// /v2/health/datastorage
func (h *DataStorageHealthHandler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	var code int
	h.clientset.RESTClient().Get().AbsPath("/healthz/etcd").Do(context.TODO()).StatusCode(&code)

	if code == http.StatusOK {
		writer.Write([]byte(`{"health":true}`))
	} else {
		writer.Write([]byte(`{"health":false}`))
	}
}
