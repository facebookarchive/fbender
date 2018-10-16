package utils

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const template = `
%s() {
	%s
}`

// BashCompletion annotates the flag with completion function and registers the
// completion function in the root command if it hasn't been added already.
func BashCompletion(cmd *cobra.Command, flags *pflag.FlagSet, flag string, fname, fbody string) error {
	if !strings.Contains(cmd.Root().BashCompletionFunction, fname) {
		cmd.Root().BashCompletionFunction += fmt.Sprintf(template, fname, fbody)
	}
	return cobra.MarkFlagCustom(flags, flag, fname)
}
