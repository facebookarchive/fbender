/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package tftp

import (
	"fmt"
	"strings"

	"github.com/pinterest/bender/tftp"
	"github.com/spf13/cobra"

	"github.com/facebookincubator/fbender/cmd/core/input"
	"github.com/facebookincubator/fbender/cmd/core/options"
	"github.com/facebookincubator/fbender/cmd/core/runner"
	tester "github.com/facebookincubator/fbender/tester/tftp"
	"github.com/facebookincubator/fbender/utils"
)

// DefaultServerPort is a default tftp server port
const DefaultServerPort = 69

func params(cmd *cobra.Command, o *options.Options) (*runner.Params, error) {
	blocksize, err := cmd.Flags().GetInt("blocksize")
	if err != nil {
		return nil, err
	}
	r, err := input.NewRequestGenerator(o.Input, inputTransformer)
	if err != nil {
		return nil, err
	}
	t := &tester.Tester{
		Target:    utils.WithDefaultPort(o.Target, DefaultServerPort),
		Timeout:   o.Timeout,
		BlockSize: blocksize,
	}
	return &runner.Params{Tester: t, RequestGenerator: r}, nil
}

func inputTransformer(input string) (interface{}, error) {
	i := strings.Index(input, " ")
	if i < 0 {
		return nil, fmt.Errorf("input must have a format of 'File Mode' got '%s'", input)
	}
	filename, mode := input[:i], input[i+1:]
	if mode != "octet" && mode != "netascii" {
		return nil, fmt.Errorf("invalid mode '%s', want (octet|netascii)", mode)
	}
	return &tftp.Request{
		Filename: filename,
		Mode:     tftp.RequestMode(mode),
	}, nil
}
