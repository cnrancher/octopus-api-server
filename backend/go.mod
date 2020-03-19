module github.com/cnrancher/edge-ui/backend

go 1.13

replace (
	github.com/containerd/cri => github.com/rancher/cri v1.3.0-k3s.3
	github.com/kubernetes-sigs/cri-tools => github.com/rancher/cri-tools v1.17.0-k3s1
	k8s.io/api => github.com/rancher/kubernetes/staging/src/k8s.io/api v1.17.2-k3s1
	k8s.io/apiextensions-apiserver => github.com/rancher/kubernetes/staging/src/k8s.io/apiextensions-apiserver v1.17.2-k3s1
	k8s.io/apimachinery => github.com/rancher/kubernetes/staging/src/k8s.io/apimachinery v1.17.2-k3s1
	k8s.io/apiserver => github.com/rancher/kubernetes/staging/src/k8s.io/apiserver v1.17.2-k3s1
	k8s.io/cli-runtime => github.com/rancher/kubernetes/staging/src/k8s.io/cli-runtime v1.17.4-k3s1
	k8s.io/client-go => k8s.io/client-go v0.17.2
	k8s.io/cloud-provider => github.com/rancher/kubernetes/staging/src/k8s.io/cloud-provider v1.17.2-k3s1
	k8s.io/cluster-bootstrap => github.com/rancher/kubernetes/staging/src/k8s.io/cluster-bootstrap v1.17.4-k3s1
	k8s.io/code-generator => github.com/rancher/kubernetes/staging/src/k8s.io/code-generator v1.17.4-k3s1
	k8s.io/component-base => github.com/rancher/kubernetes/staging/src/k8s.io/component-base v1.17.2-k3s1
	k8s.io/cri-api => github.com/rancher/kubernetes/staging/src/k8s.io/cri-api v1.17.2-k3s1
	k8s.io/csi-translation-lib => github.com/rancher/kubernetes/staging/src/k8s.io/csi-translation-lib v1.17.4-k3s1
	k8s.io/kube-aggregator => github.com/rancher/kubernetes/staging/src/k8s.io/kube-aggregator v1.17.4-k3s1
	k8s.io/kube-controller-manager => github.com/rancher/kubernetes/staging/src/k8s.io/kube-controller-manager v1.17.4-k3s1
	k8s.io/kube-proxy => github.com/rancher/kubernetes/staging/src/k8s.io/kube-proxy v1.17.4-k3s1
	k8s.io/kube-scheduler => github.com/rancher/kubernetes/staging/src/k8s.io/kube-scheduler v1.17.4-k3s1
	k8s.io/kubectl => github.com/rancher/kubernetes/staging/src/k8s.io/kubectl v1.17.2-k3s1
	k8s.io/kubelet => github.com/rancher/kubernetes/staging/src/k8s.io/kubelet v1.17.4-k3s1
	k8s.io/kubernetes => github.com/rancher/kubernetes v1.17.2-k3s1
	k8s.io/legacy-cloud-providers => github.com/rancher/kubernetes/staging/src/k8s.io/legacy-cloud-providers v1.17.4-k3s1
	k8s.io/metrics => github.com/rancher/kubernetes/staging/src/k8s.io/metrics v1.17.4-k3s1
	k8s.io/sample-apiserver => github.com/rancher/kubernetes/staging/src/k8s.io/sample-apiserver v1.17.4-k3s1
)

require (
	github.com/rancher/k3s v1.0.1
	github.com/rancher/steve v0.0.0-20200302052336-6dafed731f41
	github.com/rancher/wrangler v0.4.2-0.20200215064225-8abf292acf7b
	github.com/sirupsen/logrus v1.4.2
	github.com/urfave/cli v1.22.2
	k8s.io/apiserver v0.17.2
)
