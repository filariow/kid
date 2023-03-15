/*
Copyright © 2023 Francesco Ilario

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/

package ksa

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	mv1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

func GetLastServiceAccountSecrets(ctx context.Context, cli kubernetes.Clientset, name string, namespace string) (*corev1.Secret, error) {
	ss, err := GetServiceAccountSecrets(ctx, cli, name, namespace)
	if err != nil {
		return nil, err
	}

	if len(ss) == 0 {
		return nil, ErrSecretNotFound
	}

	fs := ss[0]
	for _, s := range ss[1:] {
		if s.GetCreationTimestamp().Compare(fs.GetCreationTimestamp().Time) >= 0 {
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

func CreateServiceAccountSecret(ctx context.Context, cli kubernetes.Clientset, name string, namespace string, sa *corev1.ServiceAccount) (*corev1.Secret, error) {
	s := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Annotations: map[string]string{
				corev1.ServiceAccountNameKey: sa.Name,
			},
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion: "v1",
					Kind:       "ServiceAccount",
					Name:       sa.Name,
					UID:        sa.UID,
				},
			},
		},
		Type: corev1.SecretTypeServiceAccountToken,
	}

	return cli.CoreV1().Secrets(namespace).Create(ctx, s, mv1.CreateOptions{})
}

func DeleteServiceAccountSecret(ctx context.Context, cli kubernetes.Clientset, name string, namespace string) error {
	return cli.CoreV1().Secrets(namespace).Delete(ctx, name, mv1.DeleteOptions{})
}

func CreateServiceAccount(ctx context.Context, cli kubernetes.Clientset, name string, namespace string) (*corev1.ServiceAccount, error) {
	c := &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
	}

	o := mv1.CreateOptions{}
	return cli.CoreV1().ServiceAccounts(namespace).Create(ctx, c, o)
}
