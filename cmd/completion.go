/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var completionCmd = &cobra.Command{
	Use:   "completion",
	Short: "Generates completion scripts",
}

var bashCompletionCmd = &cobra.Command{
	Use:   "bash",
	Short: "Generates bash completion scripts",
	Long: `To load completion run

. <(fbender completion)

To configure your bash shell to load completions for each session add to .bashrc

# ~/.bashrc or ~/.profile
. <(fbender completion)
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return Command.GenBashCompletion(os.Stdout)
	},
}

func init() {
	completionCmd.AddCommand(bashCompletionCmd)
}
