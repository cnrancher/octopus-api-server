module github.com/cnrancher/edge-api-server

go 1.13

replace k8s.io/client-go => k8s.io/client-go v0.18.0

require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/gorilla/mux v1.6.2
	github.com/pkg/errors v0.9.1
	github.com/rancher/lasso v0.0.0-20200515155337-a34e1e26ad91
	github.com/rancher/steve v0.0.0-20200518163824-9e4ed62a470a
	github.com/rancher/wrangler v0.6.2-0.20200515155908-1923f3f8ec3f
	github.com/sirupsen/logrus v1.4.2
	github.com/urfave/cli v1.22.2
	gopkg.in/yaml.v2 v2.2.8
	k8s.io/api v0.18.0
	k8s.io/apimachinery v0.18.0
	k8s.io/apiserver v0.18.0
	k8s.io/client-go v12.0.0+incompatible
	k8s.io/klog v1.0.0
)
