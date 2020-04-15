module github.com/cnrancher/edge-api-server

go 1.13

replace k8s.io/client-go => k8s.io/client-go v0.18.0

require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/gorilla/mux v1.6.2
	github.com/pkg/errors v0.8.1
	github.com/rancher/steve v0.0.0-20200326205851-420f62f642eb
	github.com/rancher/wrangler v0.6.0
	github.com/sirupsen/logrus v1.4.2
	github.com/urfave/cli v1.22.2
	k8s.io/api v0.18.0
	k8s.io/apimachinery v0.18.0
	k8s.io/apiserver v0.18.0
	k8s.io/client-go v11.0.1-0.20190409021438-1a26190bd76a+incompatible
)
