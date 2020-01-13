/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package dhcpv6_test

import (
	"testing"

	"github.com/insomniacslk/dhcp/dhcpv6"
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/suite"

	flags "github.com/facebookincubator/fbender/cmd/dhcpv6"
)

type OptionCodeSliceValueTestSuite struct {
	suite.Suite
	value pflag.Value
}

func (s *OptionCodeSliceValueTestSuite) SetupTest() {
	s.value = flags.NewOptionCodeSliceValue()
	s.Require().NotNil(s.value)
}

func (s *OptionCodeSliceValueTestSuite) TestSet_NoErrors() {
	err := s.value.Set("1,2")
	s.Require().NoError(err)
	v, err := flags.GetOptionCodesValue(s.value)
	s.Require().NoError(err)
	o := []dhcpv6.OptionCode{dhcpv6.OptionClientID, dhcpv6.OptionServerID}
	s.Assert().Equal(o, v)

	// Check if consecutive calls append values
	err = s.value.Set("3")
	s.Require().NoError(err)
	v, err = flags.GetOptionCodesValue(s.value)
	s.Require().NoError(err)
	o = append(o, dhcpv6.OptionIANA)
	s.Assert().Equal(o, v)

	err = s.value.Set("4,5")
	s.Require().NoError(err)
	v, err = flags.GetOptionCodesValue(s.value)
	s.Require().NoError(err)
	o = append(o, dhcpv6.OptionIATA, dhcpv6.OptionIAAddr)
	s.Assert().Equal(o, v)
}

func (s *OptionCodeSliceValueTestSuite) TestSet_Errors() {
	// Errors - single value
	err := s.value.Set("notanumber")
	s.Assert().EqualError(err, "strconv.ParseUint: parsing \"notanumber\": invalid syntax")
	v, err := flags.GetOptionCodesValue(s.value)
	s.Require().NoError(err)
	s.Assert().Empty(v)

	err = s.value.Set("42.5")
	s.Assert().EqualError(err, "strconv.ParseUint: parsing \"42.5\": invalid syntax")
	v, err = flags.GetOptionCodesValue(s.value)
	s.Require().NoError(err)
	s.Assert().Empty(v)

	err = s.value.Set("-10")
	s.Assert().EqualError(err, "strconv.ParseUint: parsing \"-10\": invalid syntax")
	v, err = flags.GetOptionCodesValue(s.value)
	s.Require().NoError(err)
	s.Assert().Empty(v)

	err = s.value.Set("256")
	s.Assert().EqualError(err, "strconv.ParseUint: parsing \"256\": value out of range")
	v, err = flags.GetOptionCodesValue(s.value)
	s.Require().NoError(err)
	s.Assert().Empty(v)

	// Errors - multiple values
	err = s.value.Set("42,notanumber")
	s.Assert().EqualError(err, "strconv.ParseUint: parsing \"notanumber\": invalid syntax")
	v, err = flags.GetOptionCodesValue(s.value)
	s.Require().NoError(err)
	s.Assert().Empty(v)

	err = s.value.Set("42,42.5,notanumber")
	s.Assert().EqualError(err, "strconv.ParseUint: parsing \"42.5\": invalid syntax")
	v, err = flags.GetOptionCodesValue(s.value)
	s.Require().NoError(err)
	s.Assert().Empty(v)
}

func (s *OptionCodeSliceValueTestSuite) TestType() {
	s.Assert().Equal("optionCodeSlice", s.value.Type())
}

func (s *OptionCodeSliceValueTestSuite) TestString_Known() {
	// No options
	s.Assert().Equal("", s.value.String())

	// Single option
	err := s.value.Set("1")
	s.Require().NoError(err)
	v := s.value.String()
	s.Assert().Equal("Client Identifier", v)

	// Multiple options
	err = s.value.Set("2,3")
	s.Require().NoError(err)
	v = s.value.String()
	s.Assert().Equal("Client Identifier, Server Identifier, IA_NA", v)
}

func (s *OptionCodeSliceValueTestSuite) TestString_Unknown() {
	// No options
	s.Assert().Equal("", s.value.String())

	// Single option
	err := s.value.Set("10")
	s.Require().NoError(err)
	v := s.value.String()
	s.Assert().Equal("unknown (10)", v)

	// Multiple options
	err = s.value.Set("35,255")
	s.Require().NoError(err)
	v = s.value.String()
	s.Assert().Equal("unknown (10), unknown (35), unknown (255)", v)
}

func (s *OptionCodeSliceValueTestSuite) TestGetOptionCodes() {
	f := pflag.NewFlagSet("Test FlagSet", pflag.ExitOnError)
	f.Var(s.value, "optioncodes", "set option codes")

	err := s.value.Set("39")
	s.Require().NoError(err)

	v, err := flags.GetOptionCodes(f, "optioncodes")
	s.Require().NoError(err)
	s.Assert().Equal([]dhcpv6.OptionCode{dhcpv6.OptionFQDN}, v)

	// Check error when flag does not exist
	_, err = flags.GetOptionCodes(f, "nonexistent")
	s.Assert().EqualError(err, "flag nonexistent accessed but not defined")

	// Check error when value is of different type
	f.Int("myint", 0, "set myint")
	_, err = flags.GetOptionCodes(f, "myint")
	s.Assert().EqualError(err, "trying to get option codes value of flag of type int")
}

func (s *OptionCodeSliceValueTestSuite) TestGetOptionCodesValue() {
	err := s.value.Set("39")
	s.Require().NoError(err)

	v, err := flags.GetOptionCodesValue(s.value)
	s.Require().NoError(err)
	s.Assert().Equal([]dhcpv6.OptionCode{dhcpv6.OptionFQDN}, v)

	// Check error when value is of different type
	f := pflag.NewFlagSet("Test FlagSet", pflag.ExitOnError)
	f.Int("myint", 0, "set myint")
	flag := f.Lookup("myint")
	s.Require().NotNil(flag)
	_, err = flags.GetOptionCodesValue(flag.Value)
	s.Assert().EqualError(err, "trying to get option codes value of flag of type int")
}

func TestOptionCodeSliceValueTestSuite(t *testing.T) {
	suite.Run(t, new(OptionCodeSliceValueTestSuite))
}
