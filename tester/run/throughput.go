/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package run

import (
	"time"

	"github.com/pinterest/bender"

	"github.com/facebookincubator/fbender/tester"
)

// LoadTestThroughputFixed runs predefined set of throughput tests
func LoadTestThroughputFixed(r tester.ThroughputRunner, o interface{}, qs ...tester.QPS) error {
	t := r.Tester()
	if err := t.Before(o); err != nil {
		return err
	}
	defer t.After(o)

	for _, qps := range qs {
		if err := loadTestThroughput(r, t, o, qps); err != nil {
			return err
		}
	}
	return nil
}

// LoadTestThroughputConstraints automatically tries to find a breakpoint based on provided constraints checks
func LoadTestThroughputConstraints(r tester.ThroughputRunner, o interface{}, start tester.QPS, g tester.Growth,
	cs ...*tester.Constraint) error {

	t := r.Tester()
	if err := t.Before(o); err != nil {
		return err
	}
	defer t.After(o)

	qps := start
	for qps > 0 {
		startTime := time.Now()
		if err := loadTestThroughput(r, t, o, qps); err != nil {
			return err
		}
		duration := time.Since(startTime)
		if checkConstraints(startTime, duration, cs...) {
			qps = g.OnSuccess(qps)
		} else {
			qps = g.OnFail(qps)
		}
	}
	return nil
}

// loadTestThroughput runs a single test for a desired QPS.
func loadTestThroughput(r tester.ThroughputRunner, t tester.Tester, o interface{}, qps tester.QPS) error {
	if err := t.BeforeEach(o); err != nil {
		return err
	}
	defer t.AfterEach(o)

	if err := r.Before(qps, o); err != nil {
		return err
	}
	defer r.After(qps, o)

	executor, err := t.RequestExecutor(o)
	if err != nil {
		return err
	}
	bender.LoadTestThroughput(r.Intervals(), r.Requests(), executor, r.Recorder())
	bender.Record(r.Recorder(), r.Recorders()...)
	return nil
}
