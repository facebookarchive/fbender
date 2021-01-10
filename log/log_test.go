/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package log_test

import (
	"bytes"
	"testing"

	"github.com/facebookincubator/fbender/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockedWriter struct {
	mock.Mock
}

func (m *MockedWriter) Write(p []byte) (int, error) {
	args := m.Called(p)

	return args.Int(0), args.Error(1)
}

func TestFprintf_NoError(t *testing.T) {
	w := new(bytes.Buffer)
	log.Fprintf(w, "Hello %s", "World")
	assert.Equal(t, "Hello World", w.String())
}

func TestFprintf_Error(t *testing.T) {
	w := new(MockedWriter)
	w.On("Write", "Hello World").Return(0, assert.AnError)

	defer func() {
		if r := recover(); r == nil {
			assert.Fail(t, "Expected Fprintf to panic")
		}
	}()

	log.Fprintf(w, "Hello %s", "World")
}

func TestPrintf_NoError(t *testing.T) {
	w := new(bytes.Buffer)
	log.Stdout = w
	log.Printf("Hello %s", "World")
	assert.Equal(t, "Hello World", w.String())
}

func TestPrintf_Error(t *testing.T) {
	w := new(MockedWriter)
	log.Stdout = w
	w.On("Write", "Hello World").Return(0, assert.AnError)

	defer func() {
		if r := recover(); r == nil {
			assert.Fail(t, "Expected Fprintf to panic")
		}
	}()

	log.Printf("Hello %s", "World")
}

func TestErrorf_NoError(t *testing.T) {
	w := new(bytes.Buffer)
	log.Stderr = w
	log.Errorf("Hello %s", "World")
	assert.Equal(t, "Hello World", w.String())
}

func TestErrorf_Error(t *testing.T) {
	w := new(MockedWriter)
	log.Stderr = w
	w.On("Write", "Hello World").Return(0, assert.AnError)

	defer func() {
		if r := recover(); r == nil {
			assert.Fail(t, "Expected Fprintf to panic")
		}
	}()

	log.Errorf("Hello %s", "World")
}
