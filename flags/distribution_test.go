/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package flags_test

import (
	"testing"

	"github.com/pinterest/bender"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/facebookincubator/fbender/flags"
)

func TestNewDefaultDistribution(t *testing.T) {
	distribution := flags.NewDefaultDistribution()
	require.NotNil(t, distribution)
	assert.Equal(t, distribution.Name, "uniform")
	assert.Equal(t, distribution.String(), "uniform")
	assertPointerEqual(t, distribution.Get(), bender.UniformIntervalGenerator, "Expected uniform distribution")
}

func TestDistributionChoices(t *testing.T) {
	expected := []string{"uniform", "exponential"}
	assert.ElementsMatch(t, flags.DistributionChoices(), expected)
}

func TestDistribution__String(t *testing.T) {
	distribution := new(flags.Distribution)
	err := distribution.Set("uniform")
	require.NoError(t, err)
	assert.Equal(t, distribution.String(), "uniform")
	err = distribution.Set("exponential")
	require.NoError(t, err)
	assert.Equal(t, distribution.String(), "exponential")
	// Check if the name is full even if we do prefix match
	err = distribution.Set("uni")
	require.NoError(t, err)
	assert.Equal(t, distribution.String(), "uniform")
	err = distribution.Set("exp")
	require.NoError(t, err)
	assert.Equal(t, distribution.String(), "exponential")
}

func TestDistribution__Set(t *testing.T) {
	distribution := new(flags.Distribution)
	// Setting known distribution
	err := distribution.Set("uniform")
	require.NoError(t, err)
	assert.Equal(t, distribution.Name, "uniform")
	assertPointerEqual(t, distribution.Get(), bender.UniformIntervalGenerator, "Expected uniform distribution")
	err = distribution.Set("exponential")
	require.NoError(t, err)
	assert.Equal(t, distribution.Name, "exponential")
	assertPointerEqual(t, distribution.Get(), bender.ExponentialIntervalGenerator, "Expected exponential distribution")
	// Setting known distribution through an unambiguous prefix
	err = distribution.Set("u")
	require.NoError(t, err)
	assert.Equal(t, distribution.Name, "uniform")
	assertPointerEqual(t, distribution.Get(), bender.UniformIntervalGenerator, "Expected uniform distribution")
	err = distribution.Set("e")
	require.NoError(t, err)
	assert.Equal(t, distribution.Name, "exponential")
	assertPointerEqual(t, distribution.Get(), bender.ExponentialIntervalGenerator, "Expected exponential distribution")
	// Setting unknown distribution should fail
	err = distribution.Set("unknown")
	assert.EqualError(t, err, "generator must be one of (exponential|uniform), 'unknown' given")
	// Setting known distribution through an ambiguous prefix should fail
	err = distribution.Set("")
	// We cannot check directly because map is unordered
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ambiguous generator '' matches")
}

func TestDistribution__Type(t *testing.T) {
	distribution := new(flags.Distribution)
	assert.Equal(t, distribution.Type(), "distribution")
}

func TestDistribution__Get(t *testing.T) {
	distribution := new(flags.Distribution)
	err := distribution.Set("uniform")
	require.NoError(t, err)
	assertPointerEqual(t, distribution.Get(), bender.UniformIntervalGenerator, "Expected uniform distribution")
	err = distribution.Set("exponential")
	require.NoError(t, err)
	assertPointerEqual(t, distribution.Get(), bender.ExponentialIntervalGenerator, "Expected exponential distribution")
	// Check if the distribution function is properly set if we do prefix match
	err = distribution.Set("uni")
	require.NoError(t, err)
	assertPointerEqual(t, distribution.Get(), bender.UniformIntervalGenerator, "Expected uniform distribution")
	err = distribution.Set("exp")
	require.NoError(t, err)
	assertPointerEqual(t, distribution.Get(), bender.ExponentialIntervalGenerator, "Expected exponential distribution")
}

func TestGetDistribution(t *testing.T) {
	distribution := new(flags.Distribution)
	f := pflag.NewFlagSet("Test FlagSet", pflag.ExitOnError)
	f.Var(distribution, "distribution", "set distribution")
	err := distribution.Set("uniform")
	require.NoError(t, err)
	dist, err := flags.GetDistribution(f, "distribution")
	require.NoError(t, err)
	assertPointerEqual(t, dist, bender.UniformIntervalGenerator)
	// Check if value changes
	err = distribution.Set("exponential")
	require.NoError(t, err)
	dist, err = flags.GetDistribution(f, "distribution")
	require.NoError(t, err)
	assertPointerEqual(t, dist, bender.ExponentialIntervalGenerator)
	// Check if error when flag does not exist
	_, err = flags.GetDistribution(f, "nonexistent")
	assert.EqualError(t, err, "flag nonexistent accessed but not defined")
	// Check if error when value is of different type
	f.Int("myint", 0, "set myint")
	_, err = flags.GetDistribution(f, "myint")
	assert.EqualError(t, err, "trying to get distribution value of flag of type int")
}

func TestGetDistributionValue(t *testing.T) {
	distribution := new(flags.Distribution)
	err := distribution.Set("uniform")
	require.NoError(t, err)
	dist, err := flags.GetDistributionValue(distribution)
	require.NoError(t, err)
	assertPointerEqual(t, dist, bender.UniformIntervalGenerator)
	// Check if value changes
	err = distribution.Set("exponential")
	require.NoError(t, err)
	dist, err = flags.GetDistributionValue(distribution)
	require.NoError(t, err)
	assertPointerEqual(t, dist, bender.ExponentialIntervalGenerator)
	// Check if error when value is of different type
	f := pflag.NewFlagSet("Test FlagSet", pflag.ExitOnError)
	f.Int("myint", 0, "set myint")
	flag := f.Lookup("myint")
	require.NotNil(t, flag)
	_, err = flags.GetDistributionValue(flag.Value)
	assert.EqualError(t, err, "trying to get distribution value of flag of type int")
}

func TestBashCompletionDistribution(t *testing.T) {
	c := &cobra.Command{}
	d := flags.NewDefaultDistribution()
	// Check no error when applied to distribution flag
	f := c.Flags().VarPF(d, "distribution", "d", "set distribution")
	err := flags.BashCompletionDistribution(c, c.Flags(), "distribution")
	require.NoError(t, err)
	require.Contains(t, f.Annotations, "cobra_annotation_bash_completion_custom")
	assert.Equal(t, []string{"__fbender_handle_distribution_flag"},
		f.Annotations["cobra_annotation_bash_completion_custom"])
	// Check error when flag is not defined
	err = flags.BashCompletionDistribution(c, c.Flags(), "nonexistent")
	assert.EqualError(t, err, "flag nonexistent accessed but not defined")
	// Check error when flag is not a distribution
	c.Flags().Int("myint", 0, "set myint")
	err = flags.BashCompletionDistribution(c, c.Flags(), "myint")
	assert.EqualError(t, err, "trying to autocomplete distribution on flag of type int")
}
