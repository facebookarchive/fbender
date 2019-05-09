/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package dhcpv6

import (
	"net"

	"github.com/insomniacslk/dhcp/dhcpv6"
	"github.com/spf13/cobra"

	"github.com/facebookincubator/fbender/cmd/core/input"
	"github.com/facebookincubator/fbender/cmd/core/options"
	"github.com/facebookincubator/fbender/cmd/core/runner"
	tester "github.com/facebookincubator/fbender/tester/dhcpv6"
	"github.com/facebookincubator/fbender/utils"
)

func params(cmd *cobra.Command, o *options.Options) (*runner.Params, error) {
	optionCodes, err := GetOptionCodes(cmd.Flags(), "oro")
	if err != nil {
		return nil, err
	}
	r, err := input.NewRequestGenerator(o.Input, inputTransformer(optionCodes))
	if err != nil {
		return nil, err
	}
	t := &tester.Tester{
		Target:     utils.WithDefaultPort(o.Target, dhcpv6.DefaultServerPort),
		Timeout:    o.Timeout,
		BufferSize: o.BufferSize,
	}
	return &runner.Params{Tester: t, RequestGenerator: r}, nil
}

func inputTransformer(optionCodes []dhcpv6.OptionCode) input.Transformer {
	return func(input string) (interface{}, error) {
		mac, err := net.ParseMAC(input)
		if err != nil {
			return nil, err
		}
		return dhcpv6.NewSolicit(mac, dhcpv6.WithRequestedOptions(optionCodes...))
	}
}
