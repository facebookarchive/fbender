/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package core

import (
	"fmt"
	"strings"

	"github.com/facebookincubator/fbender/cmd/core/errors"
	"github.com/facebookincubator/fbender/cmd/core/options"
	"github.com/facebookincubator/fbender/cmd/core/runner"
	"github.com/spf13/cobra"
)

// CommandParams is used to generate params for the runner.
type CommandParams func(*cobra.Command, *options.Options) (*runner.Params, error)

// CommandTemplate groups parameters used to generate the test command.
type CommandTemplate struct {
	// Command name it will be invoked with
	Name string
	// Short and long help messages
	Short, Long string
	// Usage examples can containt {test} string which will be replaced with the
	// actual test name (throughput|concurrency)
	Fixed, Constraints string
}

// NewTestCommand creates a protocol command with all the test subcommands.
// ${protocol} throughput fixed
// ${protocol} throughput constraints
// ${protocol} concurrency fixed
// ${protocol} concurrency constraints
//nolint:funlen
func NewTestCommand(c *CommandTemplate, p CommandParams) *cobra.Command {
	// Help messages
	var (
		tShort  = fmt.Sprintf("%s throughput (QPS)", c.Short)
		tfShort = fmt.Sprintf("%s with fixed amount of QPS", tShort)
		tcShort = fmt.Sprintf("%s with constraints", tShort)

		cShort  = fmt.Sprintf("%s concurrent connections", c.Short)
		cfShort = fmt.Sprintf("%s with fixed number of connections", cShort)
		ccShort = fmt.Sprintf("%s with constraints", cShort)
	)

	// Examples
	var (
		tfExamples = strings.ReplaceAll(c.Fixed, "{test}", "throughput")
		tcExamples = strings.ReplaceAll(c.Constraints, "{test}", "throughput")
		tExamples  = fmt.Sprintf("%s\n%s", tfExamples, tcExamples)

		cfExamples = strings.ReplaceAll(c.Fixed, "{test}", "concurrency")
		ccExamples = strings.ReplaceAll(c.Constraints, "{test}", "concurrency")
		cExamples  = fmt.Sprintf("%s\n%s", cfExamples, ccExamples)

		examples = fmt.Sprintf("%s\n%s", tExamples, cExamples)
	)

	// Top level command
	command := &cobra.Command{
		Use:     c.Name,
		Short:   c.Short,
		Long:    fmt.Sprintf("%s.\n%s", c.Short, c.Long),
		Example: examples,
	}

	// Subcommands
	tCommand := &cobra.Command{
		Use:     "throughput",
		Short:   tShort,
		Long:    fmt.Sprintf("%s.\n%s", tShort, c.Long),
		Example: tExamples,
	}

	cCommand := &cobra.Command{
		Use:     "concurrency",
		Short:   cShort,
		Long:    fmt.Sprintf("%s.\n%s", cShort, c.Long),
		Example: cExamples,
	}

	command.AddCommand(tCommand)
	command.AddCommand(cCommand)

	// Throughput subcommands
	tfCommand := &cobra.Command{
		Use:     "fixed",
		Short:   tfShort,
		Long:    fmt.Sprintf("%s.\n%s", tfShort, c.Long),
		Example: tfExamples,
		Args:    fixedArgs,
		RunE:    RunLoadTestThroughputFixed(p),
	}

	tcCommand := &cobra.Command{
		Use:     "constraints",
		Short:   tcShort,
		Long:    fmt.Sprintf("%s.\n%s\n%s", tcShort, c.Long, ConstraintsHelp),
		Example: tcExamples,
		Args:    constraintsArgs,
		RunE:    RunLoadTestThroughputConstraints(p),
	}

	tcCommand.PersistentFlags().AddFlagSet(ConstraintsFlags)
	tCommand.AddCommand(tfCommand)
	tCommand.AddCommand(tcCommand)

	// Concurrency subcommands
	cfCommand := &cobra.Command{
		Use:     "fixed",
		Short:   cfShort,
		Long:    fmt.Sprintf("%s.\n%s", cfShort, c.Long),
		Example: cfExamples,
		Args:    fixedArgs,
		RunE:    RunLoadTestConcurrencyFixed(p),
	}

	ccCommand := &cobra.Command{
		Use:     "constraints",
		Short:   ccShort,
		Long:    fmt.Sprintf("%s.\n%s\n%s", ccShort, c.Long, ConstraintsHelp),
		Example: ccExamples,
		Args:    constraintsArgs,
		RunE:    RunLoadTestConcurrencyConstraints(p),
	}

	ccCommand.PersistentFlags().AddFlagSet(ConstraintsFlags)
	cCommand.AddCommand(cfCommand)
	cCommand.AddCommand(ccCommand)

	return command
}

// fixedArgs validates arguments for a fixed QPS test.
func fixedArgs(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("%w: requires at least one test value", errors.ErrInvalidArgument)
	}

	_, err := ExtractTests(args)

	return err
}

// constraintsArgs validates arguments for a constraints tests.
func constraintsArgs(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("%w: requires starting test value", errors.ErrInvalidArgument)
	}

	_, err := ExtractTests(args)

	return err
}
