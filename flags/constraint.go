/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package flags

import (
	"encoding/csv"
	"fmt"
	"strings"

	"github.com/facebookincubator/fbender/tester"
	"github.com/spf13/pflag"
)

// ConstraintSliceValue is a pflag value storing constraints.
type ConstraintSliceValue struct {
	Parsers []tester.MetricParser

	value   *[]*tester.Constraint
	changed bool
}

// NewConstraintSliceValue creates a new constraint slice value for pflag.
func NewConstraintSliceValue(parsers ...tester.MetricParser) *ConstraintSliceValue {
	v := []*tester.Constraint{}

	return &ConstraintSliceValue{
		Parsers: parsers,
		value:   &v,
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

// Set validates given string given constraints and parses them to constraint
// structures using metric parsers.
func (c *ConstraintSliceValue) Set(value string) error {
	values, err := readAsCSV(value)
	if err != nil {
		return err
	}

	constraints := []*tester.Constraint{}

	for _, v := range values {
		constraint, err := tester.ParseConstraint(v, c.Parsers...)
		if err != nil {
			return fmt.Errorf("error parsing constraint %q: %w", v, err)
		}

		constraints = append(constraints, constraint)
	}

	if !c.changed {
		*c.value = constraints
	} else {
		*c.value = append(*c.value, constraints...)
	}

	c.changed = true

	return nil
}

// Type returns the ConstraintSliceValue Type.
func (c *ConstraintSliceValue) Type() string {
	return "constraints"
}

func (c *ConstraintSliceValue) String() string {
	return fmt.Sprintf("%+v", *c.value)
}

// GetConstraints returns a constraints from a pflag set.
func GetConstraints(f *pflag.FlagSet, name string) ([]*tester.Constraint, error) {
	flag := f.Lookup(name)
	if flag == nil {
		return nil, fmt.Errorf("%w: %q", ErrUndefined, name)
	}

	return GetConstraintsValue(flag.Value)
}

// GetConstraintsValue returns a constraints from a pflag value.
func GetConstraintsValue(v pflag.Value) ([]*tester.Constraint, error) {
	if constraints, ok := v.(*ConstraintSliceValue); ok {
		return *constraints.value, nil
	}

	return nil, fmt.Errorf("%w, want: constraints, got: %s", ErrInvalidType, v.Type())
}
