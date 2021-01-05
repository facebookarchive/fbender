/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package tester_test

import (
	"math"
	"testing"

	"github.com/facebookincubator/fbender/tester"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLinearGrowth__OnSuccess(t *testing.T) {
	// test checks if the linear growth follows the arithmetic sequence starting
	// at s, growing by r for n steps
	test := func(s, r, n int) {
		g := &tester.LinearGrowth{Increase: r}
		a := s

		for i := 1; i < n; i++ {
			a = g.OnSuccess(a)
			assert.Equal(t, s+r*i, a)
		}
	}

	test(10, 3, 100)
	test(42, 42, 42)
	test(100, 1, 100)
}

func TestLinearGrowth__OnFail(t *testing.T) {
	test := func(s, r, n int) {
		g := &tester.LinearGrowth{Increase: r}
		a := s

		for i := 1; i < n; i++ {
			a = g.OnSuccess(a)
		}

		a = g.OnFail(a)
		assert.Equal(t, 0, a)
	}

	test(10, 3, 100)
	test(42, 42, 42)
	test(100, 1, 100)
}

func TestLinearGrowth__String(t *testing.T) {
	g := &tester.LinearGrowth{Increase: 3}
	assert.Equal(t, g.String(), "+3")

	g = &tester.LinearGrowth{Increase: 42}
	assert.Equal(t, g.String(), "+42")

	g = &tester.LinearGrowth{Increase: 100}
	assert.Equal(t, g.String(), "+100")
}

func TestParseGrowth_LinearGrowth(t *testing.T) {
	// Valid linear growth
	g, err := tester.ParseGrowth("+200")
	require.NoError(t, err)
	assert.IsType(t, new(tester.LinearGrowth), g)
	assert.Equal(t, 200, g.(*tester.LinearGrowth).Increase)
	// Invalid value
	_, err = tester.ParseGrowth("+abcdef")
	assert.EqualError(t, err, "strconv.Atoi: parsing \"abcdef\": invalid syntax")

	_, err = tester.ParseGrowth("+99.9")
	assert.EqualError(t, err, "strconv.Atoi: parsing \"99.9\": invalid syntax")
}

func TestPercentageGrowth__OnSuccess(t *testing.T) {
	// test checks if the linear growth follows the arithmetic sequence starting
	// at s, growing by r for n steps
	test := func(s int, r float64, n int) {
		g := &tester.PercentageGrowth{Increase: r}
		a := s

		for i := 1; i < n; i++ {
			a = g.OnSuccess(a)
			expected := int(float64(s) * math.Pow((100+r)/100., float64(i)))
			// We're reounding every result to int so it may not be equal
			assert.InDelta(t, expected, a, float64(s)*r/100.)
		}
	}

	test(10, 100, 100)
	test(2, 200, 10)
	test(100, 100, 10)
}

func TestPercentageGrowth__OnFail(t *testing.T) {
	test := func(s int, r float64, n int) {
		g := &tester.PercentageGrowth{Increase: r}
		a := s

		for i := 1; i < n; i++ {
			a = g.OnSuccess(a)
		}

		a = g.OnFail(a)
		assert.Equal(t, 0, a)
	}

	test(10, 3, 100)
	test(42, 42, 42)
	test(100, 1, 100)
}

func TestPercentageGrowth__String(t *testing.T) {
	g := &tester.PercentageGrowth{Increase: 102.5}
	assert.Equal(t, g.String(), "%102.50")

	g = &tester.PercentageGrowth{Increase: 66.66}
	assert.Equal(t, g.String(), "%66.66")

	g = &tester.PercentageGrowth{Increase: 100}
	assert.Equal(t, g.String(), "%100.00")
}

func TestParseGrowth_PercentageGrowth(t *testing.T) {
	// Valid linear growth
	g, err := tester.ParseGrowth("%100.50")
	require.NoError(t, err)

	assert.IsType(t, new(tester.PercentageGrowth), g)
	assert.Equal(t, 100.50, g.(*tester.PercentageGrowth).Increase)

	// Invalid value
	_, err = tester.ParseGrowth("%abcdef")
	assert.EqualError(t, err, "strconv.ParseFloat: parsing \"abcdef\": invalid syntax")
}

func TestExponentialGrowth__OnSuccess(t *testing.T) {
	g := &tester.ExponentialGrowth{Precision: 10}

	// When not bound it should double the test value
	assert.Equal(t, 40, g.OnSuccess(20))
	assert.Equal(t, 80, g.OnSuccess(40))
	assert.Equal(t, 160, g.OnSuccess(80))
	// When bound it should return a middle value
	assert.Equal(t, 120, g.OnFail(160))
	assert.Equal(t, 140, g.OnSuccess(120))
	assert.Equal(t, 150, g.OnSuccess(140))
	// Finally return 0 when precision is met
	assert.Equal(t, 0, g.OnSuccess(150))
}

func TestExponentialGrowth__OnFail(t *testing.T) {
	g := &tester.ExponentialGrowth{Precision: 10}

	assert.Equal(t, 40, g.OnSuccess(20))
	assert.Equal(t, 80, g.OnSuccess(40))
	assert.Equal(t, 160, g.OnSuccess(80))
	// It should set the upper bound and return middle value
	assert.Equal(t, 120, g.OnFail(160))
	assert.Equal(t, 100, g.OnFail(120))
	assert.Equal(t, 90, g.OnFail(100))
	// Finally return 0 when precision is met
	assert.Equal(t, 0, g.OnSuccess(90))
}

func TestExponentialGrowth__String(t *testing.T) {
	g := &tester.ExponentialGrowth{Precision: 3}
	assert.Equal(t, g.String(), "^3")

	g = &tester.ExponentialGrowth{Precision: 42}
	assert.Equal(t, g.String(), "^42")

	g = &tester.ExponentialGrowth{Precision: 100}
	assert.Equal(t, g.String(), "^100")
}

func TestParseGrowth_ExponentialGrowth(t *testing.T) {
	// Valid linear growth
	g, err := tester.ParseGrowth("^10")
	require.NoError(t, err)

	assert.IsType(t, new(tester.ExponentialGrowth), g)
	assert.Equal(t, 10, g.(*tester.ExponentialGrowth).Precision)

	// Invalid value
	_, err = tester.ParseGrowth("^abcdef")
	assert.EqualError(t, err, "strconv.Atoi: parsing \"abcdef\": invalid syntax")

	_, err = tester.ParseGrowth("^99.9")
	assert.EqualError(t, err, "strconv.Atoi: parsing \"99.9\": invalid syntax")
}

func TestParseGrowth(t *testing.T) {
	g, err := tester.ParseGrowth("@200")
	assert.Nil(t, g)
	assert.Equal(t, tester.ErrInvalidGrowth, err)
}
