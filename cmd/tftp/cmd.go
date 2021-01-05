/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package tftp

import (
	"github.com/facebookincubator/fbender/cmd/core"
)

//nolint:gochecknoglobals
var template = &core.CommandTemplate{
	Name:  "tftp",
	Short: "Test TFTP",
	Long: `
The specified timeout applies to a single datagram in a tftp transfer rather
than to the whole session.

Target: ipv4:port, [ipv6]:port, hostname:port.

Input format: "Filename octet" or "Filename netascii"
  /my/file octet
  /my/otherfile octet
  /another netascii`,
	Fixed: `  fbender tftp {test} fixed -t $TARGET 10 20
  fbender tftp {test} fixed -t $TARGET -d 5m 50`,
	Constraints: `  fbender tftp {test} constraints -t $TARGET -b 1500 -c "AVG(latency)<10" 20
  fbender tftp {test} constraints -t $TARGET -g ^10 -c "MAX(errors)<10" 40`,
}

// Command is the TFTP subcommand.
//nolint:gochecknoglobals
var Command = core.NewTestCommand(template, params)

//nolint:gochecknoinits
func init() {
	Command.PersistentFlags().IntP("blocksize", "s", 512, "blocksize option as in RFC2348")
}
