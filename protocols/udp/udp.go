/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package udp

import (
	"errors"
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/facebookincubator/fbender/log"
	"github.com/pinterest/bender"
)

// MaxResponseSize is the max response size for UDP test.
const MaxResponseSize = 2048

// Datagram represents a udp datagram to be sent.
type Datagram struct {
	Port int
	Data []byte
}

// ErrInvalidType is raised when object type mismatch.
var ErrInvalidType = errors.New("invalid type")

// ResponseValidator validates a udp response.
type ResponseValidator func(request *Datagram, response []byte) error

// CreateExecutor creates a new UDP RequestExecutor.
func CreateExecutor(timeout time.Duration, validator ResponseValidator, hosts ...string) bender.RequestExecutor {
	var i int

	return func(_ int64, request interface{}) (interface{}, error) {
		datagram, ok := request.(*Datagram)
		if !ok {
			return nil, fmt.Errorf("%w, want: *Datagram, got: %T", ErrInvalidType, request)
		}

		addr := net.JoinHostPort(hosts[i], strconv.Itoa(datagram.Port))
		i = (i + 1) % len(hosts)

		// Setup connection
		conn, err := net.Dial("udp", addr)
		if err != nil {
			//nolint:wrapcheck
			return nil, err
		}

		defer func() {
			if err = conn.Close(); err != nil {
				log.Errorf("Error closing connection: %v\n", err)
			}
		}()

		if err = conn.SetWriteDeadline(time.Now().Add(timeout)); err != nil {
			//nolint:wrapcheck
			return nil, err
		}

		_, err = conn.Write(datagram.Data)
		if err != nil {
			//nolint:wrapcheck
			return nil, err
		}

		buffer := make([]byte, MaxResponseSize)

		if err = conn.SetReadDeadline(time.Now().Add(timeout)); err != nil {
			//nolint:wrapcheck
			return nil, err
		}

		n, err := conn.Read(buffer)
		if err != nil {
			//nolint:wrapcheck
			return nil, err
		}

		if err = validator(datagram, buffer[:n]); err != nil {
			return nil, err
		}

		return buffer[:n], nil
	}
}
