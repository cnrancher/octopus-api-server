package framework

import (
	"os"

	"github.com/cnrancher/octopus-api-server/tests/integration/cluster"
)

var (
	kubeConfig string
)

func init() {
	path := os.Getenv("HOME")
	kubeConfig = path + "/.config/k3d/" + cluster.ClusterName + "/kubeconfig.yaml"
}
