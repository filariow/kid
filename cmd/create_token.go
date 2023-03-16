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
	"github.com/filariow/kid/pkg/kid"
	"github.com/spf13/cobra"
)

// createTokenCmd represents the token command
var createTokenCmd = &cobra.Command{
	Use:   "token <identity>",
	Short: "Creates a new token for the given identity",
	Long: `A new secret is created for the given identity.
The name of the secret is built starting from the ones present on the cluster.
Identity secrets respect the format <identity>-key-<number>.
The new secret will have the name <identity>-key<number+1>.`,
	Args: cobra.MatchAll(cobra.ExactArgs(1)),
	RunE: func(cmd *cobra.Command, args []string) error {
		cli, err := kid.GetCurrentContextClient()
		if err != nil {
			return err
		}

		ctx := cmd.Context()
		name := args[0]
		s, err := identity.CreateNewTokenVersion(ctx, *cli, name, namespace)
		if err != nil {
			return err
		}

		fmt.Printf("created secret '%s/%s'\n", s.Namespace, s.Name)
		return nil
	},
}

func init() {
	createCmd.AddCommand(createTokenCmd)
}
