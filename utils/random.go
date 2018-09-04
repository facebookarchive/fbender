/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package utils

import (
	"crypto/rand"
	"encoding/hex"
)

// RandomHex generates a random hex string of specified length
func RandomHex(n int) (string, error) {
	b := make([]byte, (n+1)/2)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b)[:n], nil
}
