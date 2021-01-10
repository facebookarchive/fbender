/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package dhcpv6

import (
	"fmt"
	"net"
	"time"

	"github.com/insomniacslk/dhcp/dhcpv6"
	"github.com/insomniacslk/dhcp/dhcpv6/async"
	"github.com/pinterest/bender"
	protocol "github.com/pinterest/bender/dhcpv6"
)

// Tester is a load tester for DHCPv6.
type Tester struct {
	Target     string
	Timeout    time.Duration
	BufferSize int
	client     *async.Client
}

// Before is called before the first test.
func (t *Tester) Before(options interface{}) error {
	target, err := net.ResolveUDPAddr("udp6", t.Target)
	if err != nil {
		return fmt.Errorf("unable to set up the tester: %w", err)
	}

	ip, err := dhcpv6.GetGlobalAddr("eth0")
	if err != nil {
		return fmt.Errorf("unable to set up the tester: %w", err)
	}

	t.client = &async.Client{
		ReadTimeout:  t.Timeout,
		WriteTimeout: t.Timeout,
		LocalAddr:    &net.UDPAddr{IP: ip, Port: dhcpv6.DefaultServerPort, Zone: ""},
		RemoteAddr:   target,
		IgnoreErrors: true,
	}

	return nil
}

// After is called after all tests are finished.
func (t *Tester) After(_ interface{}) {}

// BeforeEach is called before every test.
func (t *Tester) BeforeEach(options interface{}) error {
	return t.client.Open(t.BufferSize)
}

// AfterEach is called after every test.
func (t *Tester) AfterEach(_ interface{}) {
	t.client.Close()
}

func validator(req, res dhcpv6.DHCPv6) error {
	return nil
}

// RequestExecutor returns a request executor.
func (t *Tester) RequestExecutor(_ interface{}) (bender.RequestExecutor, error) {
	return protocol.CreateExecutor(t.client, validator), nil
}
