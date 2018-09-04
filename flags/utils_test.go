/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package flags_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/facebookincubator/fbender/flags"
)

func TestChoicesString(t *testing.T) {
	// Empty list
	assert.Equal(t, "()", flags.ChoicesString([]string{}))
	// Single element list
	assert.Equal(t, "(one)", flags.ChoicesString([]string{"one"}))
	// Multiple elements list
	assert.Equal(t, "(one|two)", flags.ChoicesString([]string{"one", "two"}))
	assert.Equal(t, "(one|two|three)", flags.ChoicesString([]string{"one", "two", "three"}))
}
