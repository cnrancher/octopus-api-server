//go:generate go run pkg/codegen/cleanup/main.go
//go:generate /bin/rm -rf pkg/generated
//go:generate go run pkg/codegen/main.go

package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/cnrancher/edge-api-server/pkg/auth"
	edgeserver "github.com/cnrancher/edge-api-server/pkg/server"
	"github.com/cnrancher/edge-api-server/pkg/server/router"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

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
	app.Name = "edge-api-server"
	app.Version = version.FriendlyVersion()
	app.Usage = "run k3s edge UI api server!"
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

	s, err := newSteveServer(steveConfig, ctx)
	if err != nil {
		return err
	}
	return s.ListenAndServe(ctx, steveConfig.HTTPSListenPort, steveConfig.HTTPListenPort, nil)
}

func initKubeClient(kubeconfig string) (*kubernetes.Clientset, error) {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, fmt.Errorf("kubeconfig error %s\n", err.Error())
	}

	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("kubernetes clientset create error: %s", err.Error())
	}

	return clientSet, nil
}

func newSteveServer(c stevecli.Config, ctx context.Context) (*steveserver.Server, error) {
	restConfig, err := kubeconfig.GetInteractiveClientConfig(c.KubeConfig).ClientConfig()
	if err != nil {
		return nil, err
	}

	client, err := initKubeClient(c.KubeConfig)
	if err != nil {
		return nil, err
	}

	restConfig.RateLimiter = ratelimit.None

	a := auth.NewK3sAuthenticator(restConfig.Host, client, ctx)
	edgeServer := &edgeserver.EdgeServer{
		RestConfig: restConfig,
		Client:     client,
		Context:    ctx,
	}

	handler := router.New(edgeServer)

	return &steveserver.Server{
		RestConfig:     restConfig,
		AuthMiddleware: auth.ToAuthMiddleware(a),
		DashboardURL: func() string {
			return c.DashboardURL
		},
		Next: handler,
	}, nil
}
