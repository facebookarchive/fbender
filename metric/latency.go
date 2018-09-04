/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package metric

import (
	"sync"
	"time"

	"github.com/pinterest/bender"

	"github.com/facebookincubator/fbender/tester"
)

// LatencyMetric fetches data from statistics
type LatencyMetric struct {
	mutex  sync.Mutex
	points []tester.DataPoint
}

// LatencyMetricOptions represents errors metric options
type LatencyMetricOptions interface {
	AddRecorder(bender.Recorder)
	GetUnit() time.Duration
}

// Setup prepares errors metric
func (m *LatencyMetric) Setup(options interface{}) error {
	opts, ok := options.(LatencyMetricOptions)
	if !ok {
		return tester.ErrInvalidOptions
	}
	unit := opts.GetUnit()
	opts.AddRecorder(func(msg interface{}) {
		switch msg := msg.(type) {
		case *bender.StartEvent:
			m.points = make([]tester.DataPoint, 0)
		case *bender.EndRequestEvent:
			m.mutex.Lock()
			m.points = append(m.points, tester.DataPoint{
				Time:  time.Unix(msg.Start, 0),
				Value: float64(msg.End-msg.Start) / float64(unit),
			})
			m.mutex.Unlock()
		}
	})
	return nil
}

// Fetch calculates the errors percentage for the statistics
func (m *LatencyMetric) Fetch(start time.Time, duration time.Duration) ([]tester.DataPoint, error) {
	return m.points, nil
}

// Name returns the name of the errors statistic
func (m *LatencyMetric) Name() string {
	return "latency"
}
