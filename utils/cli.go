package utils

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog/v2"
)

type KubernetesClient struct {
	Clientset  *kubernetes.Clientset
	RestConfig *rest.Config
}

func NewKubernetesClient() *KubernetesClient {
	cli := &KubernetesClient{}
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	cli.RestConfig = config
	cli.Clientset = clientSet
	return cli
}

func NewKubernetesClientOutSide(kubeconfig string) *KubernetesClient {
	cli := &KubernetesClient{}

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	restConfig, err := NewRestConfig(kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	cli.Clientset =  clientset
	cli.RestConfig = restConfig
	return cli
}

func NewRestConfig(kubeConfig string) (*rest.Config, error) {
	clientConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{
			ExplicitPath: kubeConfig,
		},
		&clientcmd.ConfigOverrides{},
	)
	clusterConfig, err := clientConfig.ClientConfig()
	if err != nil {
		klog.Warningf("create new k8s rest client for %s got error:%s", kubeConfig, err)
		return nil, err
	}

	//clusterConfig.ContentConfig.GroupVersion = &v1alpha1.SchemeGroupVersion
	clusterConfig.APIPath = "/apis"
	//clusterConfig.NegotiatedSerializer = serializer.DirectCodecFactory{CodecFactory: scheme.Codecs}
	clusterConfig.UserAgent = rest.DefaultKubernetesUserAgent()
	return clusterConfig, nil
}