/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package tester_test

import (
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/facebookincubator/fbender/tester"
)

func TestLessThan__Compare(t *testing.T) {
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)
	for i := 0; i < 4096; i++ {
		x, y := r.Float64(), r.Float64()
		assert.Equal(t, x < y, tester.LessThan.Compare(x, y), "%f < %f", x, y)
	}
}

func TestLessThan__Name(t *testing.T) {
	assert.Equal(t, "<", tester.LessThan.Name())
}

func TestParseComparator_LessThan(t *testing.T) {
	cmp, err := tester.ParseComparator("<")
	assert.NoError(t, err)
	assertPointerEqual(t, tester.LessThan, cmp)
}

func TestGreaterThan__Compare(t *testing.T) {
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)
	for i := 0; i < 4096; i++ {
		x, y := r.Float64(), r.Float64()
		assert.Equal(t, x > y, tester.GreaterThan.Compare(x, y), "%f > %f", x, y)
	}
}

func TestGreaterThan__Name(t *testing.T) {
	assert.Equal(t, ">", tester.GreaterThan.Name())
}

func TestParseComparator_GreaterThan(t *testing.T) {
	cmp, err := tester.ParseComparator(">")
	assert.NoError(t, err)
	assertPointerEqual(t, tester.GreaterThan, cmp)
}

func TestParseComparator(t *testing.T) {
	cmp, err := tester.ParseComparator("!")
	assert.Nil(t, cmp)
	assert.Equal(t, tester.ErrInvalidComparator, err)
}
