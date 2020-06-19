package framework

import (
	"os"
)

var (
	kubeConfig  string
)

func init(){
	kubeConfig = os.Getenv("KUBECONFIG")
}