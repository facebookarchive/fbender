/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package utils

import (
	"regexp"
)

// NamedRegex is a regex which supports named capture groups.
type NamedRegex struct {
	*regexp.Regexp
}

// FindStringSubmatchMap returns a map of named capture groups.
func (r *NamedRegex) FindStringSubmatchMap(s string) map[string]string {
	captures := make(map[string]string)
	match := r.FindStringSubmatch(s)

	if match == nil {
		return captures
	}

	for i, name := range r.SubexpNames() {
		// Ignore the whole regexp match and unnamed groups
		if i == 0 || name == "" {
			continue
		}

		captures[name] = match[i]
	}

	return captures
}

// MustCompile compiles a string to a named regexp.
func MustCompile(s string) NamedRegex {
	return NamedRegex{Regexp: regexp.MustCompile(s)}
}
