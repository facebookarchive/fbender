/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

/*
Package log implements standard logging.
*/
package log

import (
	"fmt"
	"io"
	"os"
)

// Standard writers for the log package
var (
	Stderr io.Writer = os.Stderr
	Stdout io.Writer = os.Stdout
)

// Fprintf formats according to a format specifier and writes to w.
// It panics if any write error encountered.
func Fprintf(w io.Writer, format string, args ...interface{}) {
	if _, err := fmt.Fprintf(w, format, args...); err != nil {
		panic(err)
	}
}

// Printf formats according to a format specifier and writes to standard output.
// It panics if any write error encountered.
func Printf(format string, args ...interface{}) {
	Fprintf(Stdout, format, args...)
}

// Errorf formats according to a format specifier and writes to standard error.
// It panics if any write error encountered.
func Errorf(format string, args ...interface{}) {
	Fprintf(Stderr, format, args...)
}
