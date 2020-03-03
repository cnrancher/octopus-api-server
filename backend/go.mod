module github.com/cnrancher/edge-ui/backend

go 1.13

replace k8s.io/client-go => k8s.io/client-go v0.17.2

require (
	github.com/rancher/steve v0.0.0-20200302052336-6dafed731f41
	github.com/rancher/wrangler v0.4.2-0.20200215064225-8abf292acf7b
	github.com/sirupsen/logrus v1.4.2
	github.com/urfave/cli v1.22.2
	k8s.io/apimachinery v0.17.2 // indirect
	k8s.io/client-go v11.0.1-0.20190409021438-1a26190bd76a+incompatible // indirect
)
