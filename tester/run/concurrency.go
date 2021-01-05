/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package run

import (
	"time"

	"github.com/facebookincubator/fbender/tester"
	"github.com/pinterest/bender"
)

// LoadTestConcurrencyFixed runs predefined set of throughput tests.
func LoadTestConcurrencyFixed(r tester.ConcurrencyRunner, o interface{}, ws ...tester.Workers) error {
	t := r.Tester()
	if err := t.Before(o); err != nil {
		return err
	}

	defer t.After(o)

	for _, workers := range ws {
		if err := loadTestConcurrency(r, t, o, workers); err != nil {
			return err
		}
	}

	return nil
}

// LoadTestConcurrencyConstraints automatically tries to find a breakpoint based on provided constraints checks.
func LoadTestConcurrencyConstraints(r tester.ConcurrencyRunner, o interface{}, start tester.Workers, g tester.Growth,
	cs ...*tester.Constraint) error {
	t := r.Tester()
	if err := t.Before(o); err != nil {
		return err
	}

	defer t.After(o)

	workers := start
	for workers > 0 {
		startTime := time.Now()

		if err := loadTestConcurrency(r, t, o, workers); err != nil {
			return err
		}

		duration := time.Since(startTime)

		if checkConstraints(startTime, duration, cs...) {
			workers = g.OnSuccess(workers)
		} else {
			workers = g.OnFail(workers)
		}
	}

	return nil
}

// loadTestConcurrency runs a single test for a desired QPS.
func loadTestConcurrency(r tester.ConcurrencyRunner, t tester.Tester, o interface{}, workers tester.Workers) error {
	if err := t.BeforeEach(o); err != nil {
		return err
	}
	defer t.AfterEach(o)

	if err := r.Before(workers, o); err != nil {
		return err
	}
	defer r.After(workers, o)

	executor, err := t.RequestExecutor(o)
	if err != nil {
		return err
	}

	bender.LoadTestConcurrency(r.WorkerSemaphore(), r.Requests(), executor, r.Recorder())
	bender.Record(r.Recorder(), r.Recorders()...)

	return nil
}
