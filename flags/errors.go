/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package flags

import (
	"errors"
)

// ErrUndefined is raised when user tries to access an undefined flag.
var ErrUndefined = errors.New("flag accessed but not defined")

// ErrInvalidType is raised when user tries to get a value from a flag of a different type.
var ErrInvalidType = errors.New("accessed flag type does not match")
