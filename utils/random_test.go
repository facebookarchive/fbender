/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package utils_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/facebookincubator/fbender/utils"
)

func TestRandomHex(t *testing.T) {
	// Even length
	hex, err := utils.RandomHex(16)
	assert.NoError(t, err)
	assert.Len(t, hex, 16)
	// Odd length
	hex, err = utils.RandomHex(9)
	assert.NoError(t, err)
	assert.Len(t, hex, 9)
	// Corner case
	hex, err = utils.RandomHex(0)
	assert.NoError(t, err)
	assert.Len(t, hex, 0)
}
