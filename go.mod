module github.com/cnrancher/edge-api-server

go 1.13

replace k8s.io/client-go => k8s.io/client-go v0.18.0

require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/gorilla/mux v1.6.2
	github.com/pkg/errors v0.8.1
	github.com/rancher/lasso v0.0.0-20200417051414-b55b9620e2e7
	github.com/rancher/steve v0.0.0-20200417063946-685dea747a25
	github.com/rancher/wrangler v0.6.2-0.20200417063009-962aed6a55dc
	github.com/sirupsen/logrus v1.4.2
	github.com/urfave/cli v1.22.2
	gopkg.in/yaml.v2 v2.2.8
	k8s.io/api v0.18.0
	k8s.io/apimachinery v0.18.0
	k8s.io/apiserver v0.18.0
	k8s.io/client-go v12.0.0+incompatible
	k8s.io/klog v1.0.0
)
