/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package dns

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/facebookincubator/fbender/utils"
)

// protocols is a set of available protocols
var protocols = map[string]struct{}{
	"udp": {},
	"tcp": {},
}

type protocolValue struct {
	value string
}

// NewProtocolValue returns new Protocol flag with a default value.
func NewProtocolValue() pflag.Value {
	return &protocolValue{value: "udp"}
}

func (s *protocolValue) Set(value string) error {
	if _, ok := protocols[value]; ok {
		s.value = value
		return nil
	}
	return fmt.Errorf("unknown protocol %q, want: \"udp\" or \"tcp\"", value)
}

func (s *protocolValue) Type() string {
	return "protocol"
}

func (s *protocolValue) String() string {
	return s.value
}

// GetProtocol returns a protocol from a pflag set
func GetProtocol(f *pflag.FlagSet, name string) (string, error) {
	flag := f.Lookup(name)
	if flag == nil {
		return "", fmt.Errorf("flag %s accessed but not defined", name)
	}
	return GetProtocolValue(flag.Value)
}

// GetProtocolValue returns a protocol from a pflag value
func GetProtocolValue(v pflag.Value) (string, error) {
	if protocol, ok := v.(*protocolValue); ok {
		return protocol.value, nil
	}
	return "", fmt.Errorf("trying to get protocol value of flag of type %s", v.Type())
}

// Bash completion function constants
const (
	fname = "__fbender_handle_dns_protocol_flag"
	fbody = `COMPREPLY=($(compgen -W "udp tcp" -- "${cur}"))`
)

// BashCompletionProtocol adds bash completion to a protocol flag
func BashCompletionProtocol(cmd *cobra.Command, flags *pflag.FlagSet, name string) error {
	return utils.BashCompletion(cmd, flags, name, fname, fbody)
}
