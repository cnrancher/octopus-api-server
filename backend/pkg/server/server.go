package server

import (
	"context"

	"k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
)

type EdgeServer struct {
	RestConfig *restclient.Config
	Client     *kubernetes.Clientset
	Context    context.Context
}
