module github.com/cnrancher/octopus-api-server

go 1.13

replace (
	github.com/crewjam/saml => github.com/rancher/saml v0.0.0-20180713225824-ce1532152fde
	github.com/rancher/apiserver => github.com/rancheredge/apiserver v0.0.0-20200708064345-bd2aa8f9b7d1
	github.com/rancher/steve => github.com/rancheredge/steve v0.0.0-20200708031911-f69e0f4820b4
	k8s.io/client-go => k8s.io/client-go v0.18.0
)

require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/gorilla/mux v1.7.3
	github.com/pkg/errors v0.9.1
	github.com/rancher/apiserver v0.0.0-20200622174841-b4d1106a9883
	github.com/rancher/lasso v0.0.0-20200515155337-a34e1e26ad91
	github.com/rancher/steve v0.0.0-20200622175150-3dbc369174fb
	github.com/rancher/wrangler v0.6.2-0.20200515155908-1923f3f8ec3f
	github.com/rancher/wrangler-api v0.6.1-0.20200515193802-dcf70881b087
	github.com/sirupsen/logrus v1.4.2
	github.com/stretchr/testify v1.5.1
	github.com/urfave/cli v1.22.2
	gopkg.in/yaml.v2 v2.2.8
	k8s.io/api v0.18.0
	k8s.io/apimachinery v0.18.0
	k8s.io/apiserver v0.18.0
	k8s.io/client-go v12.0.0+incompatible
	k8s.io/klog v1.0.0
)
