/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package dns_test

import (
	"strings"
	"testing"

	dnsflags "github.com/facebookincubator/fbender/cmd/dns"
	"github.com/facebookincubator/fbender/flags"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/suite"
)

type ProtocolValueTestSuite struct {
	suite.Suite
	value pflag.Value
}

func (s *ProtocolValueTestSuite) SetupTest() {
	s.value = dnsflags.NewProtocolValue()
	s.Require().NotNil(s.value)
}

func (s *ProtocolValueTestSuite) TestSet_NoErrors() {
	err := s.value.Set("udp")
	s.Require().NoError(err)

	v, err := dnsflags.GetProtocolValue(s.value)
	s.Require().NoError(err)
	s.Assert().Equal("udp", v)

	err = s.value.Set("tcp")
	s.Require().NoError(err)

	v, err = dnsflags.GetProtocolValue(s.value)
	s.Require().NoError(err)
	s.Assert().Equal("tcp", v)
}

func (s *ProtocolValueTestSuite) TestSet_Errors() {
	// Save original flag value
	o, err := dnsflags.GetProtocolValue(s.value)
	s.Require().NoError(err)

	// Try invalid value
	err = s.value.Set("unknown")
	s.Assert().ErrorIs(err, dnsflags.ErrInvalidProtocol)
	s.Assert().EqualError(err, "invalid protocol, want: \"udp\" or \"tcp\", got: \"unknown\"")

	// The value shouldn't change
	v, err := dnsflags.GetProtocolValue(s.value)
	s.Require().NoError(err)
	s.Assert().Equal(o, v)
}

func (s *ProtocolValueTestSuite) TestType() {
	s.Assert().Equal("protocol", s.value.Type())
}

func (s *ProtocolValueTestSuite) TestGetProtocol() {
	f := pflag.NewFlagSet("Test FlagSet", pflag.ExitOnError)
	f.Var(s.value, "protocol", "set protocol")

	err := s.value.Set("tcp")
	s.Require().NoError(err)

	v, err := dnsflags.GetProtocol(f, "protocol")
	s.Require().NoError(err)
	s.Assert().Equal("tcp", v)

	// Check error when flag does not exist
	_, err = dnsflags.GetProtocol(f, "nonexistent")
	s.Assert().ErrorIs(err, flags.ErrUndefined)
	s.Assert().EqualError(err, "flag accessed but not defined: \"nonexistent\"")

	// Check error when value is of different type
	f.Int("myint", 0, "set myint")
	_, err = dnsflags.GetProtocol(f, "myint")
	s.Assert().ErrorIs(err, flags.ErrInvalidType)
	s.Assert().EqualError(err, "accessed flag type does not match, want: protocol, got: int")
}

func (s *ProtocolValueTestSuite) TestGetProtocolValue() {
	err := s.value.Set("tcp")
	s.Require().NoError(err)

	v, err := dnsflags.GetProtocolValue(s.value)
	s.Require().NoError(err)
	s.Assert().Equal("tcp", v)

	// Check error when value is of different type
	f := pflag.NewFlagSet("Test FlagSet", pflag.ExitOnError)
	f.Int("myint", 0, "set myint")

	flag := f.Lookup("myint")
	s.Require().NotNil(flag)

	_, err = dnsflags.GetProtocolValue(flag.Value)
	s.Assert().ErrorIs(err, flags.ErrInvalidType)
	s.Assert().EqualError(err, "accessed flag type does not match, want: protocol, got: int")
}

func (s *ProtocolValueTestSuite) TestBashCompletionProtocol() {
	c := &cobra.Command{}
	f := c.Flags().VarPF(s.value, "protocol", "p", "set protocol")

	// Check if the complete function is appended
	err := dnsflags.BashCompletionProtocol(c, c.Flags(), "protocol")
	s.Require().NoError(err)
	s.Assert().Contains(c.BashCompletionFunction, "__fbender_handle_dns_protocol_flag")

	// Check if the flag has the bash
	s.Require().Contains(f.Annotations, "cobra_annotation_bash_completion_custom")
	s.Assert().Equal([]string{"__fbender_handle_dns_protocol_flag"},
		f.Annotations["cobra_annotation_bash_completion_custom"])

	// Check if the function is appended only once
	err = dnsflags.BashCompletionProtocol(c, c.Flags(), "protocol")
	s.Require().NoError(err)

	count := strings.Count(c.BashCompletionFunction, "__fbender_handle_dns_protocol_flag")
	s.Assert().Equal(1, count, "Completion function should be added only once")
}

func TestProtocolValueTestSuite(t *testing.T) {
	suite.Run(t, new(ProtocolValueTestSuite))
}
