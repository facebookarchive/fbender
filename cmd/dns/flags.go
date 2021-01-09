/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package dns

import (
	"errors"
	"fmt"

	"github.com/facebookincubator/fbender/flags"
	"github.com/facebookincubator/fbender/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// protocols is a set of available protocols.
//nolint:gochecknoglobals
var protocols = map[string]struct{}{
	"udp": {},
	"tcp": {},
}

// ErrInvalidProtocol is raised when an unknown protocol is set.
var ErrInvalidProtocol = errors.New("invalid protocol")

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

	return fmt.Errorf("%w, want: \"udp\" or \"tcp\", got: %q", ErrInvalidProtocol, value)
}

func (s *protocolValue) Type() string {
	return "protocol"
}

func (s *protocolValue) String() string {
	return s.value
}

// GetProtocol returns a protocol from a pflag set.
func GetProtocol(f *pflag.FlagSet, name string) (string, error) {
	flag := f.Lookup(name)
	if flag == nil {
		return "", fmt.Errorf("%w: %q", flags.ErrUndefined, name)
	}

	return GetProtocolValue(flag.Value)
}

// GetProtocolValue returns a protocol from a pflag value.
func GetProtocolValue(v pflag.Value) (string, error) {
	if protocol, ok := v.(*protocolValue); ok {
		return protocol.value, nil
	}

	return "", fmt.Errorf("%w, want: protocol, got: %s", flags.ErrInvalidType, v.Type())
}

// Bash completion function constants.
const (
	fname = "__fbender_handle_dns_protocol_flag"
	fbody = `COMPREPLY=($(compgen -W "udp tcp" -- "${cur}"))`
)

// BashCompletionProtocol adds bash completion to a protocol flag.
func BashCompletionProtocol(cmd *cobra.Command, flags *pflag.FlagSet, name string) error {
	return utils.BashCompletion(cmd, flags, name, fname, fbody)
}
