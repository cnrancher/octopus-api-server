package framework

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/json"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type K8sCli struct {
	Clientset *kubernetes.Clientset
	Dyclient  dynamic.Interface
	Cfg       *rest.Config
}

func NewK8sCli() (*K8sCli, error) {
	config, err := clientcmd.BuildConfigFromFlags("", kubeConfig)
	if err != nil {
		return nil, err
	}
	k8s := K8sCli{}
	k8s.Clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	k8s.Cfg = config
	dyClient, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	k8s.Dyclient = dyClient
	return &k8s, nil
}

func ConvertToUnstruct(in metav1.Object) (*unstructured.Unstructured, error) {
	uns := &unstructured.Unstructured{
		Object: make(map[string]interface{}),
	}
	data, err := json.Marshal(in)
	if err != nil {
		return uns, err
	}
	if err = uns.UnmarshalJSON(data); err != nil {
		return uns, err
	}
	return uns, nil
}
