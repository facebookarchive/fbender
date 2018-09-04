/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package tester

import (
	"errors"
)

// Comparator allows to compare given values
type Comparator interface {
	Compare(x, y float64) bool
	Name() string
}

type comparator struct {
	repr string
	cmp  func(x, y float64) bool
}

func (c *comparator) Compare(x, y float64) bool {
	return c.cmp(x, y)
}

func (c *comparator) Name() string {
	return c.repr
}

// Available comparators
var (
	LessThan    Comparator = &comparator{repr: "<", cmp: func(x, y float64) bool { return x < y }}
	GreaterThan            = &comparator{repr: ">", cmp: func(x, y float64) bool { return x > y }}
)

// Comparators is a map of comparators representation to the actual comparator
var Comparators = map[string]Comparator{
	LessThan.Name():    LessThan,
	GreaterThan.Name(): GreaterThan,
}

// ErrInvalidComparator is returned when a comparator cannot be found
var ErrInvalidComparator = errors.New("invalid comparator")

// ParseComparator returns a comparator based on the given string
func ParseComparator(name string) (Comparator, error) {
	if cmp, ok := Comparators[name]; ok {
		return cmp, nil
	}
	return nil, ErrInvalidComparator
}
