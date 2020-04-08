//go:generate go run pkg/codegen/cleanup/main.go
//go:generate /bin/rm -rf pkg/generated
//go:generate go run pkg/codegen/main.go

package main

import (
	"context"
	"flag"
	"github.com/cnrancher/edge-ui/backend/pkg/auth"
	"github.com/cnrancher/edge-ui/backend/pkg/server/router"
	"os"

	"github.com/rancher/steve/pkg/debug"
	steveserver "github.com/rancher/steve/pkg/server"
	stevecli "github.com/rancher/steve/pkg/server/cli"
	"github.com/rancher/steve/pkg/version"
	"github.com/rancher/wrangler/pkg/kubeconfig"
	"github.com/rancher/wrangler/pkg/ratelimit"
	"github.com/rancher/wrangler/pkg/signals"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var (
	Version     = "v0.0.1"
	GitCommit   = "HEAD"
	KubeConfig  string
	steveConfig stevecli.Config
	debugConfig debug.Config
)

func main() {
	app := cli.NewApp()
	app.Name = "edge-api-controller"
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

func run(_ *cli.Context) error {
	flag.Parse()

	logrus.Info("Starting controller")
	ctx := signals.SetupSignalHandler(context.Background())

	debugConfig.MustSetupDebug()

	s, err := newSteveServer(steveConfig)
	if err != nil {
		return err
	}
	return s.ListenAndServe(ctx, steveConfig.HTTPSListenPort, steveConfig.HTTPListenPort, nil)
}

func newSteveServer(c stevecli.Config) (*steveserver.Server, error) {
	restConfig, err := kubeconfig.GetInteractiveClientConfig(c.KubeConfig).ClientConfig()
	if err != nil {
		return nil, err
	}

	restConfig.RateLimiter = ratelimit.None
	a := auth.NewK3sAuthenticator(restConfig.Host)
	handler := router.New(restConfig)

	return &steveserver.Server{
		RestConfig: restConfig,
		AuthMiddleware: auth.ToAuthMiddleware(a),
		DashboardURL: func() string {
			return c.DashboardURL
		},
		Next: handler,
	}, nil
}
