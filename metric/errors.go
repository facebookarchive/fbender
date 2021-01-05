/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package metric

import (
	"time"

	"github.com/facebookincubator/fbender/recorders"
	"github.com/facebookincubator/fbender/tester"
	"github.com/pinterest/bender"
)

// ErrorsMetric fetches data from statistics.
type ErrorsMetric struct {
	Statistics recorders.Statistics
}

// ErrorsMetricOptions represents errors metric options.
type ErrorsMetricOptions interface {
	AddRecorder(bender.Recorder)
}

// Setup prepares errors metric.
func (m *ErrorsMetric) Setup(options interface{}) error {
	opts, ok := options.(ErrorsMetricOptions)
	if !ok {
		return tester.ErrInvalidOptions
	}

	opts.AddRecorder(recorders.NewStatisticsRecorder(&m.Statistics))

	return nil
}

// Fetch calculates the errors percentage for the statistics.
func (m *ErrorsMetric) Fetch(start time.Time, duration time.Duration) ([]tester.DataPoint, error) {
	errorsPct := float64(m.Statistics.Errors) / float64(m.Statistics.Requests) * 100.0

	// return a single point with time equal to end of the test
	return []tester.DataPoint{
		{Time: start.Add(duration), Value: errorsPct},
	}, nil
}

// Name returns the name of the errors statistic.
func (m *ErrorsMetric) Name() string {
	return "errors"
}
