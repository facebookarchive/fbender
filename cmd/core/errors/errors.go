/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package errors

import "errors"

// ErrInvalidFormat is raised when the input does not match the desired format.
var ErrInvalidFormat = errors.New("invalid format")

// ErrInvalidType is raised when provided object is not of the desired type.
var ErrInvalidType = errors.New("invalid type")

// ErrInvalidArgument is raised when the command is given invalid arguments.
var ErrInvalidArgument = errors.New("invalid argument")
