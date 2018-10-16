/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package dns

import (
	"fmt"
	"strings"
	"time"

	"github.com/miekg/dns"
	"github.com/spf13/cobra"

	"github.com/facebookincubator/fbender/cmd/core/input"
	"github.com/facebookincubator/fbender/cmd/core/options"
	"github.com/facebookincubator/fbender/cmd/core/runner"
	tester "github.com/facebookincubator/fbender/tester/dns"
	"github.com/facebookincubator/fbender/utils"
)

// DefaultServerPort is a default dns server port
const DefaultServerPort = 53

func params(cmd *cobra.Command, o *options.Options) (*runner.Params, error) {
	randomize, err := cmd.Flags().GetBool("randomize")
	if err != nil {
		return nil, err
	}
	protocol, err := cmd.Flags().GetString("protocol")
	if err != nil {
		return nil, err
	}
	if protocol != "tcp" && protocol != "udp" {
		return nil, fmt.Errorf("unknown protocol (%s), expecting one of 'tcp' or 'udp'", protocol)
	}
	r, err := input.NewRequestGenerator(o.Input, inputTransformer, getModifiers(randomize)...)
	if err != nil {
		return nil, err
	}
	t := &tester.Tester{
		Target:   utils.WithDefaultPort(o.Target, DefaultServerPort),
		Timeout:  o.Timeout,
		Protocol: protocol,
	}
	return &runner.Params{Tester: t, RequestGenerator: r}, nil
}

func inputTransformer(input string) (interface{}, error) {
	var domain, typeString string
	n, err := fmt.Sscanf(input, "%s %s", &domain, &typeString)
	if err != nil || n < 2 {
		return nil, fmt.Errorf("invalid input: %q, want: \"Domain QType\"", input)
	}
	msgTyp, ok := dns.StringToType[strings.ToUpper(typeString)]
	if !ok {
		return nil, fmt.Errorf("invalid QType: %q", typeString)
	}
	msg := new(dns.Msg)
	msg.SetQuestion(dns.Fqdn(domain), msgTyp)
	return msg, nil
}

func getModifiers(randomize bool) []input.Modifier {
	if randomize {
		return []input.Modifier{randomPrefixModifier}
	}
	return []input.Modifier{}
}

const prefixLength = 16

func randomPrefixModifier(request interface{}) (interface{}, error) {
	msg, ok := request.(*dns.Msg)
	if !ok {
		return nil, fmt.Errorf("invalid request type: %T, want: *dns.Msg", request)
	}
	hex, err := utils.RandomHex(prefixLength)
	if err != nil {
		return nil, err
	}
	// Create a new message so we don't destroy the original to avoid recursive prefixing
	modified := new(dns.Msg)
	domain := fmt.Sprintf("%d.%s.%s", time.Now().Unix(), hex, msg.Question[0].Name)
	msgTyp := msg.Question[0].Qtype
	modified.SetQuestion(dns.Fqdn(domain), msgTyp)
	return modified, nil
}
