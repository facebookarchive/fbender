/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package tester

import (
	"errors"
	"time"
)

// DataPoint represents a sample of data.
type DataPoint struct {
	Time  time.Time
	Value float64
}

// Metric provides a function to get data points of this metric.
type Metric interface {
	// Setup is used to setup the metric before the test
	Setup(options interface{}) error
	// Fetch is used to get metric measures.
	Fetch(start time.Time, duration time.Duration) ([]DataPoint, error)
	// Name is used to serialize metrics.
	Name() string
}

// Aggregator provides a function to aggregate data points.
type Aggregator interface {
	// Aggregate calculates a single value for a metric datapoints.
	Aggregate([]DataPoint) float64
	// Name is used to serialize metrics.
	Name() string
}

type metricAggregator struct {
	repr string
	aggr func([]DataPoint) float64
}

func (m *metricAggregator) Aggregate(points []DataPoint) float64 {
	return m.aggr(points)
}

func (m *metricAggregator) Name() string {
	return m.repr
}

// MinimumAggregator returns the smallest datapoint value.
//nolint:gochecknoglobals
var MinimumAggregator Aggregator = &metricAggregator{
	repr: "MIN",
	aggr: func(points []DataPoint) float64 {
		if len(points) == 0 {
			return 0.
		}

		x := points[0].Value
		for _, point := range points {
			if point.Value < x {
				x = point.Value
			}
		}

		return x
	},
}

// MaximumAggregator returns the smallest datapoint value.
//nolint:gochecknoglobals
var MaximumAggregator Aggregator = &metricAggregator{
	repr: "MAX",
	aggr: func(points []DataPoint) float64 {
		if len(points) == 0 {
			return 0.
		}

		x := points[0].Value
		for _, point := range points {
			if point.Value > x {
				x = point.Value
			}
		}

		return x
	},
}

// AverageAggregator returns the average data point value.
//nolint:gochecknoglobals
var AverageAggregator Aggregator = &metricAggregator{
	repr: "AVG",
	aggr: func(points []DataPoint) float64 {
		if len(points) == 0 {
			return 0.
		}

		sum := 0.
		for _, point := range points {
			sum = sum + point.Value
		}

		return sum / float64(len(points))
	},
}

// Aggregators is a map of aggregators representation to the actual aggregator.
//nolint:gochecknoglobals
var Aggregators = map[string]Aggregator{
	MinimumAggregator.Name(): MinimumAggregator,
	MaximumAggregator.Name(): MaximumAggregator,
	AverageAggregator.Name(): AverageAggregator,
}

// ErrInvalidAggregator is returned when a metric aggregator cannot be found.
var ErrInvalidAggregator = errors.New("invalid aggregator")

// ParseAggregator returns metric aggregator from its name.
func ParseAggregator(name string) (Aggregator, error) {
	if aggregator, ok := Aggregators[name]; ok {
		return aggregator, nil
	}

	return nil, ErrInvalidAggregator
}
