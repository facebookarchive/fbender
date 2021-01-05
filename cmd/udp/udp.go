/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package udp

import (
	"encoding/base64"
	"fmt"

	"github.com/facebookincubator/fbender/cmd/core/input"
	"github.com/facebookincubator/fbender/cmd/core/options"
	"github.com/facebookincubator/fbender/cmd/core/runner"
	"github.com/facebookincubator/fbender/protocols/udp"
	tester "github.com/facebookincubator/fbender/tester/udp"
	"github.com/spf13/cobra"
)

func params(cmd *cobra.Command, o *options.Options) (*runner.Params, error) {
	r, err := input.NewRequestGenerator(o.Input, inputTransformer)
	if err != nil {
		return nil, err
	}

	t := &tester.Tester{
		Target:  o.Target,
		Timeout: o.Timeout,
	}

	return &runner.Params{Tester: t, RequestGenerator: r}, nil
}

func inputTransformer(input string) (interface{}, error) {
	var encodedData string

	datagram := new(udp.Datagram)

	n, err := fmt.Sscanf(input, "%d %s", &datagram.Port, &encodedData)
	if err != nil || n < 2 {
		return nil, fmt.Errorf("invalid datagram: %q, want: \"Port Base64Payload\"", input)
	}

	datagram.Data, err = base64.StdEncoding.DecodeString(encodedData)
	if err != nil {
		return nil, err
	}

	return datagram, nil
}
