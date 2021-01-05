/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package tester

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/facebookincubator/fbender/utils"
)

// Constraint represents a constraint tests should meet to be considered
// successful.
type Constraint struct {
	Metric     Metric
	Aggregator Aggregator
	Comparator Comparator
	Threshold  float64
}

func (c *Constraint) String() string {
	return fmt.Sprintf("%s(%s) %s %.2f",
		c.Aggregator.Name(), c.Metric.Name(), c.Comparator.Name(), c.Threshold)
}

// Check fetches metric and checks if the constraint has been satisfied.
func (c *Constraint) Check(start time.Time, duration time.Duration) error {
	points, err := c.Metric.Fetch(start, duration)
	if err != nil {
		return err
	}

	if points == nil {
		return errors.New("no data points")
	}

	value := c.Aggregator.Aggregate(points)
	if !c.Comparator.Compare(value, c.Threshold) {
		return fmt.Errorf("unsatisfied %.4f %s %.4f", value, c.Comparator.Name(), c.Threshold)
	}

	return nil
}

// ConstraintsHelp is an help message on how to use constraints.
const ConstraintsHelp = `
Constraints follow the syntax:
  Constraint ::= <Aggregator>(<Metric>)<Cmp><Threshold>
  Aggregator ::= "MIN" | "MAX"
  Metric     ::= <string>
  Cmp        ::= "<" | ">"
  Threshold  ::= <float>

Constraints examples:
  MIN(metric) < 20.5
  MAX(metric) > 0.45
  MIN(metric) < 123

` + GrowthHelp

// ErrNotParsed should be returned when a parser did not parse a constraint.
var ErrNotParsed = errors.New("constraint could not be parsed")

// MetricParser is used to parse string values to a metric.
// Parsers should return a metric and error if it successfully parsed
// a metric string, or a fatal error occurred. Otherwise it should return
// ErrNotParsed which will result in trying next parser from the list.
type MetricParser func(string) (Metric, error)

// Named capture groups of the constraints matching regexp.
const (
	aggregatorMatch = `(?P<aggregator>\w+)`
	metricMatch     = `(?P<metric>\S+)`
	comparatorMatch = `(?P<comparator>[<>=~!@#$%^&?]+)`
	thresholdMatch  = `(?P<threshold>[-+]?\d*\.?\d+)`
)

//nolint:gochecknoglobals
var constraintRegexp = utils.MustCompile(
	fmt.Sprintf(
		`^\s*%s\(%s\)\s*%s\s*%s\s*$`,
		aggregatorMatch, metricMatch, comparatorMatch, thresholdMatch,
	),
)

// ParseConstraint creates a constraint from a string representation.
func ParseConstraint(s string, parsers ...MetricParser) (*Constraint, error) {
	if !constraintRegexp.MatchString(s) {
		return nil, errors.New("invalid constraint format")
	}

	match := constraintRegexp.FindStringSubmatchMap(s)

	aggregator, err := ParseAggregator(match["aggregator"])
	if err != nil {
		return nil, err
	}

	metric, err := parseMetric(match["metric"], parsers...)
	if err != nil {
		return nil, err
	}

	comparator, err := ParseComparator(match["comparator"])
	if err != nil {
		return nil, err
	}

	threshold, err := strconv.ParseFloat(match["threshold"], 64)
	if err != nil {
		return nil, err
	}

	return &Constraint{
		Metric:     metric,
		Aggregator: aggregator,
		Comparator: comparator,
		Threshold:  threshold,
	}, nil
}

func parseMetric(name string, parsers ...MetricParser) (Metric, error) {
	for _, parser := range parsers {
		metric, err := parser(name)
		if err == nil {
			return metric, nil
		} else if !errors.Is(err, ErrNotParsed) {
			return nil, err
		}
	}

	return nil, ErrNotParsed
}
