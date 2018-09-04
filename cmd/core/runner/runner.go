/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package runner

import (
	"runtime"

	"github.com/gosuri/uiprogress"
	"github.com/pinterest/bender"
	"github.com/pinterest/bender/hist"
	"github.com/sirupsen/logrus"

	"github.com/facebookincubator/fbender/cmd/core/options"
	"github.com/facebookincubator/fbender/log"
	"github.com/facebookincubator/fbender/recorders"
	"github.com/facebookincubator/fbender/tester"
	"github.com/facebookincubator/fbender/utils"
)

// RequestGenerator is used to generate requests
type RequestGenerator func(i int) interface{}

// Params represents test parameters for the runner
type Params struct {
	Tester           tester.Tester
	RequestGenerator RequestGenerator
}

// runner groups fields used in both runners
type runner struct {
	requests chan interface{}
	recorder chan interface{}

	recorders []bender.Recorder
	histogram *hist.Histogram
	progress  *uiprogress.Progress
	bar       *uiprogress.Bar

	Params *Params
}

// reset "frees" all the runner fields
func (r *runner) reset() {
	r.requests = nil
	r.recorder = nil
	r.recorders = nil
	r.histogram = nil
	r.progress = nil
	r.bar = nil
}

// Before initializes all common fields
func (r *runner) Before(test int, opts interface{}) error {
	o, ok := opts.(*options.Options)
	if !ok {
		return tester.ErrInvalidOptions
	}

	cancel := utils.NewBackgroundSpinner("Cleaning up the memory", 0)
	r.reset()
	runtime.GC()
	cancel()

	cancel = utils.NewBackgroundSpinner("Preparing the test", 0)
	r.recorder = make(chan interface{}, o.BufferSize)
	r.recorders = []bender.Recorder{
		recorders.NewLogrusRecorder(logrus.StandardLogger(), logrus.Fields{"test": test}),
	}
	r.recorders = append(r.recorders, o.Recorders...)
	if !o.NoStatistics {
		r.histogram = hist.NewHistogram(2*int(o.Timeout), int(o.Unit))
		r.recorders = append(r.recorders, bender.NewHistogramRecorder(r.histogram))
	}
	cancel()

	log.Printf("Running test: %d\n", test)
	return nil
}

// After cleans up after the test
func (r *runner) After(test int, options interface{}) {
	if r.histogram != nil {
		log.Printf("%s", r.histogram.String())
	}
}

// Tester returns the protocol tester
func (r *runner) Tester() tester.Tester {
	return r.Params.Tester
}

// Requests returns the requests channel
func (r *runner) Requests() chan interface{} {
	return r.requests
}

// Recorder returns the recoreder
func (r *runner) Recorder() chan interface{} {
	return r.recorder
}

// Recorders returns a list of recoreders
func (r *runner) Recorders() []bender.Recorder {
	return r.recorders
}
