//go:generate go run pkg/codegen/cleanup/main.go
//go:generate /bin/rm -rf pkg/generated
//go:generate go run pkg/codegen/main.go

package main

import (
	"context"
	"flag"
	"os"

	"github.com/rancher/steve/pkg/debug"
	stevecli "github.com/rancher/steve/pkg/server/cli"
	"github.com/rancher/steve/pkg/version"
	"github.com/rancher/wrangler/pkg/signals"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var (
	Version    = "v0.0.1"
	GitCommit  = "HEAD"
	KubeConfig string
	steveConfig      stevecli.Config
	debugConfig debug.Config
)

func main() {
	app := cli.NewApp()
	app.Name = "edge-ui-backend"
	app.Version = version.FriendlyVersion()
	app.Usage = "run k3s edge UI api!"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "kubeconfig",
			EnvVar:      "KUBECONFIG",
			Destination: &KubeConfig,
		},
	}
	app.Flags = append(
		stevecli.Flags(&steveConfig),
		debug.Flags(&debugConfig)...)
	app.Action = run
	if err := app.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}
}

func run(_ *cli.Context) error{
	flag.Parse()

	logrus.Info("Starting controller")
	ctx := signals.SetupSignalHandler(context.Background())

	debugConfig.MustSetupDebug()
	s := steveConfig.MustServer()
	return s.ListenAndServe(ctx, steveConfig.HTTPSListenPort, steveConfig.HTTPListenPort, nil)
}
