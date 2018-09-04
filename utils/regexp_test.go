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

func TestFindStringSubmatchMap(t *testing.T) {
	r := utils.MustCompile(`^(?P<name>[A-Z][a-z]*) (?P<surname>[A-Z][a-z]*)$`)
	// this should match
	match := r.FindStringSubmatchMap("Mikolaj Walczak")
	expected := map[string]string{"name": "Mikolaj", "surname": "Walczak"}
	assert.Equal(t, expected, match)
	// this should fail
	match = r.FindStringSubmatchMap("12345")
	expected = map[string]string{}
	assert.Equal(t, expected, match)

	// With optional fields
	r = utils.MustCompile(`^(?P<name>[A-Z][a-z]*)( (?P<surname>[A-Z][a-z]*))?$`)
	// this should still match
	match = r.FindStringSubmatchMap("Mikolaj Walczak")
	expected = map[string]string{"name": "Mikolaj", "surname": "Walczak"}
	assert.Equal(t, expected, match)
	// this should also match and have an empty surname
	match = r.FindStringSubmatchMap("Mikolaj")
	expected = map[string]string{"name": "Mikolaj", "surname": ""}
	assert.Equal(t, expected, match)
	// this should fail
	match = r.FindStringSubmatchMap("12345")
	expected = map[string]string{}
	assert.Equal(t, expected, match)
}
