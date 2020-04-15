package main

import (
	"os"

	controllergen "github.com/rancher/wrangler/pkg/controller-gen"
	"github.com/rancher/wrangler/pkg/controller-gen/args"
)

func main() {
	os.Unsetenv("GOPATH")
	controllergen.Run(args.Options{
		OutputPackage: "github.com/cnrancher/edge-api-server/pkg/generated",
		Boilerplate:   "scripts/boilerplate.go.txt",
		Groups:        map[string]args.Group{},
	})
}
