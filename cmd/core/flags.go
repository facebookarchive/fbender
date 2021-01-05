/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package core

import (
	"strings"

	"github.com/facebookincubator/fbender/flags"
	"github.com/facebookincubator/fbender/metric"
	"github.com/facebookincubator/fbender/tester"
	"github.com/spf13/pflag"
)

//nolint:gochecknoglobals
var (
	// ConstraintsFlags contains flags for specifying constraints tests options.
	ConstraintsFlags = pflag.NewFlagSet("Constraints test flags", pflag.ExitOnError)
	// ConstraintsValue is a pflag value for constraints.
	ConstraintsValue = flags.NewConstraintSliceValue(metric.Parser)
	// ConstraintsHelp is a help message on how to use constraints.
	ConstraintsHelp = strings.Join([]string{tester.ConstraintsHelp, metric.Help}, "\n")
)

//nolint:gochecknoinits
func init() {
	growth := &flags.GrowthValue{Growth: &tester.PercentageGrowth{Increase: 100.}}

	ConstraintsFlags.VarP(ConstraintsValue, "constraints", "c", "constraints to be checked after each test")
	ConstraintsFlags.VarP(growth, "growth", "g", "growth used to determinate the next test (+AMOUNT|%PERCENT|^PRECISION)")
}
