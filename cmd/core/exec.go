/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package core

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/facebookincubator/fbender/cmd/core/options"
	"github.com/facebookincubator/fbender/cmd/core/runner"
	"github.com/facebookincubator/fbender/log"
	"github.com/facebookincubator/fbender/tester/run"
)

// cobraRunE is just an alias for function type which implements cobra.RunE
type cobraRunE = func(cmd *cobra.Command, args []string) error

// executor invokes actual test function with proper params
type executor func(p *runner.Params, o *options.Options) error

func exec(p CommandParams, e executor, gs ...OptionsGenerator) cobraRunE {
	return func(cmd *cobra.Command, args []string) error {
		o, err := GenerateOptions(cmd, args, gs...)
		if err != nil {
			return err
		}
		params, err := p(cmd, o)
		if err != nil {
			return err
		}
		// We want runtime errors to be logged and not trigger help message
		if err := e(params, o); err != nil {
			log.Errorf("Error: %v\n", err)
			os.Exit(1)
		}
		return nil
	}
}

func setupConstraints(o *options.Options, cmd *cobra.Command, args []string) (*options.Options, error) {
	for _, constraint := range o.Constraints {
		if err := constraint.Metric.Setup(o); err != nil {
			return nil, err
		}
	}
	return o, nil
}

var fixedOptionsGenerators = []OptionsGenerator{
	ExtractArgs,
	ExtractOptions,
}

var constraintsOptionsGenerators = []OptionsGenerator{
	ExtractArgs,
	ExtractOptions,
	ExtractConstraintsOptions,
	setupConstraints,
}

// RunLoadTestThroughputFixed returns a new cobra RunE method for the load
// tester with fixed QPS tests
func RunLoadTestThroughputFixed(p CommandParams) cobraRunE {
	return exec(p, fixedThroughputExecutor, fixedOptionsGenerators...)
}

func fixedThroughputExecutor(p *runner.Params, o *options.Options) error {
	return run.LoadTestThroughputFixed(runner.NewThroughputRunner(p), o, o.Tests...)
}

// RunLoadTestThroughputConstraints returns a new cobra RunE method for the QPS
// load tester with constraint checks
func RunLoadTestThroughputConstraints(p CommandParams) cobraRunE {
	return exec(p, constraintsThroughputExecutor, constraintsOptionsGenerators...)
}

func constraintsThroughputExecutor(p *runner.Params, o *options.Options) error {
	return run.LoadTestThroughputConstraints(runner.NewThroughputRunner(p), o, o.Start, o.Growth, o.Constraints...)
}

// RunLoadTestConcurrencyFixed returns a new cobra RunE method for the load
// tester with fixed concurrent connections count
func RunLoadTestConcurrencyFixed(p CommandParams) cobraRunE {
	return exec(p, fixedConcurrencyExecutor, fixedOptionsGenerators...)
}

func fixedConcurrencyExecutor(p *runner.Params, o *options.Options) error {
	return run.LoadTestConcurrencyFixed(runner.NewConcurrencyRunner(p), o, o.Tests...)
}

// RunLoadTestConcurrencyConstraints returns a new cobra RunE method for the
// concurrency load tester with constraint checks
func RunLoadTestConcurrencyConstraints(p CommandParams) cobraRunE {
	return exec(p, constraintsConcurrencyExecutor, constraintsOptionsGenerators...)
}

func constraintsConcurrencyExecutor(p *runner.Params, o *options.Options) error {
	return run.LoadTestConcurrencyConstraints(runner.NewConcurrencyRunner(p), o, o.Start, o.Growth, o.Constraints...)
}
