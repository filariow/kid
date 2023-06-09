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
	"os"

	"github.com/filariow/kid/pkg/kube"
	"github.com/spf13/cobra"
)

var namespace string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "kid",
	Short: "KId - Kubernetes Identity",
	Long: `KId (Kubernetes Identity) is a Command Line Application to create and manage Identities.
Each identity is a Service Account.
This tool helps managing JWT tokens (create, revoke, rotate, rollback) and exporting kubeconfig.

It uses heavily the convention over configuration paradigm and its not meant to
provide a solid and constraining workflow.It gives you a lot of freedom, so be wise.
Do not create tokens with the same version of revoked/leaked ones!`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	ns, err := kube.GetConfigDefaultNamespace()
	if err != nil {
		fmt.Fprint(os.Stderr, "can not parse namespace from kubeconfig, using default")
		*ns = "default"
	}
	rootCmd.PersistentFlags().StringVarP(&namespace, "namespace", "n", *ns, "the namespace where to operate")
}
