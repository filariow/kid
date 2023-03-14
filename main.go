package main

import (
	"fmt"
	"log"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	fmt.Println(corev1.ServiceAccountNameKey)
	fmt.Println(corev1.SecretTypeServiceAccountToken)

	_, err := buildClient(nil)
	if err != nil {
		return err
	}

	return nil
}

func buildClient(kfg []byte) (*kubernetes.Clientset, error) {
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
