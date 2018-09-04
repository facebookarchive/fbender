/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package runner

import (
	"context"
	"time"

	"github.com/pinterest/bender"

	"github.com/facebookincubator/fbender/cmd/core/options"
	"github.com/facebookincubator/fbender/recorders"
	"github.com/facebookincubator/fbender/tester"
	"github.com/facebookincubator/fbender/utils"
)

// ConcurrencyRunner is a test runner for load test concurrency commands
type ConcurrencyRunner struct {
	runner
	workerSem     *bender.WorkerSemaphore
	spinnerCancel context.CancelFunc
}

// NewConcurrencyRunner returns new ConcurrencyRunner
func NewConcurrencyRunner(params *Params) *ConcurrencyRunner {
	return &ConcurrencyRunner{
		runner: runner{
			Params: params,
		},
	}
}

// Before prepares requests, recorders and interval generator
func (r *ConcurrencyRunner) Before(workers tester.Workers, opts interface{}) error {
	if err := r.runner.Before(workers, opts); err != nil {
		return err
	}
	o, ok := opts.(*options.Options)
	if !ok {
		return tester.ErrInvalidOptions
	}

	r.workerSem = bender.NewWorkerSemaphore()
	go func() { r.workerSem.Signal(workers) }()

	r.requests = make(chan interface{})
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		for i := 0; ; i++ {
			select {
			case <-ctx.Done():
				close(r.requests)
				return
			default:
				r.requests <- r.Params.RequestGenerator(i)
			}
		}
	}()

	// We want bar to measure time
	count := int(o.Duration/time.Second) * 10
	r.progress, r.bar = recorders.NewLoadTestProgress(count)
	r.progress.Start()
	go func() {
		for i := 0; i < count; i++ {
			time.Sleep(time.Second / 10)
			r.bar.Incr()
		}
		cancel()
		r.progress.Stop()
		r.spinnerCancel = utils.NewBackgroundSpinner("Waiting for tests to finish", 0)
	}()
	return nil
}

// After cleans up after the test
func (r *ConcurrencyRunner) After(test int, options interface{}) {
	r.spinnerCancel()
	r.runner.After(test, options)
}

// WorkerSemaphore returns a worker semaphore for concurrency test
func (r *ConcurrencyRunner) WorkerSemaphore() *bender.WorkerSemaphore {
	return r.workerSem
}
