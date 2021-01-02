/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package recorders_test

import (
	"io/ioutil"
	"testing"

	"github.com/facebookincubator/fbender/recorders"
	"github.com/gosuri/uiprogress"
	"github.com/pinterest/bender"
	"github.com/stretchr/testify/suite"
)

type ProgressBarRecorderTestSuite struct {
	suite.Suite
	progress         *uiprogress.Progress
	bar              *uiprogress.Bar
	recorder         chan interface{}
	progressRecorder bender.Recorder
}

func (s *ProgressBarRecorderTestSuite) SetupTest() {
	s.progress = uiprogress.New()
	s.progress.SetOut(ioutil.Discard)
	s.progress.Start()
	s.bar = s.progress.AddBar(10)
	s.recorder = make(chan interface{}, 1)
	s.progressRecorder = recorders.NewProgressBarRecorder(s.bar)
}

func (s *ProgressBarRecorderTestSuite) TearDownTest() {
	s.progress.Stop()
}

func (s *ProgressBarRecorderTestSuite) recordSingleEvent(event interface{}) {
	s.recorder <- event
	close(s.recorder)
	bender.Record(s.recorder, s.progressRecorder)
}

func (s *ProgressBarRecorderTestSuite) TestStartEvent() {
	s.recordSingleEvent(new(bender.StartEvent))
	s.Equal(0, s.bar.Current())
}

func (s *ProgressBarRecorderTestSuite) TestEndEvent() {
	s.recordSingleEvent(new(bender.EndEvent))
	s.Equal(0, s.bar.Current())
}

func (s *ProgressBarRecorderTestSuite) TestWaitEvent() {
	s.recordSingleEvent(new(bender.WaitEvent))
	s.Equal(0, s.bar.Current())
}

func (s *ProgressBarRecorderTestSuite) TestStartRequestEvent() {
	s.recordSingleEvent(new(bender.StartRequestEvent))
	s.Equal(0, s.bar.Current())
}

func (s *ProgressBarRecorderTestSuite) TestEndRequestEvent() {
	s.recordSingleEvent(new(bender.EndRequestEvent))
	s.Equal(1, s.bar.Current())
}

func TestProgressBarRecorderTestSuite(t *testing.T) {
	suite.Run(t, new(ProgressBarRecorderTestSuite))
}
