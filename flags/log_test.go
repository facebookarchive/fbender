/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package flags_test

import (
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/facebookincubator/fbender/flags"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//nolint:gochecknoglobals
var (
	rand   uint32
	randmu sync.Mutex
)

func tempFilename(dir, prefix string) string {
	if dir == "" {
		dir = os.TempDir()
	}

	randmu.Lock()

	r := rand
	if r == 0 {
		r = uint32(time.Now().UnixNano() + int64(os.Getpid()))
	}

	r = r*1664525 + 1013904223 // constants from Numerical Recipes
	rand = r

	randmu.Unlock()

	return filepath.Join(dir, prefix+strconv.Itoa(int(1e9 + r%1e9))[1:])
}

func TestLogLevelChoices(t *testing.T) {
	expected := []string{"panic", "fatal", "error", "warning", "info", "debug", "trace"}
	assert.ElementsMatch(t, flags.LogLevelChoices(), expected)
}

func TestLogLevel__String(t *testing.T) {
	logLevel := &flags.LogLevel{Logger: logrus.New()}
	logLevel.Logger.SetLevel(logrus.TraceLevel)
	assert.Equal(t, "trace", logLevel.String())
	logLevel.Logger.SetLevel(logrus.DebugLevel)
	assert.Equal(t, "debug", logLevel.String())
	logLevel.Logger.SetLevel(logrus.InfoLevel)
	assert.Equal(t, "info", logLevel.String())
	logLevel.Logger.SetLevel(logrus.WarnLevel)
	assert.Equal(t, "warning", logLevel.String())
	logLevel.Logger.SetLevel(logrus.ErrorLevel)
	assert.Equal(t, "error", logLevel.String())
	logLevel.Logger.SetLevel(logrus.FatalLevel)
	assert.Equal(t, "fatal", logLevel.String())
	logLevel.Logger.SetLevel(logrus.PanicLevel)
	assert.Equal(t, "panic", logLevel.String())
}

//nolint:funlen
func TestLogLevel__Set(t *testing.T) {
	logger := logrus.New()
	logLevel := &flags.LogLevel{Logger: logger}

	// Setting proper level by name should change the Level value
	err := logLevel.Set("trace")
	assert.NoError(t, err)

	err = logLevel.Set("debug")
	assert.NoError(t, err)
	assert.Equal(t, logrus.DebugLevel, logger.Level)

	err = logLevel.Set("info")
	assert.NoError(t, err)
	assert.Equal(t, logrus.InfoLevel, logger.Level)

	err = logLevel.Set("warning")
	assert.NoError(t, err)
	assert.Equal(t, logrus.WarnLevel, logger.Level)

	err = logLevel.Set("error")
	assert.NoError(t, err)
	assert.Equal(t, logrus.ErrorLevel, logger.Level)

	err = logLevel.Set("fatal")
	assert.NoError(t, err)
	assert.Equal(t, logrus.FatalLevel, logger.Level)

	err = logLevel.Set("panic")
	assert.NoError(t, err)
	assert.Equal(t, logrus.PanicLevel, logger.Level)

	// Setting proper level by number should change the Level value
	err = logLevel.Set("6")
	assert.NoError(t, err)

	err = logLevel.Set("5")
	assert.NoError(t, err)
	assert.Equal(t, logrus.DebugLevel, logger.Level)

	err = logLevel.Set("4")
	assert.NoError(t, err)
	assert.Equal(t, logrus.InfoLevel, logger.Level)

	err = logLevel.Set("3")
	assert.NoError(t, err)
	assert.Equal(t, logrus.WarnLevel, logger.Level)

	err = logLevel.Set("2")
	assert.NoError(t, err)
	assert.Equal(t, logrus.ErrorLevel, logger.Level)

	err = logLevel.Set("1")
	assert.NoError(t, err)
	assert.Equal(t, logrus.FatalLevel, logger.Level)

	err = logLevel.Set("0")
	assert.NoError(t, err)
	assert.Equal(t, logrus.PanicLevel, logger.Level)

	// Setting unknown level string should fail
	err = logLevel.Set("unknown")
	assert.EqualError(t, err, "level must be an integer [0..6] or (panic|fatal|error|warning|info|debug|trace)")

	// Setting too large integer
	err = logLevel.Set("7")
	assert.EqualError(t, err, "level must be an integer [0..6] or (panic|fatal|error|warning|info|debug|trace)")
}

func TestLogLevel__Type(t *testing.T) {
	logLevel := new(flags.LogLevel)
	assert.Equal(t, "level", logLevel.Type())
}

func TestBashCompletionLogLevel(t *testing.T) {
	c := &cobra.Command{}
	l := &flags.LogLevel{Logger: logrus.New()}

	// Check no error when applied to level flag
	f := c.Flags().VarPF(l, "level", "l", "set level")

	err := flags.BashCompletionLogLevel(c, c.Flags(), "level")
	require.NoError(t, err)

	require.Contains(t, f.Annotations, "cobra_annotation_bash_completion_custom")
	assert.Equal(t, []string{"__fbender_handle_loglevel_flag"},
		f.Annotations["cobra_annotation_bash_completion_custom"])
	assert.Equal(t, `
__fbender_handle_loglevel_flag() {
	COMPREPLY=($(compgen -W "panic fatal error warning info debug trace" -- "${cur}"))
}`, c.BashCompletionFunction)

	// Check error when flag is not defined
	err = flags.BashCompletionLogLevel(c, c.Flags(), "nonexistent")
	assert.EqualError(t, err, "flag nonexistent accessed but not defined")

	// Check error when flag is not a level
	c.Flags().Int("myint", 0, "set myint")

	err = flags.BashCompletionLogLevel(c, c.Flags(), "myint")
	assert.EqualError(t, err, "trying to autocomplete level on flag of type int")
}

func TestLogFormatChoices(t *testing.T) {
	expected := []string{"text", "json"}
	assert.ElementsMatch(t, flags.LogFormatChoices(), expected)
}

func TestLogFormat__String(t *testing.T) {
	logFormat := &flags.LogFormat{Format: "json"}
	assert.Equal(t, "json", logFormat.String())

	logFormat.Format = "text"
	assert.Equal(t, "text", logFormat.String())
}

func TestLogFormat__Set(t *testing.T) {
	logger := logrus.New()
	logFormat := &flags.LogFormat{Logger: logger}

	// Setting proper level by name should change the Formatter
	err := logFormat.Set("text")
	require.NoError(t, err)
	assert.Equal(t, "text", logFormat.Format)
	assert.IsType(t, new(logrus.TextFormatter), logger.Formatter)

	err = logFormat.Set("json")
	require.NoError(t, err)
	assert.Equal(t, "json", logFormat.Format)
	assert.IsType(t, new(logrus.JSONFormatter), logger.Formatter)

	// Setting unknown format string should fail
	err = logFormat.Set("unknown")
	assert.EqualError(t, err, "logformat must be one of (json|text)")
}

func TestLogFormat__Type(t *testing.T) {
	logFormat := new(flags.LogFormat)
	assert.Equal(t, "format", logFormat.Type())
}

func TestBashCompletionLogFormat(t *testing.T) {
	c := &cobra.Command{}
	l := &flags.LogFormat{Logger: logrus.New()}

	// Check no error when applied to format flag
	f := c.Flags().VarPF(l, "format", "f", "set format")
	err := flags.BashCompletionLogFormat(c, c.Flags(), "format")
	require.NoError(t, err)
	require.Contains(t, f.Annotations, "cobra_annotation_bash_completion_custom")
	assert.Equal(t, []string{"__fbender_handle_logformat_flag"},
		f.Annotations["cobra_annotation_bash_completion_custom"])
	assert.Equal(t, `
__fbender_handle_logformat_flag() {
	COMPREPLY=($(compgen -W "json text" -- "${cur}"))
}`, c.BashCompletionFunction)

	// Check error when flag is not defined
	err = flags.BashCompletionLogFormat(c, c.Flags(), "nonexistent")
	assert.EqualError(t, err, "flag nonexistent accessed but not defined")

	// Check error when flag is not a format
	c.Flags().Int("myint", 0, "set myint")
	err = flags.BashCompletionLogFormat(c, c.Flags(), "myint")
	assert.EqualError(t, err, "trying to autocomplete format on flag of type int")
}

func TestNewLogOutput(t *testing.T) {
	logger := logrus.New()
	logOutput := flags.NewLogOutput(logger)
	assertPointerEqual(t, logger, logOutput.Logger)
	assert.Equal(t, "<stdout>", logOutput.String())
	assert.Equal(t, os.Stdout, logger.Out)
}

func TestLogOutput__String(t *testing.T) {
	logOutput := flags.NewLogOutput(logrus.New())
	assert.Equal(t, "<stdout>", logOutput.String())

	err := logOutput.Set("")
	require.NoError(t, err)
	assert.Equal(t, "<stdout>", logOutput.String())

	logOutput.Out = os.Stderr
	assert.Equal(t, "<stderr>", logOutput.String())

	filename := tempFilename("", "testlogoutput__string")

	err = logOutput.Set(filename)
	require.NoError(t, err)

	defer func() {
		err = logOutput.Out.Close()
		require.NoError(t, err)
	}()

	assert.Equal(t, filename, logOutput.String())
}

func TestLogOutput__Set(t *testing.T) {
	logOutput := flags.NewLogOutput(logrus.New())
	filename := tempFilename("", "testlogoutput__set")

	err := logOutput.Set(filename)
	require.NoError(t, err)

	// Check if logOutput opened proper file
	file, err := os.Open(filename)
	require.NoError(t, err)

	fileStat, err := file.Stat()
	require.NoError(t, err)

	err = file.Close()
	require.NoError(t, err)

	logFileStat, err := logOutput.Out.Stat()
	require.NoError(t, err)

	assert.True(t, os.SameFile(fileStat, logFileStat), "file opened by the logger is not the expected file")
	// Check if setting a different file closes the first one
	logFile := logOutput.Out

	err = logOutput.Set("")
	require.NoError(t, err)

	err = logFile.Close()
	require.Error(t, err)

	require.Contains(t, err.Error(), "file already closed")
}

func TestLogOutput__Type(t *testing.T) {
	logOutput := flags.NewLogOutput(logrus.New())
	assert.Equal(t, "output", logOutput.Type())
}
