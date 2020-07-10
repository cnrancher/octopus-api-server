package main

import (
	"os"

	"github.com/cnrancher/octopus-api-server/pkg/apis/octopusapi.cattle.io/v1alpha1"
	controllergen "github.com/rancher/wrangler/pkg/controller-gen"
	"github.com/rancher/wrangler/pkg/controller-gen/args"
)

func main() {
	os.Unsetenv("GOPATH")
	controllergen.Run(args.Options{
		OutputPackage: "github.com/cnrancher/octopus-api-server/pkg/generated",
		Boilerplate:   "scripts/boilerplate.go.txt",
		Groups: map[string]args.Group{
			"octopusapi.cattle.io": {
				PackageName: "octopusapi.cattle.io",
				Types: []interface{}{
					v1alpha1.Catalog{},
					v1alpha1.DeviceTemplate{},
					v1alpha1.DeviceTemplateRevision{},
					v1alpha1.Setting{},
				},
				GenerateTypes:   true,
				GenerateClients: true,
			},
		},
	})
}
