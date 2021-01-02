/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package flags

import (
	"fmt"
	"sort"
	"strings"

	"github.com/facebookincubator/fbender/utils"
	"github.com/pinterest/bender"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// DistributionGenerator represents distribution generator function.
type DistributionGenerator = func(float64) bender.IntervalGenerator

const (
	uniformGenerator     = "uniform"
	exponentialGenerator = "exponential"
)

//nolint:gochecknoglobals
var generators = map[string]DistributionGenerator{
	uniformGenerator:     bender.UniformIntervalGenerator,
	exponentialGenerator: bender.ExponentialIntervalGenerator,
}

// Distribution represents a interval generator flag value.
type Distribution struct {
	Name      string
	generator DistributionGenerator
}

// NewDefaultDistribution returns new distribution flag with default values.
func NewDefaultDistribution() *Distribution {
	return &Distribution{
		Name:      uniformGenerator,
		generator: generators[uniformGenerator],
	}
}

// DistributionChoices returns a string representation of available generators.
func DistributionChoices() []string {
	choices := []string{}

	for key := range generators {
		choices = append(choices, key)
	}

	sort.Strings(choices)

	return choices
}

func (d *Distribution) String() string {
	return d.Name
}

// Set validates a given value and sets distribution (allows prefix matching).
func (d *Distribution) Set(value string) error {
	matches := []string{}

	for key := range generators {
		if strings.HasPrefix(key, value) {
			matches = append(matches, key)
		}
	}

	if len(matches) == 0 {
		choices := ChoicesString(DistributionChoices())

		return fmt.Errorf("generator must be one of %s, '%s' given", choices, value)
	} else if len(matches) > 1 {
		return fmt.Errorf("ambiguous generator '%s' matches %s", value, ChoicesString(matches))
	}

	generator := matches[0]
	d.Name = generator
	d.generator = generators[generator]

	return nil
}

// Type returns a distribution type.
func (d *Distribution) Type() string {
	return "distribution"
}

// Get returns a distibution generator.
func (d *Distribution) Get() DistributionGenerator {
	return d.generator
}

// GetDistribution returns a distribution from a pflag set.
func GetDistribution(f *pflag.FlagSet, name string) (DistributionGenerator, error) {
	flag := f.Lookup(name)
	if flag == nil {
		return nil, fmt.Errorf("flag %s accessed but not defined", name)
	}

	return GetDistributionValue(flag.Value)
}

// GetDistributionValue returns a distribution from a pflag value.
func GetDistributionValue(v pflag.Value) (DistributionGenerator, error) {
	if distribution, ok := v.(*Distribution); ok {
		return distribution.Get(), nil
	}

	return nil, fmt.Errorf("trying to get distribution value of flag of type %s", v.Type())
}

// Bash completion function constants.
const (
	fnameDistribution = "__fbender_handle_distribution_flag"
	fbodyDistribution = `COMPREPLY=($(compgen -W "uniform exponential" -- "${cur}"))`
)

// BashCompletionDistribution adds bash completion to a distribution flag.
func BashCompletionDistribution(cmd *cobra.Command, f *pflag.FlagSet, name string) error {
	flag := f.Lookup(name)
	if flag == nil {
		return fmt.Errorf("flag %s accessed but not defined", name)
	}

	if _, ok := flag.Value.(*Distribution); !ok {
		return fmt.Errorf("trying to autocomplete distribution on flag of type %s", flag.Value.Type())
	}

	return utils.BashCompletion(cmd, f, name, fnameDistribution, fbodyDistribution)
}
