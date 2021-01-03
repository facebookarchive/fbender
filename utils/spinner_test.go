/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package utils_test

import (
	"bytes"
	"testing"
	"time"

	"github.com/facebookincubator/fbender/log"
	"github.com/facebookincubator/fbender/utils"
	"github.com/stretchr/testify/assert"
)

func TestNewBackgroundSpinner(t *testing.T) {
	w := new(bytes.Buffer)
	log.Stderr = w

	cancel := utils.NewBackgroundSpinner("Testing", 100*time.Millisecond)

	time.Sleep(1 * time.Second)
	cancel()

	v := w.String()

	assert.Contains(t, v, "\rTesting... |")
	assert.Contains(t, v, "\rTesting... /")
	assert.Contains(t, v, "\rTesting... -")
	assert.Contains(t, v, "\rTesting... \\")
	assert.Contains(t, v, "\rTesting... Done.")
}
