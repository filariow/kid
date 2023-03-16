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

package cmd

import (
	"fmt"

	"github.com/filariow/kid/pkg/identity"
	"github.com/filariow/kid/pkg/kube"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const (
	getKubeconfigTargetNamespaceLongParam string = "target-namespace"
	getKubeconfigServerUrlLongParam       string = "server-url"
	getKubeconfigUserLongParam            string = "user"
)

var (
	getKubeconfigTargetNamespace string
	getKubeconfigServerUrl       string
	getKubeconfigUser            string
)

// getKubeconfigCmd represents the kubeconfig command
var getKubeconfigCmd = &cobra.Command{
	Use:   "kubeconfig <identity>",
	Short: "Display the kubeconfig for authenticating as the given identity",
	Long: `Creates and prints to stdout the kubeconfig for authenticating as the given identity.
The token embedded in the kubeconfig is the last one created.`,
	Args: cobra.MatchAll(cobra.ExactArgs(1)),
	RunE: func(cmd *cobra.Command, args []string) error {
		cli, err := kube.GetCurrentContextClient()
		if err != nil {
			return err
		}

		ctx := cmd.Context()

		name := args[0]
		s, err := kube.GetLastServiceAccountSecrets(ctx, *cli, name, namespace)
		if err != nil {
			return err
		}

		tkn, err := identity.GetToken(s)
		if err != nil {
			return err
		}

		o := getKubeconfigOptionsFromFlags(cmd.Flags())
		kfg, err := identity.GetKubeconfig(*cli, tkn, o)
		if err != nil {
			return err
		}

		fmt.Println(string(kfg))

		return nil
	},
}

func init() {
	getCmd.AddCommand(getKubeconfigCmd)

	getKubeconfigCmd.Flags().StringVarP(&getKubeconfigTargetNamespace, getKubeconfigTargetNamespaceLongParam, "t", "", "Target namespace to set in kubeconfig")
	getKubeconfigCmd.Flags().StringVarP(&getKubeconfigServerUrl, getKubeconfigServerUrlLongParam, "s", "", "if set overrides the cluster server URL")
	getKubeconfigCmd.Flags().StringVarP(&getKubeconfigUser, getKubeconfigUserLongParam, "u", "", "if set overrides the user")
}

func getKubeconfigOptionsFromFlags(ff *pflag.FlagSet) identity.GetKubeconfigOptions {
	o := identity.GetKubeconfigOptions{}

	if ff.Changed(getKubeconfigTargetNamespaceLongParam) {
		o.Namespace = &getKubeconfigTargetNamespace
	}

	if ff.Changed(getKubeconfigServerUrlLongParam) {
		o.OverrideHost = &getKubeconfigServerUrl
	}

	if ff.Changed(getKubeconfigUserLongParam) {
		o.User = &getKubeconfigUser
	}

	return o
}
