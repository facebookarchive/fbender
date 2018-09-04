/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package recorders_test

import (
	"errors"
	"testing"

	"github.com/pinterest/bender"
	"github.com/stretchr/testify/suite"

	"github.com/facebookincubator/fbender/recorders"
)

type StatisticsRecorderTestSuite struct {
	suite.Suite
	statistics         *recorders.Statistics
	recorder           chan interface{}
	statisticsRecorder bender.Recorder
}

func (s *StatisticsRecorderTestSuite) SetupTest() {
	s.recorder = make(chan interface{}, 1)
	s.statistics = new(recorders.Statistics)
	s.statisticsRecorder = recorders.NewStatisticsRecorder(s.statistics)
}

func (s *StatisticsRecorderTestSuite) recordSingleEvent(event interface{}) {
	s.recorder <- event
	close(s.recorder)
	bender.Record(s.recorder, s.statisticsRecorder)
}

func (s *StatisticsRecorderTestSuite) TestStartEvent() {
	s.statistics.Requests = 42
	s.statistics.Errors = 6
	s.recordSingleEvent(new(bender.StartEvent))
	s.Equal(int64(0), s.statistics.Requests)
	s.Equal(int64(0), s.statistics.Errors)
}

func (s *StatisticsRecorderTestSuite) TestEndEvent() {
	s.recordSingleEvent(new(bender.EndEvent))
	s.Equal(int64(0), s.statistics.Requests)
	s.Equal(int64(0), s.statistics.Errors)
}

func (s *StatisticsRecorderTestSuite) TestWaitEvent() {
	s.recordSingleEvent(new(bender.WaitEvent))
	s.Equal(int64(0), s.statistics.Requests)
	s.Equal(int64(0), s.statistics.Errors)
}

func (s *StatisticsRecorderTestSuite) TestStartRequestEvent() {
	s.recordSingleEvent(new(bender.StartRequestEvent))
	s.Equal(int64(0), s.statistics.Requests)
	s.Equal(int64(0), s.statistics.Errors)
}

func (s *StatisticsRecorderTestSuite) TestEndRequestEvent_NoError() {
	s.recordSingleEvent(&bender.EndRequestEvent{Err: nil})
	s.Equal(int64(1), s.statistics.Requests)
	s.Equal(int64(0), s.statistics.Errors)
}

func (s *StatisticsRecorderTestSuite) TestEndRequestEvent_Error() {
	s.recordSingleEvent(&bender.EndRequestEvent{Err: errors.New("with error")})
	s.Equal(int64(1), s.statistics.Requests)
	s.Equal(int64(1), s.statistics.Errors)
}

func TestStatisticsRecorderTestSuite(t *testing.T) {
	suite.Run(t, new(StatisticsRecorderTestSuite))
}
