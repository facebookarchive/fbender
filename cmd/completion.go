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
