/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package run_test

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/facebookincubator/fbender/log"
	"github.com/facebookincubator/fbender/tester"
	"github.com/facebookincubator/fbender/tester/run"
	"github.com/pinterest/bender"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type MockedThroughputRunner struct {
	mock.Mock
}

func (m *MockedThroughputRunner) Before(qps tester.QPS, options interface{}) error {
	args := m.Called(qps, options)

	return args.Error(0)
}

func (m *MockedThroughputRunner) After(qps tester.QPS, options interface{}) {
	m.Called(qps, options)
}

func (m *MockedThroughputRunner) Tester() tester.Tester {
	args := m.Called()
	if t, ok := args.Get(0).(tester.Tester); ok {
		return t
	}

	panic(fmt.Sprintf("assert: arguments: Tester(0) failed because object wasn't correct type: %v", args.Get(0)))
}

func (m *MockedThroughputRunner) Intervals() bender.IntervalGenerator {
	args := m.Called()
	if intervalGenerator, ok := args.Get(0).(bender.IntervalGenerator); ok {
		return intervalGenerator
	}

	panic(fmt.Sprintf("assert: arguments: Intervals(0) failed because object wasn't correct type: %v", args.Get(0)))
}

func (m *MockedThroughputRunner) Requests() chan interface{} {
	args := m.Called()
	if channel, ok := args.Get(0).(chan interface{}); ok {
		return channel
	}

	panic(fmt.Sprintf("assert: arguments: Channel(0) failed because object wasn't correct type: %v", args.Get(0)))
}

func (m *MockedThroughputRunner) Recorder() chan interface{} {
	args := m.Called()
	if channel, ok := args.Get(0).(chan interface{}); ok {
		return channel
	}

	panic(fmt.Sprintf("assert: arguments: Channel(0) failed because object wasn't correct type: %v", args.Get(0)))
}

func (m *MockedThroughputRunner) Recorders() []bender.Recorder {
	args := m.Called()
	if recorders, ok := args.Get(0).([]bender.Recorder); ok {
		return recorders
	}

	panic(fmt.Sprintf("assert: arguments: RecorderSlice(0) failed because object wasn't correct type: %v", args.Get(0)))
}

type ThroughputFixedTestSuite struct {
	suite.Suite
	tester  *MockedTester
	runner  *MockedThroughputRunner
	options interface{}
}

// dummyRequests generates a buffered channel and fills it with requests
// with consecutive integers and mocks calls on executor for them to return
// a request and a predefined error.
func (s *ThroughputFixedTestSuite) dummyRequests(n int, err error) chan interface{} {
	requests := make(chan interface{}, n)
	for i := 0; i < n; i++ {
		requests <- i
		s.tester.On("DummyExecutor", mock.Anything, i).Return(i, err).Once()
	}

	close(requests)

	return requests
}

func (s *ThroughputFixedTestSuite) SetupTest() {
	s.tester = new(MockedTester)
	s.runner = new(MockedThroughputRunner)
	s.options = new(struct{})
}

func (s *ThroughputFixedTestSuite) TestTester__Before_Error() {
	s.runner.On("Tester").Return(s.tester).Once()
	s.tester.On("Before", s.options).Return(ErrDummy).Once()

	err := run.LoadTestThroughputFixed(s.runner, s.options, 10)
	s.Assert().Equal(ErrDummy, err)

	s.tester.AssertExpectations(s.T())
	s.runner.AssertExpectations(s.T())
}

func (s *ThroughputFixedTestSuite) TestTester__BeforeEach_Error() {
	// Single test
	s.runner.On("Tester").Return(s.tester).Once()
	s.tester.On("Before", s.options).Return(nil).Once()
	s.tester.On("After", s.options).Once()
	s.tester.On("BeforeEach", s.options).Return(ErrDummy).Once()

	err := run.LoadTestThroughputFixed(s.runner, s.options, 10)
	s.Assert().Equal(ErrDummy, err)

	s.tester.AssertExpectations(s.T())
	s.runner.AssertExpectations(s.T())

	// Multiple tests - Make sure all tests stop after first failed test
	s.runner.On("Tester").Return(s.tester).Once()
	s.tester.On("Before", s.options).Return(nil).Once()
	s.tester.On("After", s.options).Once()
	s.tester.On("BeforeEach", s.options).Return(ErrDummy).Once()

	err = run.LoadTestThroughputFixed(s.runner, s.options, 10, 20)
	s.Assert().Equal(ErrDummy, err)

	s.tester.AssertExpectations(s.T())
	s.runner.AssertExpectations(s.T())
}

func (s *ThroughputFixedTestSuite) TestRunner__Before_Error() {
	// Single test
	s.runner.On("Tester").Return(s.tester).Once()
	s.tester.On("Before", s.options).Return(nil).Once()
	s.tester.On("After", s.options).Once()
	s.tester.On("BeforeEach", s.options).Return(nil).Once()
	s.tester.On("AfterEach", s.options).Once()
	s.runner.On("Before", 10, s.options).Return(ErrDummy).Once()

	err := run.LoadTestThroughputFixed(s.runner, s.options, 10)
	s.Assert().Equal(ErrDummy, err)

	s.tester.AssertExpectations(s.T())
	s.runner.AssertExpectations(s.T())

	// Multiple tests - Make sure all tests stop after first failed test
	s.runner.On("Tester").Return(s.tester).Once()
	s.tester.On("Before", s.options).Return(nil).Once()
	s.tester.On("After", s.options).Once()
	s.tester.On("BeforeEach", s.options).Return(nil).Once()
	s.tester.On("AfterEach", s.options).Once()
	s.runner.On("Before", 10, s.options).Return(ErrDummy).Once()

	err = run.LoadTestThroughputFixed(s.runner, s.options, 10, 20)
	s.Assert().Equal(ErrDummy, err)

	s.tester.AssertExpectations(s.T())
	s.runner.AssertExpectations(s.T())
}

func (s *ThroughputFixedTestSuite) TestTester__RequestExecutor_Error() {
	// Single test
	s.runner.On("Tester").Return(s.tester).Once()
	s.tester.On("Before", s.options).Return(nil).Once()
	s.tester.On("After", s.options).Once()
	s.tester.On("BeforeEach", s.options).Return(nil).Once()
	s.tester.On("AfterEach", s.options).Once()
	s.runner.On("Before", 10, s.options).Return(nil).Once()
	s.runner.On("After", 10, s.options).Once()
	s.tester.On("RequestExecutor", s.options).Return(ErrDummy).Once()

	err := run.LoadTestThroughputFixed(s.runner, s.options, 10)
	s.Assert().Equal(ErrDummy, err)

	s.tester.AssertExpectations(s.T())
	s.runner.AssertExpectations(s.T())

	// Multiple tests - Make sure all tests stop after first failed test
	s.runner.On("Tester").Return(s.tester).Once()
	s.tester.On("Before", s.options).Return(nil).Once()
	s.tester.On("After", s.options).Once()
	s.tester.On("BeforeEach", s.options).Return(nil).Once()
	s.tester.On("AfterEach", s.options).Once()
	s.runner.On("Before", 10, s.options).Return(nil).Once()
	s.runner.On("After", 10, s.options).Once()
	s.tester.On("RequestExecutor", s.options).Return(ErrDummy).Once()

	err = run.LoadTestThroughputFixed(s.runner, s.options, 10, 20)
	s.Assert().Equal(ErrDummy, err)

	s.tester.AssertExpectations(s.T())
	s.runner.AssertExpectations(s.T())
}

func (s *ThroughputFixedTestSuite) TestZero() {
	s.runner.On("Tester").Return(s.tester).Once()
	s.tester.On("Before", s.options).Return(nil).Once()
	s.tester.On("After", s.options).Once()

	err := run.LoadTestThroughputFixed(s.runner, s.options)
	s.Assert().NoError(err)

	s.tester.AssertExpectations(s.T())
	s.runner.AssertExpectations(s.T())
}

func (s *ThroughputFixedTestSuite) TestSingle() {
	// Single test no failures
	s.runner.On("Tester").Return(s.tester).Once()
	s.tester.On("Before", s.options).Return(nil).Once()
	s.tester.On("After", s.options).Once()
	s.tester.On("BeforeEach", s.options).Return(nil).Once()
	s.tester.On("AfterEach", s.options).Once()
	s.runner.On("Before", 10, s.options).Return(nil).Once()
	s.runner.On("After", 10, s.options).Once()
	s.tester.On("RequestExecutor", s.options).Return(nil).Once()

	requests := s.dummyRequests(10, nil)
	s.runner.On("Requests").Return(requests).Once()
	s.runner.On("Intervals").Return(bender.UniformIntervalGenerator(100)).Once()

	recorder := make(chan interface{}, 10)
	s.runner.On("Recorder").Return(recorder).Twice()
	s.runner.On("Recorders").Return([]bender.Recorder{}).Once()

	err := run.LoadTestThroughputFixed(s.runner, s.options, 10)
	s.Assert().NoError(err)

	s.tester.AssertExpectations(s.T())
	s.runner.AssertExpectations(s.T())

	// Requests can fail and it shouldn't stop the test
	s.runner.On("Tester").Return(s.tester).Once()
	s.tester.On("Before", s.options).Return(nil).Once()
	s.tester.On("After", s.options).Once()
	s.tester.On("BeforeEach", s.options).Return(nil).Once()
	s.tester.On("AfterEach", s.options).Once()
	s.runner.On("Before", 10, s.options).Return(nil).Once()
	s.runner.On("After", 10, s.options).Once()
	s.tester.On("RequestExecutor", s.options).Return(nil).Once()

	requests = s.dummyRequests(20, ErrDummy)
	s.runner.On("Requests").Return(requests).Once()
	s.runner.On("Intervals").Return(bender.UniformIntervalGenerator(100)).Once()

	recorder = make(chan interface{}, 20)
	s.runner.On("Recorder").Return(recorder).Twice()
	s.runner.On("Recorders").Return([]bender.Recorder{}).Once()

	err = run.LoadTestThroughputFixed(s.runner, s.options, 10)
	s.Assert().NoError(err)

	s.tester.AssertExpectations(s.T())
	s.runner.AssertExpectations(s.T())
}

func (s *ThroughputFixedTestSuite) TestMultiple() {
	// Make sure Before/After gets called only once
	s.runner.On("Tester").Return(s.tester).Once()
	s.tester.On("Before", s.options).Return(nil).Once()
	s.tester.On("After", s.options).Once()

	// Make sure Before/After Each gets executed before each test
	s.tester.On("BeforeEach", s.options).Return(nil).Twice()
	s.tester.On("AfterEach", s.options).Twice()

	// Make sure Before/After run once for every test
	s.runner.On("Before", 10, s.options).Return(nil).Once()
	s.runner.On("After", 10, s.options).Once()
	s.runner.On("Before", 20, s.options).Return(nil).Once()
	s.runner.On("After", 20, s.options).Once()

	// Parameters that are always the same
	s.tester.On("RequestExecutor", s.options).Return(nil).Twice()
	s.runner.On("Intervals").Return(bender.UniformIntervalGenerator(100)).Twice()
	s.runner.On("Recorders").Return([]bender.Recorder{}).Twice()

	// Generate parameters for first test
	requests := s.dummyRequests(10, nil)
	s.runner.On("Requests").Return(requests).Once()

	recorder := make(chan interface{}, 10)
	s.runner.On("Recorder").Return(recorder).Twice()

	// Generate parameters for second test
	requests = s.dummyRequests(20, nil)
	s.runner.On("Requests").Return(requests).Once()

	recorder = make(chan interface{}, 20)
	s.runner.On("Recorder").Return(recorder).Twice()

	// Run tests
	err := run.LoadTestThroughputFixed(s.runner, s.options, 10, 20)
	s.Assert().NoError(err)

	s.tester.AssertExpectations(s.T())
	s.runner.AssertExpectations(s.T())
}

func TestThroughputFixedTestSuite(t *testing.T) {
	suite.Run(t, new(ThroughputFixedTestSuite))
}

type ThroughputConstraintsTestSuite struct {
	suite.Suite
	tester  *MockedTester
	runner  *MockedThroughputRunner
	growth  *MockedGrowth
	options interface{}
}

// dummyRequests generates a buffered channel and fills it with requests
// with consecutive integers and mocks calls on executor for them to return
// a request and a predefined error.
func (s *ThroughputConstraintsTestSuite) dummyRequests(n int, err error) chan interface{} {
	requests := make(chan interface{}, n)
	for i := 0; i < n; i++ {
		requests <- i
		s.tester.On("DummyExecutor", mock.Anything, i).Return(i, err).Once()
	}

	close(requests)

	return requests
}

func (s *ThroughputConstraintsTestSuite) SetupTest() {
	log.Stderr = ioutil.Discard
	s.tester = new(MockedTester)
	s.runner = new(MockedThroughputRunner)
	s.growth = new(MockedGrowth)
	s.options = new(struct{})
}

func (s *ThroughputConstraintsTestSuite) TestTester__Before_Error() {
	s.runner.On("Tester").Return(s.tester).Once()
	s.tester.On("Before", s.options).Return(ErrDummy).Once()

	err := run.LoadTestThroughputConstraints(s.runner, s.options, 10, s.growth)
	s.Assert().Equal(ErrDummy, err)

	s.tester.AssertExpectations(s.T())
	s.runner.AssertExpectations(s.T())
	s.growth.AssertExpectations(s.T())
}

func (s *ThroughputConstraintsTestSuite) TestTester__BeforeEach_Error() {
	s.runner.On("Tester").Return(s.tester).Once()
	s.tester.On("Before", s.options).Return(nil).Once()
	s.tester.On("After", s.options).Once()
	s.tester.On("BeforeEach", s.options).Return(ErrDummy).Once()

	err := run.LoadTestThroughputConstraints(s.runner, s.options, 10, s.growth)
	s.Assert().Equal(ErrDummy, err)

	s.tester.AssertExpectations(s.T())
	s.runner.AssertExpectations(s.T())
	s.growth.AssertExpectations(s.T())
}

func (s *ThroughputConstraintsTestSuite) TestRunner__Before_Error() {
	s.runner.On("Tester").Return(s.tester).Once()
	s.tester.On("Before", s.options).Return(nil).Once()
	s.tester.On("After", s.options).Once()
	s.tester.On("BeforeEach", s.options).Return(nil).Once()
	s.tester.On("AfterEach", s.options).Once()
	s.runner.On("Before", 10, s.options).Return(ErrDummy).Once()

	err := run.LoadTestThroughputConstraints(s.runner, s.options, 10, s.growth)
	s.Assert().Equal(ErrDummy, err)

	s.tester.AssertExpectations(s.T())
	s.runner.AssertExpectations(s.T())
	s.growth.AssertExpectations(s.T())
}

func (s *ThroughputConstraintsTestSuite) TestTester__RequestExecutor_Error() {
	s.runner.On("Tester").Return(s.tester).Once()
	s.tester.On("Before", s.options).Return(nil).Once()
	s.tester.On("After", s.options).Once()
	s.tester.On("BeforeEach", s.options).Return(nil).Once()
	s.tester.On("AfterEach", s.options).Once()
	s.runner.On("Before", 10, s.options).Return(nil).Once()
	s.runner.On("After", 10, s.options).Once()
	s.tester.On("RequestExecutor", s.options).Return(ErrDummy).Once()

	err := run.LoadTestThroughputConstraints(s.runner, s.options, 10, s.growth)
	s.Assert().Equal(ErrDummy, err)

	s.tester.AssertExpectations(s.T())
	s.runner.AssertExpectations(s.T())
	s.growth.AssertExpectations(s.T())
}

func (s *ThroughputConstraintsTestSuite) TestSingle_OnSuccess() {
	s.runner.On("Tester").Return(s.tester).Once()
	s.tester.On("Before", s.options).Return(nil).Once()
	s.tester.On("After", s.options).Once()
	s.tester.On("BeforeEach", s.options).Return(nil).Once()
	s.tester.On("AfterEach", s.options).Once()
	s.runner.On("Before", 10, s.options).Return(nil).Once()
	s.runner.On("After", 10, s.options).Once()
	s.tester.On("RequestExecutor", s.options).Return(nil).Once()

	requests := s.dummyRequests(10, nil)
	s.runner.On("Requests").Return(requests).Once()
	s.runner.On("Intervals").Return(bender.UniformIntervalGenerator(100)).Once()

	recorder := make(chan interface{}, 10)
	s.runner.On("Recorder").Return(recorder).Twice()
	s.runner.On("Recorders").Return([]bender.Recorder{}).Once()

	c := NewMockedConstraint(true)

	s.growth.On("OnSuccess", 10).Return(0).Once()

	err := run.LoadTestThroughputConstraints(s.runner, s.options, 10, s.growth, c.Constraint())
	s.Assert().NoError(err)

	s.tester.AssertExpectations(s.T())
	s.runner.AssertExpectations(s.T())
	s.growth.AssertExpectations(s.T())
	c.AssertExpectations(s.T())
}

func (s *ThroughputConstraintsTestSuite) TestSingle_OnFail() {
	s.runner.On("Tester").Return(s.tester).Once()
	s.tester.On("Before", s.options).Return(nil).Once()
	s.tester.On("After", s.options).Once()
	s.tester.On("BeforeEach", s.options).Return(nil).Once()
	s.tester.On("AfterEach", s.options).Once()
	s.runner.On("Before", 10, s.options).Return(nil).Once()
	s.runner.On("After", 10, s.options).Once()
	s.tester.On("RequestExecutor", s.options).Return(nil).Once()

	requests := s.dummyRequests(10, nil)
	s.runner.On("Requests").Return(requests).Once()
	s.runner.On("Intervals").Return(bender.UniformIntervalGenerator(100)).Once()

	recorder := make(chan interface{}, 10)
	s.runner.On("Recorder").Return(recorder).Twice()
	s.runner.On("Recorders").Return([]bender.Recorder{}).Once()

	c := NewMockedConstraint(false)

	s.growth.On("OnFail", 10).Return(0).Once()

	err := run.LoadTestThroughputConstraints(s.runner, s.options, 10, s.growth, c.Constraint())
	s.Assert().NoError(err)

	s.tester.AssertExpectations(s.T())
	s.runner.AssertExpectations(s.T())
	s.growth.AssertExpectations(s.T())
	c.AssertExpectations(s.T())
}

func (s *ThroughputConstraintsTestSuite) TestMultiple() {
	s.runner.On("Tester").Return(s.tester).Once()

	// Make sure Before/After gets called only once
	s.tester.On("Before", s.options).Return(nil).Once()
	s.tester.On("After", s.options).Once()

	// Make sure Before/After Each gets executed before each test
	s.tester.On("BeforeEach", s.options).Return(nil).Twice()
	s.tester.On("AfterEach", s.options).Twice()

	// Make sure Before/After run once for every test
	s.runner.On("Before", 10, s.options).Return(nil).Once()
	s.runner.On("After", 10, s.options).Once()
	s.runner.On("Before", 20, s.options).Return(nil).Once()
	s.runner.On("After", 20, s.options).Once()

	// Parameters that are always the same
	s.tester.On("RequestExecutor", s.options).Return(nil).Twice()
	s.runner.On("Intervals").Return(bender.UniformIntervalGenerator(100)).Twice()
	s.runner.On("Recorders").Return([]bender.Recorder{}).Twice()

	// Generate parameters for first test
	requests := s.dummyRequests(10, nil)
	s.runner.On("Requests").Return(requests).Once()

	recorder := make(chan interface{}, 10)
	s.runner.On("Recorder").Return(recorder).Twice()

	// Generate parameters for second test
	requests = s.dummyRequests(20, ErrDummy)
	s.runner.On("Requests").Return(requests).Once()

	recorder = make(chan interface{}, 20)
	s.runner.On("Recorder").Return(recorder).Twice()

	// Constraint
	c := NewMockedConstraint(true, false)

	s.growth.On("OnSuccess", 10).Return(20).Once()
	s.growth.On("OnFail", 20).Return(0).Once()

	err := run.LoadTestThroughputConstraints(s.runner, s.options, 10, s.growth, c.Constraint())
	s.Assert().NoError(err)

	s.tester.AssertExpectations(s.T())
	s.runner.AssertExpectations(s.T())
	s.growth.AssertExpectations(s.T())
	c.AssertExpectations(s.T())
}

func TestThroughputConstraintsTestSuite(t *testing.T) {
	suite.Run(t, new(ThroughputConstraintsTestSuite))
}
