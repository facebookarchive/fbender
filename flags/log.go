/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package flags

import (
	"fmt"
	"os"
	"sort"
	"strconv"

	"github.com/sirupsen/logrus"
)

// LogLevel represents a log level flag value
type LogLevel struct {
	Logger *logrus.Logger
}

// LogLevelChoices returns a string representation of available levels
func LogLevelChoices() []string {
	choices := []string{}
	for _, level := range logrus.AllLevels {
		choices = append(choices, level.String())
	}
	return choices
}

func (l *LogLevel) String() string {
	return l.Logger.Level.String()
}

// Set validates a given value and sets log level
func (l *LogLevel) Set(value string) error {
	level, err := logrus.ParseLevel(value)
	if err != nil {
		v, err := strconv.ParseUint(value, 10, 32)
		if err != nil || v >= uint64(len(logrus.AllLevels)) {
			choices := ChoicesString(LogLevelChoices())
			return fmt.Errorf("level must be an integer [0..%d] or %s", len(logrus.AllLevels)-1, choices)
		}
		level = logrus.Level(v)
	}
	l.Logger.SetLevel(level)
	return nil
}

// Type returns a log level value type
func (l *LogLevel) Type() string {
	return "level"
}

// LogFormat represents a log format flag value
type LogFormat struct {
	Logger *logrus.Logger
	Format string
}

const (
	testFormatter = "text"
	jsonFormatter = "json"
)

// We need to store a lambda function which generates formatters so we don't
// assign the same static formatter to different Loggers
var formatters = map[string](func() logrus.Formatter){
	testFormatter: func() logrus.Formatter { return new(logrus.TextFormatter) },
	jsonFormatter: func() logrus.Formatter { return new(logrus.JSONFormatter) },
}

// LogFormatChoices returns a string representation of available formats
func LogFormatChoices() []string {
	choices := []string{}
	for key := range formatters {
		choices = append(choices, key)
	}
	sort.Strings(choices)
	return choices
}

func (l *LogFormat) String() string {
	return l.Format
}

// Set validates a given value and sets log format
func (l *LogFormat) Set(value string) error {
	if formatter, ok := formatters[value]; ok {
		l.Logger.Formatter = formatter()
		l.Format = value
		return nil
	}
	return fmt.Errorf("logformat must be one of %s", ChoicesString(LogFormatChoices()))
}

// Type returns a log format value type
func (l *LogFormat) Type() string {
	return "format"
}

// LogOutput represents a log output flag value
type LogOutput struct {
	Logger   *logrus.Logger
	Filename string
	Out      *os.File
}

// NewLogOutput returns new log output flag Value
func NewLogOutput(logger *logrus.Logger) *LogOutput {
	logger.Out = os.Stdout
	return &LogOutput{
		Logger: logger,
		Out:    os.Stdout,
	}
}

func (l *LogOutput) String() string {
	if l.Out == os.Stdout {
		return "<stdout>"
	} else if l.Out == os.Stderr {
		return "<stderr>"
	}
	return l.Out.Name()
}

// Set opens the file and sets the Out of the Logger
func (l *LogOutput) Set(value string) error {
	// Close any file we have been writing to
	if l.Out != nil && l.Out != os.Stdout && l.Out != os.Stderr {
		if err := l.Out.Close(); err != nil {
			return err
		}
	}

	// If the value is not an empty string it is a file name
	if len(value) > 0 {
		file, err := os.Create(value)
		if err != nil {
			return err
		}
		l.Out = file
	} else {
		l.Out = os.Stdout
	}

	l.Logger.Out = l.Out
	return nil
}

// Type returns a log output value type
func (l *LogOutput) Type() string {
	return "output"
}
