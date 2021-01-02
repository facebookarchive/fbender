/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package core

import (
	"strconv"

	"github.com/facebookincubator/fbender/cmd/core/options"
	"github.com/facebookincubator/fbender/flags"
	"github.com/spf13/cobra"
)

// OptionsGenerator is used to generate options from command line params.
type OptionsGenerator func(o *options.Options, cmd *cobra.Command, args []string) (*options.Options, error)

// ExtractTests parses list of arguments to a list of int values.
func ExtractTests(args []string) ([]int, error) {
	values := make([]int, 0)

	for _, arg := range args {
		value, err := strconv.Atoi(arg)
		if err != nil {
			return nil, err
		}

		values = append(values, value)
	}

	return values, nil
}

// ExtractArgs extracts arguments commonly used options across all tests.
func ExtractArgs(o *options.Options, cmd *cobra.Command, args []string) (*options.Options, error) {
	var err error

	if o == nil {
		o = options.NewOptions()
	}

	o.Tests, err = ExtractTests(args)
	if err != nil {
		return nil, err
	}

	o.Start = o.Tests[0]

	return o, nil
}

// ExtractOptions extracts flags commonly used options across all commands.
func ExtractOptions(o *options.Options, cmd *cobra.Command, _ []string) (*options.Options, error) {
	var err error

	if o == nil {
		o = options.NewOptions()
	}

	o.Target, err = cmd.Flags().GetString("target")
	if err != nil {
		return nil, err
	}

	o.Duration, err = cmd.Flags().GetDuration("duration")
	if err != nil {
		return nil, err
	}

	o.Input, err = cmd.Flags().GetString("input")
	if err != nil {
		return nil, err
	}

	o.BufferSize, err = cmd.Flags().GetInt("buffer")
	if err != nil {
		return nil, err
	}

	o.Timeout, err = cmd.Flags().GetDuration("timeout")
	if err != nil {
		return nil, err
	}

	o.Distribution, err = flags.GetDistribution(cmd.Flags(), "dist")
	if err != nil {
		return nil, err
	}

	o.Unit, err = cmd.Flags().GetDuration("unit")
	if err != nil {
		return nil, err
	}

	o.NoStatistics, err = cmd.Flags().GetBool("nostats")
	if err != nil {
		return nil, err
	}

	return o, nil
}

// ExtractConstraintsOptions extracts flag commonly used options across constraints test commands.
func ExtractConstraintsOptions(o *options.Options, cmd *cobra.Command, _ []string) (*options.Options, error) {
	var err error

	if o == nil {
		o = options.NewOptions()
	}

	o.Constraints, err = flags.GetConstraints(cmd.Flags(), "constraints")
	if err != nil {
		return nil, err
	}

	o.Growth, err = flags.GetGrowth(cmd.Flags(), "growth")
	if err != nil {
		return nil, err
	}

	return o, nil
}

// GenerateOptions runs given generators for a command and returns options.
func GenerateOptions(cmd *cobra.Command, args []string, gs ...OptionsGenerator) (*options.Options, error) {
	var o *options.Options

	var err error

	for _, g := range gs {
		o, err = g(o, cmd, args)
		if err != nil {
			return nil, err
		}
	}

	return o, nil
}
