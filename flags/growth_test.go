/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package flags_test

import (
	"testing"

	"github.com/facebookincubator/fbender/flags"
	"github.com/facebookincubator/fbender/tester"
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGrowth__String(t *testing.T) {
	linGrowth := &tester.LinearGrowth{Increase: 100}
	value := &flags.GrowthValue{Growth: linGrowth}
	assert.Equal(t, linGrowth.String(), value.String())
	// Change growth values
	linGrowth.Increase = 200
	assert.Equal(t, linGrowth.String(), value.String())
	// Change growth type
	perGrowth := &tester.PercentageGrowth{Increase: 100.}
	value.Growth = perGrowth
	assert.Equal(t, perGrowth.String(), value.String())
}

func TestGrowth__Set(t *testing.T) {
	value := new(flags.GrowthValue)
	// Valid LinearGrowth string
	err := value.Set("+200")
	require.NoError(t, err)
	assert.IsType(t, new(tester.LinearGrowth), value.Growth)
	// Valid PercentageGrowth string
	err = value.Set("%100.0")
	require.NoError(t, err)
	assert.IsType(t, new(tester.PercentageGrowth), value.Growth)
	// Valid ExponentialGrowth false
	err = value.Set("^25")
	require.NoError(t, err)
	assert.IsType(t, new(tester.ExponentialGrowth), value.Growth)
	// Unknown prefix
	err = value.Set("@200")
	assert.ErrorIs(t, err, tester.ErrInvalidGrowth)
	assert.EqualError(t, err, "error parsing growth \"@200\": unknown growth, want +int, %%flaot, ^int")
}

func TestGrowth__Type(t *testing.T) {
	value := &flags.GrowthValue{Growth: &tester.LinearGrowth{Increase: 100}}
	assert.Equal(t, "growth", value.Type())
}

func TestGetGrowth(t *testing.T) {
	value := &flags.GrowthValue{Growth: &tester.LinearGrowth{Increase: 100}}
	f := pflag.NewFlagSet("Test FlagSet", pflag.ExitOnError)
	f.Var(value, "growth", "set growth")
	err := value.Set("+200")
	require.NoError(t, err)
	growth, err := flags.GetGrowth(f, "growth")
	require.NoError(t, err)
	assert.IsType(t, new(tester.LinearGrowth), growth)
	// Check if value changes
	err = value.Set("%100")
	require.NoError(t, err)
	growth, err = flags.GetGrowth(f, "growth")
	require.NoError(t, err)
	assert.IsType(t, new(tester.PercentageGrowth), growth)
	// Check if error when flag does not exist
	_, err = flags.GetGrowth(f, "nonexistent")
	assert.ErrorIs(t, err, flags.ErrUndefined)
	assert.EqualError(t, err, "flag accessed but not defined: \"nonexistent\"")
	// Check if error when value is of different type
	f.Int("myint", 0, "set myint")
	_, err = flags.GetGrowth(f, "myint")
	assert.ErrorIs(t, err, flags.ErrInvalidType)
	assert.EqualError(t, err, "accessed flag type does not match, want: *GrowthValue, got: *pflag.intValue")
}

func TestGetGrowthValue(t *testing.T) {
	value := new(flags.GrowthValue)
	err := value.Set("+200")
	require.NoError(t, err)
	growth, err := flags.GetGrowthValue(value)
	require.NoError(t, err)
	assert.IsType(t, new(tester.LinearGrowth), growth)
	// Check if value changes
	err = value.Set("%100")
	require.NoError(t, err)
	growth, err = flags.GetGrowthValue(value)
	require.NoError(t, err)
	assert.IsType(t, new(tester.PercentageGrowth), growth)
	// Check if error when value is of different type
	f := pflag.NewFlagSet("Test FlagSet", pflag.ExitOnError)
	f.Int("myint", 0, "set myint")
	flag := f.Lookup("myint")
	require.NotNil(t, flag)
	_, err = flags.GetGrowthValue(flag.Value)
	assert.ErrorIs(t, err, flags.ErrInvalidType)
	assert.EqualError(t, err, "accessed flag type does not match, want: *GrowthValue, got: *pflag.intValue")
}
