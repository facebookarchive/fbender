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

	"github.com/facebookincubator/fbender/recorders"
	"github.com/pinterest/bender"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/suite"
)

type LogrusRecorderTestSuite struct {
	suite.Suite
	logger         *logrus.Logger
	hook           *test.Hook
	recorder       chan interface{}
	logrusRecorder bender.Recorder
}

func (s *LogrusRecorderTestSuite) SetupSuite() {
	s.logger, s.hook = test.NewNullLogger()
	s.logger.Level = logrus.DebugLevel
}

func (s *LogrusRecorderTestSuite) SetupTest() {
	s.hook.Reset()
	s.recorder = make(chan interface{}, 1)
	s.logrusRecorder = recorders.NewLogrusRecorder(s.logger)
}

func (s *LogrusRecorderTestSuite) recordSingleEvent(event interface{}) {
	s.recorder <- event
	close(s.recorder)
	bender.Record(s.recorder, s.logrusRecorder)
}

func (s *LogrusRecorderTestSuite) TestStartEvent() {
	s.recordSingleEvent(new(bender.StartEvent))
	s.Assert().Len(s.hook.Entries, 0)
}

func (s *LogrusRecorderTestSuite) TestEndEvent() {
	s.recordSingleEvent(new(bender.EndEvent))
	s.Assert().Len(s.hook.Entries, 0)
}

func (s *LogrusRecorderTestSuite) TestWaitEvent() {
	s.recordSingleEvent(new(bender.WaitEvent))
	s.Assert().Len(s.hook.Entries, 0)
}

func (s *LogrusRecorderTestSuite) TestStartRequestEvent() {
	event := &bender.StartRequestEvent{
		Time:    420,
		Request: "my request",
	}
	s.recordSingleEvent(event)
	s.Require().Len(s.hook.Entries, 1)
	s.Assert().Equal(logrus.DebugLevel, s.hook.LastEntry().Level)
	s.Assert().Equal("Start", s.hook.LastEntry().Message)
	s.Assert().Equal(logrus.Fields{
		"start":   int64(420),
		"request": "my request",
	}, s.hook.LastEntry().Data)

	s.hook.Reset()
	s.recorder = make(chan interface{}, 1)
	s.logrusRecorder = recorders.NewLogrusRecorder(s.logger, logrus.Fields{
		"keyA0": "valueA0",
		"keyA1": "valueA1",
	}, logrus.Fields{
		"keyB0": "valueB0",
		"keyB1": "valueB1",
	})
	s.recordSingleEvent(event)
	s.Require().Len(s.hook.Entries, 1)
	s.Assert().Equal(logrus.DebugLevel, s.hook.LastEntry().Level)
	s.Assert().Equal("Start", s.hook.LastEntry().Message)
	s.Assert().Equal(logrus.Fields{
		"keyA0":   "valueA0",
		"keyA1":   "valueA1",
		"keyB0":   "valueB0",
		"keyB1":   "valueB1",
		"start":   int64(420),
		"request": "my request",
	}, s.hook.LastEntry().Data)
}

func (s *LogrusRecorderTestSuite) TestEndRequestEvent_NoError() {
	event := &bender.EndRequestEvent{
		Start:    420,
		End:      4200,
		Response: "my response",
		Err:      nil,
	}
	s.recordSingleEvent(event)
	s.Require().Len(s.hook.Entries, 1)
	s.Assert().Equal(logrus.InfoLevel, s.hook.LastEntry().Level)
	s.Assert().Equal("Success", s.hook.LastEntry().Message)
	s.Assert().Equal(logrus.Fields{
		"start":    int64(420),
		"end":      int64(4200),
		"elapsed":  3780,
		"response": "my response",
	}, s.hook.LastEntry().Data)

	s.hook.Reset()
	s.recorder = make(chan interface{}, 1)
	s.logrusRecorder = recorders.NewLogrusRecorder(s.logger, logrus.Fields{
		"keyA0": "valueA0",
		"keyA1": "valueA1",
	}, logrus.Fields{
		"keyB0": "valueB0",
		"keyB1": "valueB1",
	})
	s.recordSingleEvent(event)
	s.Require().Len(s.hook.Entries, 1)
	s.Assert().Equal(logrus.InfoLevel, s.hook.LastEntry().Level)
	s.Assert().Equal("Success", s.hook.LastEntry().Message)
	s.Assert().Equal(logrus.Fields{
		"keyA0":    "valueA0",
		"keyA1":    "valueA1",
		"keyB0":    "valueB0",
		"keyB1":    "valueB1",
		"start":    int64(420),
		"end":      int64(4200),
		"elapsed":  3780,
		"response": "my response",
	}, s.hook.LastEntry().Data)
}

func (s *LogrusRecorderTestSuite) TestEndRequestEvent_Error() {
	event := &bender.EndRequestEvent{
		Start:    420,
		End:      4200,
		Response: "invalid response",
		Err:      errors.New("invalid response"),
	}
	s.recordSingleEvent(event)
	s.Require().Len(s.hook.Entries, 1)
	s.Assert().Equal(logrus.WarnLevel, s.hook.LastEntry().Level)
	s.Assert().Equal("Fail", s.hook.LastEntry().Message)
	s.Assert().Equal(logrus.Fields{
		"start":    int64(420),
		"end":      int64(4200),
		"elapsed":  3780,
		"response": "invalid response",
		"error":    errors.New("invalid response"),
	}, s.hook.LastEntry().Data)

	s.hook.Reset()
	s.recorder = make(chan interface{}, 1)
	s.logrusRecorder = recorders.NewLogrusRecorder(s.logger, logrus.Fields{
		"keyA0": "valueA0",
		"keyA1": "valueA1",
	}, logrus.Fields{
		"keyB0": "valueB0",
		"keyB1": "valueB1",
	})
	s.recordSingleEvent(event)
	s.Require().Len(s.hook.Entries, 1)
	s.Assert().Equal(logrus.WarnLevel, s.hook.LastEntry().Level)
	s.Assert().Equal("Fail", s.hook.LastEntry().Message)
	s.Assert().Equal(logrus.Fields{
		"keyA0":    "valueA0",
		"keyA1":    "valueA1",
		"keyB0":    "valueB0",
		"keyB1":    "valueB1",
		"start":    int64(420),
		"end":      int64(4200),
		"elapsed":  3780,
		"response": "invalid response",
		"error":    errors.New("invalid response"),
	}, s.hook.LastEntry().Data)
}

func TestLogrusRecorderTestSuite(t *testing.T) {
	suite.Run(t, new(LogrusRecorderTestSuite))
}
