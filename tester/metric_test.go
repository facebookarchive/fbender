/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package tester_test

import (
	"math/rand"
	"testing"
	"time"

	"github.com/facebookincubator/fbender/tester"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// randomValues generates a slice filled with random values datapoints.
func randomDataPoints(n int) []tester.DataPoint {
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)

	values := make([]tester.DataPoint, n)
	for i := 0; i < 10; i++ {
		values[i] = tester.DataPoint{Value: r.Float64()}
	}

	return values
}

type MinimumAggregatorTestSuite struct {
	suite.Suite
	aggregator tester.Aggregator
}

func (s *MinimumAggregatorTestSuite) SetupTest() {
	s.aggregator = tester.MinimumAggregator
}

func (s *MinimumAggregatorTestSuite) testAggregateRandomSlice(n int) {
	points := randomDataPoints(n)
	v := s.aggregator.Aggregate(points)
	min := points[0].Value

	for _, point := range points {
		s.Assert().True(v <= point.Value, "expected %f to be less or equal to %f", v, point.Value)

		if point.Value < min {
			min = point.Value
		}
	}

	s.Assert().Equal(min, v)
}

func (s *MinimumAggregatorTestSuite) TestName() {
	s.Assert().Equal("MIN", s.aggregator.Name())
}

func (s *MinimumAggregatorTestSuite) TestAggregate() {
	// For nil and empty list we should get 0.
	v := s.aggregator.Aggregate(nil)
	s.Assert().Equal(0., v)

	v = s.aggregator.Aggregate([]tester.DataPoint{})
	s.Assert().Equal(0., v)

	// Single value is the minimum
	v = s.aggregator.Aggregate([]tester.DataPoint{
		{Value: 42.},
	})
	s.Assert().Equal(42., v)

	// Random tests to make sure we get minimum
	for i := 0; i < 128; i++ {
		s.testAggregateRandomSlice(4096)
	}
}

type MaximumAggregatorTestSuite struct {
	suite.Suite
	aggregator tester.Aggregator
}

func (s *MaximumAggregatorTestSuite) SetupTest() {
	s.aggregator = tester.MaximumAggregator
}

func (s *MaximumAggregatorTestSuite) testAggregateRandomSlice(n int) {
	points := randomDataPoints(n)
	v := s.aggregator.Aggregate(points)
	max := points[0].Value

	for _, point := range points {
		s.Assert().True(v >= point.Value, "expected %f to be greater or equal to %f", v, point.Value)

		if point.Value > max {
			max = point.Value
		}
	}

	s.Assert().Equal(max, v)
}

func (s *MaximumAggregatorTestSuite) TestName() {
	s.Assert().Equal("MAX", s.aggregator.Name())
}

func (s *MaximumAggregatorTestSuite) TestAggregate() {
	// For nil and empty list we should get 0.
	v := s.aggregator.Aggregate(nil)
	s.Assert().Equal(0., v)

	v = s.aggregator.Aggregate([]tester.DataPoint{})
	s.Assert().Equal(0., v)

	// Single value is the maximum
	v = s.aggregator.Aggregate([]tester.DataPoint{
		{Value: 42.},
	})
	s.Assert().Equal(42., v)

	// Random tests to make sure we get maximum
	for i := 0; i < 128; i++ {
		s.testAggregateRandomSlice(4096)
	}
}

type AverageAggregatorTestSuite struct {
	suite.Suite
	aggregator tester.Aggregator
}

func (s *AverageAggregatorTestSuite) SetupTest() {
	s.aggregator = tester.AverageAggregator
}

func (s *AverageAggregatorTestSuite) testAggregateRandomSlice(n int) {
	points := randomDataPoints(n)
	v := s.aggregator.Aggregate(points)
	sum := 0.

	for _, point := range points {
		sum += point.Value
	}

	s.Assert().Equal(sum/float64(n), v)
}

func (s *AverageAggregatorTestSuite) TestName() {
	s.Assert().Equal("AVG", s.aggregator.Name())
}

func (s *AverageAggregatorTestSuite) TestAggregate() {
	// For nil and empty list we should get 0.
	v := s.aggregator.Aggregate(nil)
	s.Assert().Equal(0., v)

	v = s.aggregator.Aggregate([]tester.DataPoint{})
	s.Assert().Equal(0., v)

	// Single value is the average
	v = s.aggregator.Aggregate([]tester.DataPoint{
		{Value: 42.},
	})
	s.Assert().Equal(42., v)

	// Few manual tests
	points := []tester.DataPoint{
		{Value: 10.},
		{Value: 20.},
		{Value: 30.},
	}

	v = s.aggregator.Aggregate(points)
	s.Assert().Equal(20., v)

	points = []tester.DataPoint{
		{Value: 1.},
		{Value: 2.},
		{Value: 3.},
		{Value: 4.},
		{Value: 5.},
		{Value: 6.},
		{Value: 7.},
		{Value: 8.},
		{Value: 9.},
	}

	v = s.aggregator.Aggregate(points)
	s.Assert().Equal(5., v)

	// Random tests to make sure we get maximum
	for i := 0; i < 128; i++ {
		s.testAggregateRandomSlice(4096)
	}
}

func TestParseAggregator(t *testing.T) {
	a, err := tester.ParseAggregator("MIN")
	assert.NoError(t, err)
	assertPointerEqual(t, tester.MinimumAggregator, a, "Expected minimum aggregator")

	a, err = tester.ParseAggregator("MAX")
	assert.NoError(t, err)
	assertPointerEqual(t, tester.MaximumAggregator, a, "Expected maximum aggregator")

	a, err = tester.ParseAggregator("AVG")
	assert.NoError(t, err)
	assertPointerEqual(t, tester.AverageAggregator, a, "Expected average aggregator")

	a, err = tester.ParseAggregator("Nonexistent")
	assert.Error(t, err)
	assert.Nil(t, a)
	assert.Equal(t, tester.ErrInvalidAggregator, err)
}

func TestMinimumAggregatorTestSuite(t *testing.T) {
	suite.Run(t, new(MinimumAggregatorTestSuite))
}

func TestMaximumAggregatorTestSuite(t *testing.T) {
	suite.Run(t, new(MaximumAggregatorTestSuite))
}

func TestAverageAggregatorTestSuite(t *testing.T) {
	suite.Run(t, new(AverageAggregatorTestSuite))
}
