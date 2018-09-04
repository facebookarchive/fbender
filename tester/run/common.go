/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package run

import (
	"time"

	"github.com/facebookincubator/fbender/log"
	"github.com/facebookincubator/fbender/tester"
)

// checkConstraints loops through given constraints and returns whether all of
// them have been met.
func checkConstraints(start time.Time, duration time.Duration, constraints ...*tester.Constraint) bool {
	for _, constraint := range constraints {
		if err := constraint.Check(start, duration); err != nil {
			log.Errorf("Error checking %q: %v\n", constraint.String(), err)
			return false
		}
	}
	return true
}
