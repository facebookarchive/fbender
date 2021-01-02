/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package udp

import (
	"time"

	protocol "github.com/facebookincubator/fbender/protocols/udp"
	"github.com/pinterest/bender"
)

// Tester is a load tester for UDP.
type Tester struct {
	Target    string
	Timeout   time.Duration
	Validator protocol.ResponseValidator
}

// Before is called before the first test.
func (t *Tester) Before(_ interface{}) error {
	return nil
}

// After is called after all tests are finished.
func (t *Tester) After(_ interface{}) {}

// BeforeEach is called before every test.
func (t *Tester) BeforeEach(_ interface{}) error {
	return nil
}

// AfterEach is called after every test.
func (t *Tester) AfterEach(_ interface{}) {}

func validator(_ *protocol.Datagram, _ []byte) error {
	return nil
}

// RequestExecutor returns a request executor.
func (t *Tester) RequestExecutor(options interface{}) (bender.RequestExecutor, error) {
	if t.Validator == nil {
		return protocol.CreateExecutor(t.Timeout, validator, t.Target), nil
	}

	return protocol.CreateExecutor(t.Timeout, t.Validator, t.Target), nil
}
