package main

import (
	"os"

	"github.com/cnrancher/edge-api-server/pkg/apis/edgeapi.cattle.io/v1alpha1"

	controllergen "github.com/rancher/wrangler/pkg/controller-gen"
	"github.com/rancher/wrangler/pkg/controller-gen/args"
)

func main() {
	os.Unsetenv("GOPATH")
	controllergen.Run(args.Options{
		OutputPackage: "github.com/cnrancher/edge-api-server/pkg/generated",
		Boilerplate:   "scripts/boilerplate.go.txt",
		Groups: map[string]args.Group{
			"edgeapi.cattle.io": {
				Types: []interface{}{
					v1alpha1.Foo{},
				},
				GenerateTypes: true,
			},
		},
	})
}
