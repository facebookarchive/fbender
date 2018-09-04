/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package dhcpv4

import (
	"net"

	"github.com/insomniacslk/dhcp/dhcpv4"
	"github.com/insomniacslk/dhcp/dhcpv4/async"
	"github.com/insomniacslk/dhcp/iana"
	"github.com/spf13/cobra"

	"github.com/facebookincubator/fbender/cmd/core/input"
	"github.com/facebookincubator/fbender/cmd/core/options"
	"github.com/facebookincubator/fbender/cmd/core/runner"
	tester "github.com/facebookincubator/fbender/tester/dhcpv4"
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
		Target:     utils.WithDefaultPort(o.Target, async.DefaultServerPort),
		Timeout:    o.Timeout,
		BufferSize: o.BufferSize,
	}
	return &runner.Params{Tester: t, RequestGenerator: r}, nil
}

func inputTransformer(optionCodes []dhcpv4.OptionCode) input.Transformer {
	defaultCodes := []dhcpv4.OptionCode{
		dhcpv4.OptionSubnetMask,
		dhcpv4.OptionRouter,
		dhcpv4.OptionDomainName,
		dhcpv4.OptionDomainNameServer,
	}
	return func(input string) (interface{}, error) {
		mac, err := net.ParseMAC(input)
		if err != nil {
			return nil, err
		}
		discover, err := dhcpv4.New()
		if err != nil {
			return nil, err
		}
		discover.SetHwType(iana.HwTypeEthernet)
		discover.SetHwAddrLen(uint8(len(mac)))
		discover.SetClientHwAddr(mac)
		discover.AddOption(&dhcpv4.OptMessageType{MessageType: dhcpv4.MessageTypeDiscover})
		optionCodes = append(optionCodes, defaultCodes...)
		discover = dhcpv4.WithRequestedOptions(optionCodes...)(discover)
		return discover, nil
	}
}
