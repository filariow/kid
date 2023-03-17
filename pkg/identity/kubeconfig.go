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
	"github.com/filariow/kid/pkg/kube"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

type GetKubeconfigOptions struct {
	OverrideHost *string
	User         *string
	Namespace    *string
}

func GetKubeconfig(cli kubernetes.Clientset, token *ServiceAccountToken, opts GetKubeconfigOptions) ([]byte, error) {
	cfg, err := kube.GetRESTConfig()
	if err != nil {
		return nil, err
	}

	h := cfg.Host
	if opts.OverrideHost != nil {
		h = *opts.OverrideHost
	}

	cl := map[string]*clientcmdapi.Cluster{
		"default-cluster": {
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
			Cluster:  "default-cluster",
			AuthInfo: un,
		},
	}
	if opts.Namespace != nil {
		ct["default-context"].Namespace = *opts.Namespace
	}

	cc := clientcmdapi.Config{
		Kind:           "Config",
		APIVersion:     "v1",
		Clusters:       cl,
		Contexts:       ct,
		AuthInfos:      ai,
		CurrentContext: "default-context",
	}

	return clientcmd.Write(cc)
}
