/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package flags

import (
	"fmt"
	"strings"
)

// ChoicesString converts choices list into printable text
func ChoicesString(choices []string) string {
	return fmt.Sprintf("(%s)", strings.Join(choices, "|"))
}
