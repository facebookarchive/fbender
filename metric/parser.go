/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package metric

import (
	"github.com/facebookincubator/fbender/tester"
)

// Help is a help message on available metrics.
const Help = `
Basic Metrics:
* errors - errors percentage of all requests, ignores aggregator
* latency - latency of the packets (in unit specified by --unit)
  MAX(errors) < 10.0
  MIN(errors) < 42.0
  AVG(latency) < 30`

// Parser is a parser for standard metrics.
func Parser(value string) (tester.Metric, error) {
	switch value {
	case "errors":
		return new(ErrorsMetric), nil
	case "latency":
		return new(LatencyMetric), nil
	default:
		return nil, tester.ErrNotParsed
	}
}
