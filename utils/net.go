/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package utils

import (
	"net"
	"strconv"
)

// WithDefaultPort adds a default port if no port is present.
func WithDefaultPort(hostport string, port int) string {
	if _, _, err := net.SplitHostPort(hostport); err != nil {
		return net.JoinHostPort(hostport, strconv.Itoa(port))
	}

	return hostport
}
