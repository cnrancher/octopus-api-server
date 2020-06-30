module github.com/cnrancher/edge-api-server

go 1.13

replace (
	github.com/crewjam/saml => github.com/rancher/saml v0.0.0-20180713225824-ce1532152fde
	github.com/rancher/apiserver => github.com/rancheredge/apiserver v0.0.0-20200629133558-61fb8ecfe2b8
	github.com/rancher/steve => github.com/rancheredge/steve v0.0.0-20200629092210-84d820c55b5b
	k8s.io/client-go => k8s.io/client-go v0.18.0
)

require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/gorilla/mux v1.7.3
	github.com/pkg/errors v0.9.1
	github.com/rancher/apiserver v0.0.0-20200612212259-10457317eb0b
	github.com/rancher/lasso v0.0.0-20200515155337-a34e1e26ad91
	github.com/rancher/steve v0.0.0-20200612212358-02b060294531
	github.com/rancher/wrangler v0.6.2-0.20200515155908-1923f3f8ec3f
	github.com/rancher/wrangler-api v0.6.1-0.20200515193802-dcf70881b087
	github.com/sirupsen/logrus v1.4.2
	github.com/urfave/cli v1.22.2
	gopkg.in/yaml.v2 v2.2.8
	k8s.io/api v0.18.0
	k8s.io/apimachinery v0.18.0
	k8s.io/apiserver v0.18.0
	k8s.io/client-go v12.0.0+incompatible
	k8s.io/klog v1.0.0
)
