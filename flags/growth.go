/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package flags

import (
	"fmt"

	"github.com/facebookincubator/fbender/tester"
	"github.com/spf13/pflag"
)

// GrowthValue represents growth flag value.
type GrowthValue struct {
	Growth tester.Growth
}

func (g *GrowthValue) String() string {
	return g.Growth.String()
}

// Set validates a given growth and saves it.
func (g *GrowthValue) Set(value string) error {
	var err error

	g.Growth, err = tester.ParseGrowth(value)
	if err != nil {
		return fmt.Errorf("error parsing growth %q: %w", value, err)
	}

	return nil
}

// Type returns a growth value type.
func (g *GrowthValue) Type() string {
	return "growth"
}

// GetGrowth returns a growth from a pflag set.
func GetGrowth(f *pflag.FlagSet, name string) (tester.Growth, error) {
	flag := f.Lookup(name)
	if flag == nil {
		return nil, fmt.Errorf("%w: %q", ErrUndefined, name)
	}

	return GetGrowthValue(flag.Value)
}

// GetGrowthValue returns a growth from a pflag value.
func GetGrowthValue(v pflag.Value) (tester.Growth, error) {
	if growth, ok := v.(*GrowthValue); ok {
		return growth.Growth, nil
	}

	return nil, fmt.Errorf("%w, want: *GrowthValue, got: %T", ErrInvalidType, v)
}
