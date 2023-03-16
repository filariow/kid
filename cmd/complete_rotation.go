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

// completeRotationCmd represents the rotation command
var completeRotationCmd = &cobra.Command{
	Use:   "rotation <identity>",
	Short: "Complete the key rotation for an identity",
	Long: `Key rotation is performed in two phases.
You create a new key and update your services with this new one.
Finally, you remove the old one.

This step deletes the old key`,
	Args: cobra.MatchAll(cobra.ExactArgs(1)),
	RunE: func(cmd *cobra.Command, args []string) error {
		cli, err := kid.GetCurrentContextClient()
		if err != nil {
			return err
		}

		ctx := cmd.Context()
		name := args[0]
		s, err := identity.CompleteIdentityKeyRotation(ctx, *cli, name, namespace)
		if err != nil {
			return err
		}

		fmt.Printf("deleted secret '%s/%s'\n", namespace, *s)
		return nil
	},
}

func init() {
	completeCmd.AddCommand(completeRotationCmd)
}
