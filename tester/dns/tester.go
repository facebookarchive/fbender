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

// ExtendedMsg wraps a dns.Msg with expectations
type ExtendedMsg struct {
	dns.Msg
	Rcode int
}

// Tester is a load tester for DNS
type Tester struct {
	Target   string
	Timeout  time.Duration
	Protocol string
	client   *dns.Client
}

// Before is called before the first test
func (t *Tester) Before(options interface{}) error {
	t.client = &dns.Client{
		ReadTimeout:  t.Timeout,
		DialTimeout:  t.Timeout,
		WriteTimeout: t.Timeout,
		Net:          t.Protocol,
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
	innerExecutor := protocol.CreateExecutor(t.client, validator, t.Target)

	return func(n int64, request interface{}) (interface{}, error) {
		asExtended, ok := request.(*ExtendedMsg)
		if !ok {
			return nil, fmt.Errorf("request type is not ExtendedMsg")
		}
		resp, err := innerExecutor(n, &asExtended.Msg)
		if err != nil {
			return resp, err
		}
		asMsg, ok := resp.(*dns.Msg)
		if !ok {
			return nil, fmt.Errorf("reponse type is not dns.Msg")
		}
		if asExtended.Rcode != -1 && asExtended.Rcode != asMsg.Rcode {
			return resp, fmt.Errorf(
				"invalid rcode %s, want: %s",
				dns.RcodeToString[asMsg.Rcode], dns.RcodeToString[asExtended.Rcode])
		}
		return resp, nil
	}, nil
}
