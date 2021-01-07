/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package flags_test

import (
	"testing"
	"time"

	"github.com/facebookincubator/fbender/flags"
	"github.com/facebookincubator/fbender/tester"
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
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

type MockedMetricParsers struct {
	mock.Mock
}

func (m *MockedMetricParsers) ParserA(name string) (tester.Metric, error) {
	args := m.Called(name)
	constraint := args.Get(0)
	err := args.Error(1)

	if constraint, ok := constraint.(tester.Metric); ok {
		//nolint:wrapcheck
		return constraint, err
	}

	//nolint:wrapcheck
	return nil, err
}

func (m *MockedMetricParsers) ParserB(name string) (tester.Metric, error) {
	args := m.Called(name)
	constraint := args.Get(0)
	err := args.Error(1)

	if constraint, ok := constraint.(tester.Metric); ok {
		//nolint:wrapcheck
		return constraint, err
	}

	//nolint:wrapcheck
	return nil, err
}

func (m *MockedMetricParsers) ParserC(name string) (tester.Metric, error) {
	args := m.Called(name)
	constraint := args.Get(0)
	err := args.Error(1)

	if constraint, ok := constraint.(tester.Metric); ok {
		//nolint:wrapcheck
		return constraint, err
	}

	//nolint:wrapcheck
	return nil, err
}

type ConstraintsSliceValueTestSuite struct {
	suite.Suite

	parsers *MockedMetricParsers
	value   pflag.Value
	ms      []*MockedMetric
	cs      []*tester.Constraint
}

func (s *ConstraintsSliceValueTestSuite) SetupTest() {
	s.parsers = new(MockedMetricParsers)
	s.value = flags.NewConstraintSliceValue(s.parsers.ParserA, s.parsers.ParserB, s.parsers.ParserC)
	s.Require().NotNil(s.value)

	s.ms = make([]*MockedMetric, 3)
	s.cs = make([]*tester.Constraint, 3)

	for i := 0; i < 3; i++ {
		s.ms[i] = new(MockedMetric)
		s.cs[i] = &tester.Constraint{
			Metric:     s.ms[i],
			Aggregator: tester.MaximumAggregator,
			Comparator: tester.LessThan,
			Threshold:  10. * float64(i),
		}
	}
}

func TestNewConstraintSliceValue(t *testing.T) {
	value := flags.NewConstraintSliceValue()
	assert.NotNil(t, value)
}

func (s *ConstraintsSliceValueTestSuite) TestType() {
	s.Assert().Equal("constraints", s.value.Type())
}

func (s *ConstraintsSliceValueTestSuite) TestSet_ParsersCalledInOrder() {
	// ParserA doesn't parse, returns ErrNotParsed.
	s.parsers.On("ParserA", "metric_0").Return(nil, tester.ErrNotParsed).Once()
	// ParserB parses and returns a constraint.
	s.parsers.On("ParserB", "metric_0").Return(s.ms[0], nil).Once()

	// Run the test
	err := s.value.Set("MAX(metric_0) < 0")
	s.Require().NoError(err)
	s.parsers.AssertExpectations(s.T())

	// Check if constraints were parsed
	cs, err := flags.GetConstraintsValue(s.value)
	s.Require().NoError(err)
	s.Assert().ElementsMatch(cs, s.cs[:1])
}

func (s *ConstraintsSliceValueTestSuite) TestSet_ParsersCalledForEveryValue() {
	// ParserA doesn't parse any of the metrics, returns ErrNotParsed.
	s.parsers.On("ParserA", "metric_0").Return(nil, tester.ErrNotParsed).Once()
	s.parsers.On("ParserA", "metric_1").Return(nil, tester.ErrNotParsed).Once()
	// ParserB parses metric_0.
	s.parsers.On("ParserB", "metric_0").Return(s.ms[0], nil).Once()
	s.parsers.On("ParserB", "metric_1").Return(nil, tester.ErrNotParsed).Once()
	// ParserC parses metric_1.
	s.parsers.On("ParserC", "metric_1").Return(s.ms[1], nil).Once()

	// Run the test.
	err := s.value.Set("MAX(metric_0) < 0, MAX(metric_1) < 10.0")
	s.Require().NoError(err)
	s.parsers.AssertExpectations(s.T())

	// Check if constraints were parsed.
	cs, err := flags.GetConstraintsValue(s.value)
	s.Require().NoError(err)
	s.Assert().ElementsMatch(cs, s.cs[:2])
}

func (s *ConstraintsSliceValueTestSuite) TestSet_ErrorsOnError() {
	s.parsers.On("ParserA", "metric_0").Return(nil, tester.ErrNotParsed).Once()
	s.parsers.On("ParserB", "metric_0").Return(nil, assert.AnError).Once()

	err := s.value.Set("MIN(metric_0) > 42.0")
	s.Assert().ErrorIs(err, assert.AnError)
	s.parsers.AssertExpectations(s.T())
}

func (s *ConstraintsSliceValueTestSuite) TestSet_ErrorsOnNotParsed() {
	s.parsers.On("ParserA", "metric_0").Return(nil, tester.ErrNotParsed).Once()
	s.parsers.On("ParserB", "metric_0").Return(nil, tester.ErrNotParsed).Once()
	s.parsers.On("ParserC", "metric_0").Return(nil, tester.ErrNotParsed).Once()

	err := s.value.Set("MAX(metric_0) < 0.0")
	s.Assert().ErrorIs(err, tester.ErrNotParsed)
	s.parsers.AssertExpectations(s.T())
}

func (s *ConstraintsSliceValueTestSuite) TestSet_AppendsOnConsecutiveCalls() {
	s.parsers.On("ParserA", "metric_0").Return(s.ms[0], nil).Once()
	s.parsers.On("ParserA", "metric_1").Return(s.ms[1], nil).Once()
	s.parsers.On("ParserA", "metric_2").Return(s.ms[2], nil).Once()

	// Run the test.
	err := s.value.Set("MAX(metric_0) < 0.0")
	s.Require().NoError(err)
	err = s.value.Set("MAX(metric_1) < 10.0, MAX(metric_2) < 20.0")
	s.Require().NoError(err)
	s.parsers.AssertExpectations(s.T())

	// Check if constraints were parsed.
	cs, err := flags.GetConstraintsValue(s.value)
	s.Require().NoError(err)
	s.Assert().ElementsMatch(cs, s.cs[:3])
}

func (s *ConstraintsSliceValueTestSuite) TestString() {
	// No constraints.
	s.Assert().Equal("[]", s.value.String())

	// Single constraint.
	s.parsers.On("ParserA", "metric_0").Return(s.ms[0], nil).Once()
	s.ms[0].On("Name").Return("metric_0").Once()

	err := s.value.Set("MAX(metric_0) < 0.0")
	s.Require().NoError(err)
	s.parsers.AssertExpectations(s.T())

	v := s.value.String()

	s.ms[0].AssertExpectations(s.T())
	s.Assert().Equal("[MAX(metric_0) < 0.00]", v)

	// Multiple constraints.
	s.parsers.On("ParserA", "metric_1").Return(s.ms[1], nil).Once()
	s.ms[0].On("Name").Return("metric_0").Once()
	s.ms[1].On("Name").Return("metric_1").Once()

	err = s.value.Set("MAX(metric_1) < 10.0")
	s.Require().NoError(err)
	s.parsers.AssertExpectations(s.T())

	v = s.value.String()

	s.ms[0].AssertExpectations(s.T())
	s.ms[1].AssertExpectations(s.T())
	s.Assert().Equal("[MAX(metric_0) < 0.00 MAX(metric_1) < 10.00]", v)
}

func (s *ConstraintsSliceValueTestSuite) TestGetConstraints() {
	f := pflag.NewFlagSet("Test FlagSet", pflag.ExitOnError)
	f.Var(s.value, "constraints", "set constraints")

	s.parsers.On("ParserA", "metric_0").Return(s.ms[0], nil).Once()
	err := s.value.Set("MAX(metric_0) < 0.0")
	s.Require().NoError(err)
	s.parsers.AssertExpectations(s.T())

	cs, err := flags.GetConstraints(f, "constraints")
	s.Require().NoError(err)
	s.Assert().ElementsMatch(cs, s.cs[:1])

	// Check error when flag does not exist
	_, err = flags.GetConstraints(f, "nonexistent")
	s.Assert().ErrorIs(err, flags.ErrUndefined)
	s.Assert().EqualError(err, "flag accessed but not defined: \"nonexistent\"")

	// Check error when value is of different type
	f.Int("myint", 0, "set myint")
	_, err = flags.GetConstraints(f, "myint")
	s.Assert().ErrorIs(err, flags.ErrInvalidType)
	s.Assert().EqualError(err, "accessed flag type does not match, want: constraints, got: int")
}

func (s *ConstraintsSliceValueTestSuite) TestGetConstraintsValue() {
	s.parsers.On("ParserA", "metric_0").Return(s.ms[0], nil).Once()
	err := s.value.Set("MAX(metric_0) < 0.0")
	s.Require().NoError(err)
	s.parsers.AssertExpectations(s.T())

	cs, err := flags.GetConstraintsValue(s.value)
	s.Require().NoError(err)
	s.Assert().ElementsMatch(cs, s.cs[:1])

	// Check error when value is of different type
	f := pflag.NewFlagSet("Test FlagSet", pflag.ExitOnError)
	f.Int("myint", 0, "set myint")
	flag := f.Lookup("myint")
	s.Require().NotNil(flag)
	_, err = flags.GetConstraintsValue(flag.Value)
	s.Assert().ErrorIs(err, flags.ErrInvalidType)
	s.Assert().EqualError(err, "accessed flag type does not match, want: constraints, got: int")
}

func TestConstraintsSliceValueTestSuite(t *testing.T) {
	suite.Run(t, new(ConstraintsSliceValueTestSuite))
}
