package ksa

import (
	"os"
	"path"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func GetCurrentContextClient() (*kubernetes.Clientset, error) {
	cfg, err := GetRESTConfig()
	if err != nil {
		return nil, err
	}

	return kubernetes.NewForConfig(cfg)
}

func GetConfigDefaultNamespace() (*string, error) {
	cc, err := getClientConfig()
	if err != nil {
		return nil, err
	}
	ns, _, err := cc.Namespace()
	if err != nil {
		return nil, err
	}

	return &ns, err
}

func GetRESTConfig() (*rest.Config, error) {
	cc, err := getClientConfig()
	if err != nil {
		return nil, err
	}

	cfg, err := cc.ClientConfig()
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func getClientConfig() (clientcmd.ClientConfig, error) {
	kp, err := getKubeconfigPath()
	if err != nil {
		return nil, err
	}

	kd, err := os.ReadFile(*kp)
	if err != nil {
		return nil, err
	}

	return clientcmd.NewClientConfigFromBytes(kd)
}

func getKubeconfigPath() (*string, error) {
	if env := os.Getenv("KUBECONFIG"); env != "" {
		return &env, nil
	}

	hd, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	kp := path.Join(hd, clientcmd.RecommendedHomeDir, clientcmd.RecommendedFileName)
	return &kp, nil
}

func BuildClient(kfg []byte) (*kubernetes.Clientset, error) {
	cfg, err := clientcmd.RESTConfigFromKubeConfig(kfg)
	if err != nil {
		return nil, err
	}

	cli, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, err
	}

	return cli, nil
}
