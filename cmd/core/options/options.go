/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package options

import (
	"time"

	"github.com/pinterest/bender"

	"github.com/facebookincubator/fbender/tester"
)

// Options represents common options for the Commands
type Options struct {
	Target   string
	Duration time.Duration
	Tests    []int
	Start    int

	Input string

	BufferSize   int
	Timeout      time.Duration
	Distribution func(float64) bender.IntervalGenerator
	Unit         time.Duration
	NoStatistics bool

	Constraints []*tester.Constraint
	Growth      tester.Growth

	Recorders []bender.Recorder
}

// NewOptions returns new options
func NewOptions() *Options {
	return &Options{
		Tests:       []int{},
		Constraints: []*tester.Constraint{},
		Recorders:   []bender.Recorder{},
	}
}

// GetUnit returns a unit used in tests
func (o *Options) GetUnit() time.Duration {
	return o.Unit
}

// AddRecorder adds a recorder to options
func (o *Options) AddRecorder(recorder bender.Recorder) {
	o.Recorders = append(o.Recorders, recorder)
}
