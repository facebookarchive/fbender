/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package http

import (
	"github.com/facebookincubator/fbender/cmd/core"
)

var template = &core.CommandTemplate{
	Name:  "http",
	Short: "Test HTTP",
	Long: `
Target: ipv4, ipv4:port, ipv6, [ipv6]:port, hostname, hostname:port.

Input format: "GET RelativeURL" or "POST RelativeURL FormData"
  GET index.html
  GET /
  POST echo message=Hello
  POST echo/ message=Hello&name=Mikolaj`,
	Fixed: `  fbender http {test} fixed -t $TARGET 10 20
  fbender http {test} fixed -t $TARGET -s -d 5m 50`,
	Constraints: `  fbender http {test} constraints -t $TARGET -s -c "AVG(latency)<10" 20
  fbender http {test} constraints -t $TARGET -g ^10 -c "MAX(errors)<10" 40`,
}

// Command is the HTTP subcommand
var Command = core.NewTestCommand(template, params)

func init() {
	Command.PersistentFlags().BoolP("ssl", "s", false, "enable ssl (use HTTPS)")
}
