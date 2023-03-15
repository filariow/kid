package ksa

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	mv1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	av1 "k8s.io/client-go/applyconfigurations/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

var ErrSecretNotFound = fmt.Errorf("service account's secret not found")
var ErrSecretMalformed = fmt.Errorf("service account's secret malformed")

type ServiceAccountToken struct {
	CACrt     []byte `json:"ca.crt"`
	Namespace []byte `json:"namespace"`
	Token     []byte `json:"token"`
}

type GetKubeconfigOptions struct {
	OverrideHost *string
	User         *string
	Namespace    *string
}

func GetKubeconfig(cli kubernetes.Clientset, token *ServiceAccountToken, opts GetKubeconfigOptions) ([]byte, error) {
	cfg, err := GetRESTConfig()
	if err != nil {
		return nil, err
	}

	h := cfg.Host
	if opts.OverrideHost != nil {
		h = *opts.OverrideHost
	}

	cl := map[string]*clientcmdapi.Cluster{
		cfg.ServerName: {
			Server:                   h,
			CertificateAuthorityData: token.CACrt,
		},
	}

	un := "default-user"
	if opts.User != nil {
		un = *opts.User
	}
	ai := map[string]*clientcmdapi.AuthInfo{
		un: {
			Token: string(token.Token),
		},
	}

	ct := map[string]*clientcmdapi.Context{
		"default-context": {
			Cluster:  cfg.ServerName,
			AuthInfo: un,
		},
	}
	if opts.Namespace != nil {
		ct["default-context"].Namespace = *opts.Namespace
	}

	cc := clientcmdapi.Config{
		Kind:       "Config",
		APIVersion: "v1",
		Clusters:   cl,
		Contexts:   ct,
		AuthInfos:  ai,
	}

	return clientcmd.Write(cc)
}

func GetToken(secret *corev1.Secret) (*ServiceAccountToken, error) {
	c, err := getSecretDataField(secret, "ca.crt")
	if err != nil {
		return nil, err
	}

	n, err := getSecretDataField(secret, "namespace")
	if err != nil {
		return nil, err
	}

	t, err := getSecretDataField(secret, "token")
	if err != nil {
		return nil, err
	}

	return &ServiceAccountToken{
		CACrt:     c,
		Namespace: n,
		Token:     t,
	}, nil
}

func getSecretDataField(secret *corev1.Secret, field string) ([]byte, error) {
	d, ok := secret.Data[field]
	if !ok {
		return nil, fmt.Errorf("%w: can not find '%s' in secret '%s/%s'", ErrSecretMalformed, field, secret.Namespace, secret.Name)
	}

	return d, nil
}

func GetLatestServiceAccountSecrets(ctx context.Context, cli kubernetes.Clientset, name string, namespace string) (*corev1.Secret, error) {
	ss, err := GetServiceAccountSecrets(ctx, cli, name, namespace)
	if err != nil {
		return nil, err
	}

	if len(ss) == 0 {
		return nil, ErrSecretNotFound
	}

	fs := ss[0]
	for _, s := range ss[1:] {
		if s.GetCreationTimestamp().Compare(fs.GetCreationTimestamp().Time) < 0 {
			fs = s
		}
	}

	return &fs, nil
}

func GetServiceAccountSecrets(ctx context.Context, cli kubernetes.Clientset, name string, namespace string) ([]corev1.Secret, error) {
	ss, err := cli.CoreV1().Secrets(namespace).List(ctx, mv1.ListOptions{})
	if err != nil {
		return nil, err
	}

	fss := []corev1.Secret{}
	for _, s := range ss.Items {
		if n, ok := s.Annotations[corev1.ServiceAccountNameKey]; ok && n == name {
			fss = append(fss, s)
		}
	}
	return fss, nil
}

func CreateServiceAccountSecret(ctx context.Context, cli kubernetes.Clientset, name string, namespace string, saname string) (*corev1.Secret, error) {
	s := av1.Secret(name, namespace)
	s.WithType(corev1.SecretTypeServiceAccountToken)
	s.WithAnnotations(map[string]string{
		corev1.ServiceAccountNameKey: saname,
	})

	o := mv1.ApplyOptions{FieldManager: "application/apply-patch"}
	return cli.CoreV1().Secrets(namespace).Apply(ctx, s, o)
}

func CreateServiceAccount(ctx context.Context, cli kubernetes.Clientset, name string, namespace string) (*corev1.ServiceAccount, error) {
	c := av1.ServiceAccount(name, namespace)

	o := mv1.ApplyOptions{FieldManager: "application/apply-patch"}
	return cli.CoreV1().ServiceAccounts(namespace).Apply(ctx, c, o)
}
