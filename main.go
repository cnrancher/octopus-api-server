//go:generate go run pkg/codegen/cleanup/main.go
//go:generate /bin/rm -rf pkg/generated
//go:generate go run pkg/codegen/main.go

package main

import (
	"context"
	"os"

	"github.com/cnrancher/octopus-api-server/pkg/server"
	"github.com/rancher/steve/pkg/debug"
	"github.com/rancher/steve/pkg/version"
	"github.com/rancher/wrangler/pkg/signals"
	"github.com/urfave/cli"
	"k8s.io/klog"
)

var (
	kubeConfig  string
	debugConfig debug.Config
)

func main() {
	app := cli.NewApp()
	app.Name = "octopus-api-server"
	app.Version = version.FriendlyVersion()
	app.Usage = "Run the octopus api server of k3s"

	var config server.Config
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "kubeconfig",
			Usage:       "Kube config for accessing k8s cluster",
			EnvVar:      "KUBECONFIG",
			Destination: &kubeConfig,
		},
		cli.StringFlag{
			Name:        "namespace, n",
			EnvVar:      "NAMESPACE",
			Value:       "",
			Usage:       "Namespace to watch, empty means it will watch CRDs in all namespaces.",
			Destination: &config.Namespace,
		},
		cli.IntFlag{
			Name:        "threads, t",
			EnvVar:      "THREADS",
			Value:       5,
			Usage:       "Threadiness level to set, defaults to 5.",
			Destination: &config.Threadiness,
		},
		cli.IntFlag{
			Name:        "https-listen-port",
			Value:       8443,
			Destination: &config.HTTPSListenPort,
		},
		cli.IntFlag{
			Name:        "http-listen-port",
			Value:       8080,
			Destination: &config.HTTPListenPort,
		},
	}
	app.Flags = append(app.Flags, debug.Flags(&debugConfig)...)
	app.Action = func(c *cli.Context) error {
		return run(c, config)
	}
	if err := app.Run(os.Args); err != nil {
		klog.Fatal(err)
	}
}

func run(_ *cli.Context, config server.Config) error {
	debugConfig.MustSetupDebug()
	klog.Infof("Octopus api server version %s is starting", version.FriendlyVersion())

	ctx := signals.SetupSignalHandler(context.Background())

	clientConfig, err := server.GetConfig(kubeConfig)
	if err != nil {
		return err
	}

	s, err := server.New(ctx, clientConfig, &config)
	if err != nil {
		return err
	}
	return s.ListenAndServe(ctx)
}
