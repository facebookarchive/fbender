/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package dns

import (
	"github.com/facebookincubator/fbender/cmd/core"
)

//nolint:gochecknoglobals
var template = &core.CommandTemplate{
	Name:  "dns",
	Short: "Test DNS",
	Long: `
Queries may be prefixed with timestamp and a random 16-character hex to avoid
hitting the cache. In bash this could have been achieved by running:
  $(date +%s).$(openssl rand -hex 16).domain

Target: ipv4, ipv6, hostname, ipv4:port, [ipv6]:port, hostname:port.
The port defaults to 53.

Input format: "Domain QType [Rcode]"
  example.com AAAA
  other.example.com TXT NOERROR
  mail.example.com MX
	www.doesnotexist.co.uk NXDOMAIN`,
	Fixed: `  fbender dns {test} fixed -t $TARGET 10 20
  fbender dns {test} fixed -t $TARGET -r -d 5m 50`,
	Constraints: `  fbender dns {test} constraints -t $TARGET -r -c "AVG(latency)<10" 20
  fbender dns {test} constraints -t $TARGET -g ^10 -c "MAX(errors)<10" 40`,
}

// Command is the DNS subcommand.
//nolint:gochecknoglobals
var Command = core.NewTestCommand(template, params)

//nolint:gochecknoinits
func init() {
	Command.PersistentFlags().BoolP("randomize", "r", false, "randomize queries with timestamp and a random hex")
	core.DeferPostInit(postinit)
}

func postinit() {
	protocol := NewProtocolValue()

	Command.PersistentFlags().VarP(protocol, "protocol", "p", "protocol used for DNS queries (udp|tcp)")

	if err := BashCompletionProtocol(Command, Command.PersistentFlags(), "protocol"); err != nil {
		panic(err)
	}
}
