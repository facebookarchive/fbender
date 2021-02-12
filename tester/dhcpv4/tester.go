/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package dhcpv4

import (
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/insomniacslk/dhcp/dhcpv4"
	"github.com/insomniacslk/dhcp/dhcpv4/async"
	"github.com/pinterest/bender"
	protocol "github.com/pinterest/bender/dhcpv4"
)

// Tester is a load tester for DHCPv4.
type Tester struct {
	Target     string
	Timeout    time.Duration
	BufferSize int
	client     *async.Client
}

// Before is called before the first test.
func (t *Tester) Before(options interface{}) error {
	target, err := net.ResolveUDPAddr("udp4", t.Target)
	if err != nil {
		return fmt.Errorf("unable to set up the tester: %w", err)
	}

	addr, err := getLocalIPv4("eth0")
	if err != nil {
		return fmt.Errorf("unable to set up the tester: %w", err)
	}

	t.client = &async.Client{
		ReadTimeout:  t.Timeout,
		WriteTimeout: t.Timeout,
		RemoteAddr:   target,
		LocalAddr:    &net.UDPAddr{IP: addr, Port: async.DefaultServerPort},
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

func validator(req, res *dhcpv4.DHCPv4) error {
	return nil
}

// RequestExecutor returns a request executor.
func (t *Tester) RequestExecutor(_ interface{}) (bender.RequestExecutor, error) {
	return protocol.CreateExecutor(t.client, validator)
}

// ErrNoAddress is raised when an interface has no ipv4 addresses assigned.
var ErrNoAddress = errors.New("no ipv4 address found")

// getLocalIPv4 returns the interface local IPv4 address.
func getLocalIPv4(ifname string) (net.IP, error) {
	iface, err := net.InterfaceByName(ifname)
	if err != nil {
		//nolint:wrapcheck
		return nil, err
	}

	ifaddrs, err := iface.Addrs()
	if err != nil {
		//nolint:wrapcheck
		return nil, err
	}

	for _, ifaddr := range ifaddrs {
		if ipnet, ok := ifaddr.(*net.IPNet); ok {
			if ipnet.IP.To4() != nil && !ipnet.IP.IsLoopback() && !ipnet.IP.IsLinkLocalUnicast() {
				return ipnet.IP, nil
			}
		}
	}

	return nil, fmt.Errorf("%w, interface: %s", ErrNoAddress, ifname)
}
