/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package runner

import (
	"time"

	"github.com/facebookincubator/fbender/cmd/core/options"
	"github.com/facebookincubator/fbender/recorders"
	"github.com/facebookincubator/fbender/tester"
	"github.com/pinterest/bender"
)

// ThroughputRunner is a test runner for load test throughput commands.
type ThroughputRunner struct {
	runner
	intervals bender.IntervalGenerator
}

// NewThroughputRunner returns new ThroughputRunner.
func NewThroughputRunner(params *Params) *ThroughputRunner {
	return &ThroughputRunner{
		runner: runner{
			Params: params,
		},
	}
}

// Before prepares requests, recorders and interval generator.
func (r *ThroughputRunner) Before(qps tester.QPS, opts interface{}) error {
	if err := r.runner.Before(qps, opts); err != nil {
		return err
	}

	o, ok := opts.(*options.Options)
	if !ok {
		return tester.ErrInvalidOptions
	}

	count := int(float64(qps) * float64(o.Duration/time.Second))
	r.intervals = o.Distribution(float64(qps))

	r.requests = make(chan interface{}, o.BufferSize)

	go func() {
		for i := 0; i < count; i++ {
			r.requests <- r.Params.RequestGenerator(i)
		}
		close(r.requests)
	}()

	r.progress, r.bar = recorders.NewLoadTestProgress(count)
	r.progress.Start()
	r.recorders = append(r.recorders, recorders.NewProgressBarRecorder(r.bar))

	return nil
}

// After cleans up after the test.
func (r *ThroughputRunner) After(test int, options interface{}) {
	r.progress.Stop()
	r.runner.After(test, options)
}

// Intervals returns the interval generator.
func (r *ThroughputRunner) Intervals() bender.IntervalGenerator {
	return r.intervals
}
