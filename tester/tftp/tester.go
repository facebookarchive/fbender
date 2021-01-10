/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package tftp

import (
	"fmt"
	"time"

	"github.com/pin/tftp"
	"github.com/pinterest/bender"
	protocol "github.com/pinterest/bender/tftp"
)

// Tester is a load tester for TFTP.
type Tester struct {
	Target    string
	Timeout   time.Duration
	BlockSize int

	client *tftp.Client
}

// Before is called before the first test.
func (t *Tester) Before(options interface{}) error {
	var err error

	t.client, err = tftp.NewClient(t.Target)
	if err != nil {
		return fmt.Errorf("unable to set up the tester: %w", err)
	}

	t.client.SetTimeout(t.Timeout)
	t.client.SetBlockSize(t.BlockSize)

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

// RequestExecutor returns a request executor.
func (t *Tester) RequestExecutor(_ interface{}) (bender.RequestExecutor, error) {
	return protocol.CreateExecutor(t.client, protocol.DiscardingValidator), nil
}
