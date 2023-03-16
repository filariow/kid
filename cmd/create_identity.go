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
	"encoding/json"
	"fmt"

	"github.com/filariow/kid/pkg/identity"
	"github.com/filariow/kid/pkg/kube"
	"github.com/spf13/cobra"
)

// createIdentityCmd represents the identity command
var createIdentityCmd = &cobra.Command{
	Use:   "identity <name>",
	Short: "Create a new identity",
	Long: `A new service account and a first secret are created.

The service account is named after the identity, whereas the secret name
respect the following format '<identity>-key-<number>'.`,
	Args: cobra.MatchAll(cobra.ExactArgs(1)),
	RunE: func(cmd *cobra.Command, args []string) error {
		cli, err := kube.GetCurrentContextClient()
		if err != nil {
			return err
		}

		ctx := cmd.Context()
		name := args[0]
		i, err := identity.CreateIdentity(ctx, *cli, name, namespace)
		if err != nil {
			return err
		}

		j, err := json.MarshalIndent(i, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(j))

		return nil
	},
}

func init() {
	createCmd.AddCommand(createIdentityCmd)
}
