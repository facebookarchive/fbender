/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package dns

import (
	"fmt"
	"time"

	"github.com/miekg/dns"
	"github.com/pinterest/bender"
	protocol "github.com/pinterest/bender/dns"
)

// Tester is a load tester for DHCPv6
type Tester struct {
	Target  string
	Timeout time.Duration
	client  *dns.Client
}

// Before is called before the first test
func (t *Tester) Before(options interface{}) error {
	t.client = &dns.Client{
		ReadTimeout:  t.Timeout,
		DialTimeout:  t.Timeout,
		WriteTimeout: t.Timeout,
	}
	return nil
}

// After is called after all tests are finished
func (t *Tester) After(_ interface{}) {}

// BeforeEach is called before every test
func (t *Tester) BeforeEach(_ interface{}) error {
	return nil
}

// AfterEach is called after every test
func (t *Tester) AfterEach(_ interface{}) {}

func validator(request, response *dns.Msg) error {
	if request.Id != response.Id {
		return fmt.Errorf("invalid response id: %d, want: %d", request.Id, response.Id)
	}
	return nil
}

// RequestExecutor returns a request executor
func (t *Tester) RequestExecutor(options interface{}) (bender.RequestExecutor, error) {
	return protocol.CreateExecutor(t.client, validator, t.Target), nil
}
