/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package tester

import (
	"errors"

	"github.com/pinterest/bender"
)

// ErrInvalidOptions is thrown when options don't implement the required interface.
var ErrInvalidOptions = errors.New("invalid options")

// Tester is used to setup the test for a specific endpoint.
type Tester interface {
	// Before is called once, before any tests.
	Before(options interface{}) error
	// After is called once, after all tests (or after some of them if a test fails).
	// This should be used to cleanup everything that was set up in the Before.
	After(options interface{})
	// BeforeEach is called before every test.
	BeforeEach(options interface{}) error
	// AfterEach is called after every test, even if the test fails. This should
	// be used to cleanup everything that was set up in the BeforeEach.
	AfterEach(options interface{})
	// RequestExecutor is called every time a test is to be ran to get an executor.
	RequestExecutor(options interface{}) (bender.RequestExecutor, error)
}

// QPS is the test desired queries per second.
type QPS = int

// ThroughputRunner is used to setup the test execution.
type ThroughputRunner interface {
	// Before is called before running a test.
	Before(qps QPS, options interface{}) error
	// After is called after test finishes. This should be used to clean up
	// everything that was ser up in the Before.
	After(qps QPS, options interface{})

	// Protocol tester.
	Tester() Tester

	// Params used by LoadTestThroughput function.
	Intervals() bender.IntervalGenerator
	Requests() chan interface{}
	Recorder() chan interface{}
	Recorders() []bender.Recorder
}

// Workers is the test desired concurrent workers.
type Workers = int

// ConcurrencyRunner is used to setup concurrency test execution.
type ConcurrencyRunner interface {
	// Before is called before running a test.
	Before(workers Workers, options interface{}) error
	// After is called after test finishes. This should be used to clean up
	// everything that was ser up in the Before.
	After(workers Workers, options interface{})

	// Protocol tester.
	Tester() Tester

	// Params used by LoadTestConcurrency function.
	WorkerSemaphore() *bender.WorkerSemaphore
	Requests() chan interface{}
	Recorder() chan interface{}
	Recorders() []bender.Recorder
}
