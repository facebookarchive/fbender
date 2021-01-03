/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package utils_test

import (
	"testing"

	"github.com/facebookincubator/fbender/utils"
	"github.com/stretchr/testify/assert"
)

func TestWithDefaultPort(t *testing.T) {
	// Adds default port
	assert.Equal(t, "[::1]:53", utils.WithDefaultPort("::1", 53))
	assert.Equal(t, "127.0.0.1:53", utils.WithDefaultPort("127.0.0.1", 53))
	assert.Equal(t, "example.com:53", utils.WithDefaultPort("example.com", 53))
	// Does not change port when present
	assert.Equal(t, "[::1]:5353", utils.WithDefaultPort("[::1]:5353", 53))
	assert.Equal(t, "127.0.0.1:5353", utils.WithDefaultPort("127.0.0.1:5353", 53))
	assert.Equal(t, "example.com:5353", utils.WithDefaultPort("example.com:5353", 53))
}
