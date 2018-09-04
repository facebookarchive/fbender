/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package run_test

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/pinterest/bender"
	"github.com/stretchr/testify/mock"

	"github.com/facebookincubator/fbender/tester"
)

type MockedTester struct {
	mock.Mock
}

func (m *MockedTester) Before(options interface{}) error {
	args := m.Called(options)
	return args.Error(0)
}

func (m *MockedTester) After(options interface{}) {
	m.Called(options)
}

func (m *MockedTester) BeforeEach(options interface{}) error {
	args := m.Called(options)
	return args.Error(0)
}

func (m *MockedTester) AfterEach(options interface{}) {
	m.Called(options)
}

func (m *MockedTester) RequestExecutor(options interface{}) (bender.RequestExecutor, error) {
	args := m.Called(options)
	return m.DummyExecutor, args.Error(0)
}

func (m *MockedTester) DummyExecutor(timestamp int64, request interface{}) (interface{}, error) {
	args := m.Called(timestamp, request)
	return args.Get(0), args.Error(1)
}

var ErrDummy = errors.New("dummy error")

type MockedGrowth struct {
	mock.Mock
}

func (m *MockedGrowth) String() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockedGrowth) OnSuccess(test int) int {
	args := m.Called(test)
	return args.Int(0)
}

func (m *MockedGrowth) OnFail(test int) int {
	args := m.Called(test)
	return args.Int(0)
}

type MockedMetric struct {
	mock.Mock
}

func (m *MockedMetric) Setup(options interface{}) error {
	args := m.Called(options)
	return args.Error(0)
}

func (m *MockedMetric) Fetch(start time.Time, duration time.Duration) ([]tester.DataPoint, error) {
	args := m.Called(start, duration)
	if points, ok := args.Get(0).([]tester.DataPoint); ok {
		return points, args.Error(1)
	}
	panic(fmt.Sprintf("assert: arguments: DataPointSlice(0) failed because object wasn't correct type: %v", args.Get(0)))
}

func (m *MockedMetric) Name() string {
	args := m.Called()
	return args.String(0)
}

type MockedAggregator struct {
	mock.Mock
}

func (m *MockedAggregator) Aggregate(points []tester.DataPoint) float64 {
	args := m.Called(points)
	return args.Get(0).(float64)
}

func (m *MockedAggregator) Name() string {
	args := m.Called()
	return args.String(0)
}

type MockedComparator struct {
	mock.Mock
}

func (m *MockedComparator) Compare(x, y float64) bool {
	args := m.Called(x, y)
	return args.Bool(0)
}

func (m *MockedComparator) Name() string {
	args := m.Called()
	return args.String(0)
}

type MockedConstraint struct {
	Metric     *MockedMetric
	Aggregator *MockedAggregator
	Comparator *MockedComparator
	Threshold  float64
}

// NewMockedConstraint returns a new mocked constraint with already mocked
// calls for a proper Constraint.Check function. Each call will return a result
// from the results list.
func NewMockedConstraint(results ...bool) *MockedConstraint {
	p := []tester.DataPoint{}
	n := len(results)
	c := &MockedConstraint{
		Metric:     new(MockedMetric),
		Aggregator: new(MockedAggregator),
		Comparator: new(MockedComparator),
		Threshold:  float64(100),
	}
	c.Metric.On("Fetch", mock.Anything, mock.Anything).Return(p, nil).Times(n)
	c.Aggregator.On("Aggregate", p).Return(float64(50)).Times(n)
	for _, result := range results {
		if !result {
			c.Metric.On("Name").Return("Metric").Once()
			c.Aggregator.On("Name").Return("Aggregator").Once()
			c.Comparator.On("Name").Return("?").Twice()
		}
		c.Comparator.On("Compare", float64(50), float64(100)).Return(result).Once()
	}
	return c
}

func (m *MockedConstraint) Constraint() *tester.Constraint {
	return &tester.Constraint{
		Metric:     m.Metric,
		Aggregator: m.Aggregator,
		Comparator: m.Comparator,
		Threshold:  m.Threshold,
	}
}

func (m *MockedConstraint) AssertExpectations(t *testing.T) {
	m.Metric.AssertExpectations(t)
	m.Aggregator.AssertExpectations(t)
	m.Comparator.AssertExpectations(t)
}
