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

package identity

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/filariow/ksa/pkg/ksa"
	corev1 "k8s.io/api/core/v1"
	mv1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type Instance struct {
	Namespace      string   `json:"namespace"`
	ServiceAccount string   `json:"serviceAccount"`
	Secrets        []string `json:"secrets"`
}

func CreateIdentity(ctx context.Context, cli kubernetes.Clientset, name string, namespace string) (*Instance, error) {
	ss, err := ksa.GetServiceAccountSecrets(ctx, cli, name, namespace)
	if err != nil {
		return nil, err
	}
	if len(ss) > 0 {
		return nil, fmt.Errorf("error access tokens for Service Account '%s/%s' already exists", namespace, name)
	}

	sa, err := ksa.CreateServiceAccount(ctx, cli, name, namespace)
	if err != nil {
		return nil, err
	}

	sn := createSecretName(name, 1)
	s, err := ksa.CreateServiceAccountSecret(ctx, cli, sn, namespace, sa)
	if err != nil {
		return nil, err
	}

	return &Instance{
		Namespace:      namespace,
		ServiceAccount: sa.Name,
		Secrets:        []string{s.Name},
	}, nil
}

func BeginIdentityKeyRotation(ctx context.Context, cli kubernetes.Clientset, name string, namespace string) (*corev1.Secret, error) {
	sa, err := cli.CoreV1().ServiceAccounts(namespace).Get(ctx, name, mv1.GetOptions{})
	if err != nil {
		return nil, err
	}

	s, err := ksa.GetLastServiceAccountSecrets(ctx, cli, name, namespace)
	if err != nil {
		return nil, err
	}

	sn, err := nextSecretAccountSecretName(s.Name)
	if err != nil {
		return nil, err
	}

	return ksa.CreateServiceAccountSecret(ctx, cli, *sn, namespace, sa)
}

func RollbackIdentityKey(ctx context.Context, cli kubernetes.Clientset, name string, namespace string, version uint64) (*corev1.Secret, error) {
	sa, err := cli.CoreV1().ServiceAccounts(namespace).Get(ctx, name, mv1.GetOptions{})
	if err != nil {
		return nil, err
	}

	sn := createSecretName(name, version)
	return ksa.CreateServiceAccountSecret(ctx, cli, sn, namespace, sa)
}

func RevokeIdentityKey(ctx context.Context, cli kubernetes.Clientset, name string, namespace string, version uint64) (*string, error) {
	sn := createSecretName(name, version)
	if err := ksa.DeleteServiceAccountSecret(ctx, cli, sn, namespace); err != nil {
		return nil, err
	}
	return &sn, nil
}

func CompleteIdentityKeyRotation(ctx context.Context, cli kubernetes.Clientset, name string, namespace string) (*string, error) {
	s, err := ksa.GetLastServiceAccountSecrets(ctx, cli, name, namespace)
	if err != nil {
		return nil, err
	}

	sn, err := previousSecretAccountSecretName(s.Name)
	if err != nil {
		return nil, err
	}

	return sn, cli.CoreV1().Secrets(namespace).Delete(ctx, *sn, mv1.DeleteOptions{})
}

func createSecretName(sa string, version uint64) string {
	return fmt.Sprintf("%s-key-%d", sa, version)
}

func splitServiceAccountSecretName(name string) (string, uint64, error) {
	ss := strings.Split(name, "-")
	if len(ss) <= 1 {
		return "", 0, fmt.Errorf("invalid name for Service Account name: %s", name)
	}

	p := ss[len(ss)-1]
	u, err := strconv.ParseUint(p, 10, 64)
	if err != nil {
		return "", 0, err
	}

	b := strings.Join(ss[:len(ss)-1], "-")
	return b, u, nil
}

func nextSecretAccountSecretName(name string) (*string, error) {
	b, p, err := splitServiceAccountSecretName(name)
	if err != nil {
		return nil, err
	}

	p++
	n := fmt.Sprintf("%s-%d", b, p)

	return &n, nil
}

func previousSecretAccountSecretName(name string) (*string, error) {
	b, p, err := splitServiceAccountSecretName(name)
	if err != nil {
		return nil, err
	}

	if p == 0 {
		return nil, fmt.Errorf("no prior secret to %s", name)
	}

	p--
	n := fmt.Sprintf("%s-%d", b, p)

	return &n, nil
}
