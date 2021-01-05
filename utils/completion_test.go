/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package utils_test

import (
	"strings"
	"testing"

	"github.com/facebookincubator/fbender/utils"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBashCompletion(t *testing.T) {
	const (
		fname = "__fbender_test_handle_flag"
		fbody = `COMPREPLY=($(compgen -W "42 24" -- "${cur}"))`
	)

	c := &cobra.Command{}
	c.Flags().Int("myint", 0, "set myint")

	// Check if the completion function is appended
	err := utils.BashCompletion(c, c.Flags(), "myint", fname, fbody)
	require.NoError(t, err)
	assert.Contains(t, c.BashCompletionFunction, fname)
	assert.Equal(t, `
__fbender_test_handle_flag() {
	COMPREPLY=($(compgen -W "42 24" -- "${cur}"))
}`, c.BashCompletionFunction)

	// Check if the flag annotation has been added
	f := c.Flags().Lookup("myint")
	require.NotNil(t, f)
	require.Contains(t, f.Annotations, "cobra_annotation_bash_completion_custom")
	assert.Equal(t, []string{fname}, f.Annotations["cobra_annotation_bash_completion_custom"])

	// Check if the completion function is appended only once
	err = utils.BashCompletion(c, c.Flags(), "myint", fname, fbody)
	require.NoError(t, err)
	assert.Contains(t, c.BashCompletionFunction, fname)
	count := strings.Count(c.BashCompletionFunction, fname)
	assert.Equal(t, 1, count, "Completion function should be added only once")
}
