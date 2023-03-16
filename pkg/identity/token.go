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
	"fmt"

	corev1 "k8s.io/api/core/v1"
)

var ErrSecretMalformed = fmt.Errorf("service account's secret malformed")

type ServiceAccountToken struct {
	CACrt     []byte `json:"ca.crt"`
	Namespace []byte `json:"namespace"`
	Token     []byte `json:"token"`
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
