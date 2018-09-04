/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package recorders

import (
	"sync/atomic"

	"github.com/pinterest/bender"
)

// Statistics groups statistics gathered by statistics recoreder
type Statistics struct {
	Requests int64
	Errors   int64
}

// Reset zeroes statistics
func (s *Statistics) Reset() {
	atomic.StoreInt64(&s.Requests, 0)
	atomic.StoreInt64(&s.Errors, 0)
}

// NewStatisticsRecorder creates new recorder which gathers statistics
func NewStatisticsRecorder(statistics *Statistics) bender.Recorder {
	return func(msg interface{}) {
		switch msg := msg.(type) {
		case *bender.StartEvent:
			statistics.Reset()
		case *bender.EndRequestEvent:
			atomic.AddInt64(&statistics.Requests, 1)
			if msg.Err != nil {
				atomic.AddInt64(&statistics.Errors, 1)
			}
		}
	}
}
