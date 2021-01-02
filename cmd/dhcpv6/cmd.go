/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package dhcpv6

import (
	"github.com/facebookincubator/fbender/cmd/core"
)

//nolint:gochecknoglobals
var template = &core.CommandTemplate{
	Name:  "dhcpv6",
	Short: "Test DHCPv6",
	Long: `
Target: ipv4, ipv6, hostname, ipv4:port, [ipv6]:port, hostname:port.
Port defaults to 547, unless you know what you're doing you shouldn't change it.

Input format: "DeviceMAC"
  01:23:45:67:89:ab
  E3:63:BD:7B:D2:2C
  c8:6c:2c:47:96:fd`,
	Fixed: `  fbender dhcpv6 {test} fixed -t $TARGET 10 20
  fbender dhcpv6 {test} fixed -t $TARGET -d 5m 50`,
	Constraints: `  fbender dhcpv6 {test} constraints -t $TARGET -c "AVG(latency)<10" 20
  fbender dhcpv6 {test} constraints -t $TARGET -g ^10 -c "MAX(errors)<10" 40`,
}

// Command is the TFTP subcommand.
//nolint:gochecknoglobals
var Command = core.NewTestCommand(template, params)

//nolint:gochecknoinits
func init() {
	optionCodes := NewOptionCodeSliceValue()
	Command.PersistentFlags().VarP(optionCodes, "oro", "r", "dhcpv6 requested options (ORO)")
	Command.Aliases = []string{"dhcp6"}
}
