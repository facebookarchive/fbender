/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package tester_test

import (
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/facebookincubator/fbender/metric"
	"github.com/facebookincubator/fbender/tester"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type MockedMetric struct {
	mock.Mock
}

func (m *MockedMetric) Setup(options interface{}) error {
	args := m.Called(options)

	return args.Error(0)
}

func (m *MockedMetric) Fetch(start time.Time, duration time.Duration) ([]tester.DataPoint, error) {
	args := m.Called(start, duration)

	return args.Get(0).([]tester.DataPoint), args.Error(1)
}

func (m *MockedMetric) Name() string {
	args := m.Called()

	return args.String(0)
}

type ParseConstraintTestSuite struct {
	suite.Suite

	metric  *MockedMetric
	parsers []tester.MetricParser
}

func (s *ParseConstraintTestSuite) SetupTest() {
	s.metric = new(MockedMetric)
	s.parsers = []tester.MetricParser{
		metric.Parser,
	}
}

//nolint:funlen
func (s *ParseConstraintTestSuite) TestConstructor() {
	// Fork bomb is not a valid constraint.
	c, err := tester.ParseConstraint("💣(){ 💣|💣& };💣", s.parsers...)
	s.Assert().Nil(c)
	s.Assert().Error(err)
	// Invalid aggregator 'PWN'.
	c, err = tester.ParseConstraint("PWN(time) < 24", s.parsers...)
	s.Assert().Nil(c)
	s.Assert().Equal(err, tester.ErrInvalidAggregator)
	// Valid constraint - AVG.
	c, err = tester.ParseConstraint("AVG(errors) < 10", s.parsers...)
	s.Assert().NotNil(c)
	s.Assert().NoError(err)
	// Valid constraint - MIN.
	c, err = tester.ParseConstraint("MIN(errors) < 5", s.parsers...)
	s.Assert().NotNil(c)
	s.Assert().NoError(err)
	// Valid constraint - MAX.
	c, err = tester.ParseConstraint("MAX(errors) < 20", s.parsers...)
	s.Assert().NotNil(c)
	s.Assert().NoError(err)
	// Invalid metric '0xdeadbeef'.
	c, err = tester.ParseConstraint("MAX(0xdeadbeef) < 10", s.parsers...)
	s.Assert().Nil(c)
	s.Assert().Equal(err, tester.ErrNotParsed)
	// Should return error if metric parser returned error.
	c, err = tester.ParseConstraint("MAX(mymetrics) > 1", func(_ string) (tester.Metric, error) {
		return nil, errors.New("error")
	})
	s.Assert().Nil(c)
	s.Assert().Error(err)
	// Valid metric - errors.
	c, err = tester.ParseConstraint("MAX(errors) < 10", s.parsers...)
	s.Assert().NotNil(c)
	s.Assert().NoError(err)
	// Valid metric - latency.
	c, err = tester.ParseConstraint("MAX(latency) < 100", s.parsers...)
	s.Assert().NotNil(c)
	s.Assert().NoError(err)
	// Invalid comparator '@@'.
	c, err = tester.ParseConstraint("MAX(errors) @@ 10", s.parsers...)
	s.Assert().Nil(c)
	s.Assert().Equal(err, tester.ErrInvalidComparator)
	// Valid threshold.
	c, err = tester.ParseConstraint("MIN(latency) < 1.337", s.parsers...)
	s.Assert().NotNil(c)
	s.Assert().NoError(err)
	// Invalid threshold - comma.
	c, err = tester.ParseConstraint("MAX(latency) < 1,337", s.parsers...)
	s.Assert().Nil(c)
	s.Assert().Error(err)
	// Invalid threshold - hex.
	c, err = tester.ParseConstraint("MAX(latency) < 0x414141", s.parsers...)
	s.Assert().Nil(c)
	s.Assert().Error(err)
	// Invalid threshold - overflow.
	number := `179769313486231570814527423731704356798070567525844996598917476803157260780028538760589558632766878171540
458953514382464234321326889464182768467546703537516986049910576551282076245490090389328944075868508455133942304583236903
2229481658085593321233482747978262041447231687381771809192998812504040261841248583681337`
	c, err = tester.ParseConstraint(
		fmt.Sprintf(
			"MAX(latency) < %s",
			strings.ReplaceAll(number, "\n", ""),
		),
		s.parsers...,
	)
	s.Assert().Nil(c)
	s.Require().Error(err)
	s.Assert().Contains(err.Error(), "value out of range")
}

func (s *ParseConstraintTestSuite) TestCheck() {
	now := time.Now()
	// No datapoints should result in error.
	s.metric.On("Fetch", now, time.Second).
		Return([]tester.DataPoint(nil), nil).Once()

	c := &tester.Constraint{Metric: s.metric}
	err := c.Check(now, time.Second)

	s.Assert().Error(err)
	s.metric.AssertExpectations(s.T())
	// If metric.Fetch returned error, Check should result in error.
	s.metric.On("Fetch", now, time.Second).
		Return([]tester.DataPoint(nil), errors.New("error")).Once()

	c = &tester.Constraint{Metric: s.metric}
	err = c.Check(now, time.Second)

	s.Assert().Error(err)
	s.metric.AssertExpectations(s.T())
	// Should pass - (10 < 10 + 1).
	s.metric.On("Fetch", now, time.Second).
		Return(
			[]tester.DataPoint{{Value: 10}},
			nil,
		).Once()

	c = &tester.Constraint{
		Metric:     s.metric,
		Aggregator: tester.MinimumAggregator,
		Comparator: tester.LessThan,
		Threshold:  10 + 1,
	}
	err = c.Check(now, time.Second)

	s.Assert().NoError(err)
	s.metric.AssertExpectations(s.T())
	// Should not pass - (10 > 10 + 1).
	s.metric.On("Fetch", now, time.Second).
		Return(
			[]tester.DataPoint{
				{Value: 10},
			},
			nil).Once()

	c = &tester.Constraint{
		Metric:     s.metric,
		Aggregator: tester.MinimumAggregator,
		Comparator: tester.GreaterThan,
		Threshold:  10 + 1,
	}
	err = c.Check(now, time.Second)

	s.Assert().Error(err)
	s.metric.AssertExpectations(s.T())
}

func TestParseConstraintTestSuite(t *testing.T) {
	suite.Run(t, new(ParseConstraintTestSuite))
}
