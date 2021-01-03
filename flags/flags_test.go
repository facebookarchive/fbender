/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package flags_test

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// assertPointerEqual checks whether two pointers are equal.
func assertPointerEqual(t *testing.T, value, expected interface{}, args ...interface{}) {
	t.Helper()

	valuePointer := reflect.ValueOf(value).Pointer()
	expectedPointer := reflect.ValueOf(expected).Pointer()
	assert.Equal(t, valuePointer, expectedPointer, args...)
}
