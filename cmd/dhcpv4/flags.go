/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package dhcpv4

import (
	"encoding/csv"
	"fmt"
	"strconv"
	"strings"

	"github.com/insomniacslk/dhcp/dhcpv4"
	"github.com/spf13/pflag"
)

type optionCodeSliceValue struct {
	value   dhcpv4.OptionCodeList
	changed bool
}

// NewOptionCodeSliceValue creates a new option code slice value for pflag.
func NewOptionCodeSliceValue() pflag.Value {
	return &optionCodeSliceValue{
		changed: false,
	}
}

func readAsCSV(val string) ([]string, error) {
	if val == "" {
		return []string{}, nil
	}
	stringReader := strings.NewReader(val)
	csvReader := csv.NewReader(stringReader)
	return csvReader.Read()
}

func (s *optionCodeSliceValue) Set(value string) error {
	values, err := readAsCSV(value)
	if err != nil {
		return err
	}
	var buf []byte
	var optcodes dhcpv4.OptionCodeList
	for _, v := range values {
		optcode, err := strconv.ParseUint(v, 10, 8)
		if err != nil {
			return err
		}
		buf = append(buf, byte(uint8(optcode)))
	}
	err = optcodes.FromBytes(buf)
	if err != nil {
		return err
	}
	if !s.changed {
		s.value = optcodes
	} else {
		s.value.Add(optcodes...)
	}
	s.changed = true
	return nil
}

func (s *optionCodeSliceValue) Type() string {
	return "optionCodeSlice"
}

func (s *optionCodeSliceValue) String() string {
	return s.value.String()
}

// GetOptionCodes returns an option code slice from a pflag set
func GetOptionCodes(f *pflag.FlagSet, name string) (dhcpv4.OptionCodeList, error) {
	flag := f.Lookup(name)
	if flag == nil {
		return nil, fmt.Errorf("flag %s accessed but not defined", name)
	}
	return GetOptionCodesValue(flag.Value)
}

// GetOptionCodesValue returns an option code slice from a pflag value
func GetOptionCodesValue(v pflag.Value) (dhcpv4.OptionCodeList, error) {
	if optcodes, ok := v.(*optionCodeSliceValue); ok {
		return optcodes.value, nil
	}
	return nil, fmt.Errorf("trying to get option codes value of flag of type %s", v.Type())
}
