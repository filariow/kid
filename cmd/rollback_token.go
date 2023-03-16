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
	"strconv"

	"github.com/filariow/kid/pkg/identity"
	"github.com/filariow/kid/pkg/kube"
	"github.com/spf13/cobra"
)

// rollbackTokenCmd represents the token command
var rollbackTokenCmd = &cobra.Command{
	Use:   "token <identity> <version>",
	Short: "Rollback a token",
	Long:  `Rollback the token with a given version for the given identity.`,
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		cli, err := kube.GetCurrentContextClient()
		if err != nil {
			return err
		}

		uv, err := strconv.ParseUint(args[1], 10, 64)
		if err != nil {
			return err
		}

		s, err := identity.RollbackIdentityKey(cmd.Context(), *cli, args[0], namespace, uv)
		if err != nil {
			return err
		}

		fmt.Printf("secret rolled back '%s/%s'\n", s.Namespace, s.Name)
		return nil
	},
}

func init() {
	rollbackCmd.AddCommand(rollbackTokenCmd)
}
