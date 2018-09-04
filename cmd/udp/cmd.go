/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package udp

import (
	"github.com/facebookincubator/fbender/cmd/core"
)

var template = &core.CommandTemplate{
	Name:  "udp",
	Short: "Test UDP",
	Long: `
Target: ipv4, ipv6, hostname.

Input format: "DstPort Base64EncodedeData"
  2545 TG9yZW0=
  7346 aXBzdW0gZG9sb3Igc2l0
  5012 YW1ldCBpbg==`,
	Fixed: `  fbender udp {test} fixed -t $TARGET 10 20
  fbender udp {test} fixed -t $TARGET -d 5m 50`,
	Constraints: `  fbender udp {test} constraints -t $TARGET -c "AVG(latency)<10" 20
  fbender udp {test} constraints -t $TARGET -g ^10 -c "MAX(errors)<10" 40`,
}

// Command is the UDP subcommand
var Command = core.NewTestCommand(template, params)
