/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package tester_test

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// assertPointerEqual checks whether two pointers are equal
func assertPointerEqual(t *testing.T, expected, value interface{}, args ...interface{}) {
	expectedPointer := reflect.ValueOf(expected).Pointer()
	valuePointer := reflect.ValueOf(value).Pointer()
	assert.Equal(t, expectedPointer, valuePointer, args...)
}
