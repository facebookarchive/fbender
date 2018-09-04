/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package recorders

import (
	"github.com/pinterest/bender"
	"github.com/sirupsen/logrus"
)

// NewLogrusRecorder creates a new logrus.Logger-based recorder.
func NewLogrusRecorder(l *logrus.Logger, defaults ...logrus.Fields) bender.Recorder {
	return func(msg interface{}) {
		log := logrus.NewEntry(l)
		for _, fields := range defaults {
			for key, value := range fields {
				log = log.WithField(key, value)
			}
		}
		switch msg := msg.(type) {
		case *bender.StartRequestEvent:
			logStartRequestEvent(log, msg)
		case *bender.EndRequestEvent:
			logEndRequestEvent(log, msg)
		}
		// Ignore other events
	}
}

func logStartRequestEvent(log *logrus.Entry, msg *bender.StartRequestEvent) {
	log.WithFields(logrus.Fields{
		"start":   msg.Time,
		"request": msg.Request,
	}).Debug("Start")
}

func logEndRequestEvent(log *logrus.Entry, msg *bender.EndRequestEvent) {
	log = log.WithFields(logrus.Fields{
		"start":    msg.Start,
		"end":      msg.End,
		"elapsed":  int(msg.End - msg.Start),
		"response": msg.Response,
	})
	if msg.Err != nil {
		log.WithError(msg.Err).Warn("Fail")
	} else {
		log.Info("Success")
	}
}
