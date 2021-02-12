/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package tester

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// Growth is used to determine what test should be ran next.
type Growth interface {
	OnSuccess(test int) int
	OnFail(test int) int
	String() string
}

// LinearGrowth increases test by a specified amount with every successful test.
type LinearGrowth struct {
	Increase int
}

// LinearGrowthPrefix prefix used in linear growth string representation.
const LinearGrowthPrefix = "+"

func (g *LinearGrowth) String() string {
	return fmt.Sprintf("%s%d", LinearGrowthPrefix, g.Increase)
}

// OnSuccess increases test by a specified amount.
func (g *LinearGrowth) OnSuccess(test int) int {
	return test + g.Increase
}

// OnFail stops the tests.
func (g *LinearGrowth) OnFail(test int) int {
	return 0
}

// PercentageGrowth increases test by a specified percentage with every successful test.
type PercentageGrowth struct {
	Increase float64
}

// PercentageGrowthPrefix prefix used in percentage growth string representation.
const PercentageGrowthPrefix = "%"

func (g *PercentageGrowth) String() string {
	return fmt.Sprintf("%s%.2f", PercentageGrowthPrefix, g.Increase)
}

// OnSuccess increases test by a specified percentage.
func (g *PercentageGrowth) OnSuccess(test int) int {
	return int((100. + g.Increase) / 100. * float64(test))
}

// OnFail stops the tests.
func (g *PercentageGrowth) OnFail(test int) int {
	return 0
}

// ExponentialGrowth performs binary search up to a given precision.
type ExponentialGrowth struct {
	Precision int

	left, right int
	bound       bool
}

// ExponentialGrowthPrefix prefix used in exponential growth string representation.
const ExponentialGrowthPrefix = "^"

func (g *ExponentialGrowth) String() string {
	return fmt.Sprintf("%s%d", ExponentialGrowthPrefix, g.Precision)
}

// OnSuccess sets the lower bound to the last test and returns (left+right) / 2
// unless the precision has been achieved.
func (g *ExponentialGrowth) OnSuccess(test int) int {
	g.left = test
	if !g.bound {
		return test * 2
	}

	if g.right-g.left <= g.Precision {
		return 0
	}

	return int(float64(g.right+g.left) / 2)
}

// OnFail sets the upper bound to the last test and returns (left+right) / 2
// unless the precision has been achieved.
func (g *ExponentialGrowth) OnFail(test int) int {
	g.right = test
	g.bound = true

	if g.right-g.left <= g.Precision {
		return 0
	}

	return int(float64(g.right+g.left) / 2)
}

// GrowthHelp provides usage help about the growth.
const GrowthHelp = `Growth determines what will be the next value used for a test.
* linear growth (+int) increases test value by a fixed amount after each success,
  stops immediately after the first failure
* percentage growth (%float) increases test value by a fixed percentage after
  each success, stops immediately after the first failure
* exponential growth (^int) first doubles the test value after each success to
  find an upper bound, then performs a binary search up to a given precision`

// ErrInvalidGrowth is returned when a growth cannot be found.
var ErrInvalidGrowth = errors.New("unknown growth, want +int, %%flaot, ^int")

// ParseGrowth creates a growth from its string representation.
func ParseGrowth(value string) (Growth, error) {
	switch {
	case strings.HasPrefix(value, LinearGrowthPrefix):
		inc, err := strconv.Atoi(strings.TrimPrefix(value, LinearGrowthPrefix))
		if err != nil {
			//nolint:wrapcheck
			return nil, err
		}

		return &LinearGrowth{Increase: inc}, nil

	case strings.HasPrefix(value, PercentageGrowthPrefix):
		inc, err := strconv.ParseFloat(strings.TrimPrefix(value, PercentageGrowthPrefix), 64)
		if err != nil {
			//nolint:wrapcheck
			return nil, err
		}

		return &PercentageGrowth{Increase: inc}, nil

	case strings.HasPrefix(value, ExponentialGrowthPrefix):
		prec, err := strconv.Atoi(strings.TrimPrefix(value, ExponentialGrowthPrefix))
		if err != nil {
			//nolint:wrapcheck
			return nil, err
		}

		return &ExponentialGrowth{Precision: prec}, nil

	default:
		return nil, ErrInvalidGrowth
	}
}
