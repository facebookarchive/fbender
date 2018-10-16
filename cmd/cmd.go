/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/facebookincubator/fbender/cmd/core"
	"github.com/facebookincubator/fbender/cmd/dhcpv4"
	"github.com/facebookincubator/fbender/cmd/dhcpv6"
	"github.com/facebookincubator/fbender/cmd/dns"
	"github.com/facebookincubator/fbender/cmd/http"
	"github.com/facebookincubator/fbender/cmd/tftp"
	"github.com/facebookincubator/fbender/cmd/udp"
	"github.com/facebookincubator/fbender/flags"
)

// Subcommands are the protocol subcommands
var Subcommands = []*cobra.Command{
	dhcpv4.Command,
	dhcpv6.Command,
	dns.Command,
	http.Command,
	tftp.Command,
	udp.Command,
}

// Command is the root command for the CLI
var Command = &cobra.Command{
	Use: "fbender",
	Long: `FBender is a load tester tool for various protocols. It provides two different
approaches to load testing: Throughput and Concurrency and each of them can have
either fixed or constraints based test values. Throughput tests give the tester
control over the throughput (QPS), but not over concurrency. The second gives
the user control over the concurrency but not over the throughput.
  * fixed - runs a single test for each of the specified values.
  * constraint - runs tests adjusting load based on the growth and constraints.

Target:
Target format may vary depending on the protocol, however most of them accept
ipv4, ipv6, hostname with an optional port. Use "fbender protocol --help" to get
the documentation on the target format for a specific protocol.

Input:
Unless explicitly stated in the command documentation one request is generated
per input line, skipping the lines with improper format. Use "fbender
protocol help" to get the documentation on the input format for a specific
protocol. The generated requests are then reused in a round-robin manner.

Output:
All important information is printed to stdout. Test logs can be redirected
using the output flag. They can also be filtered based on the message verbosity
level. Note that this filters/redirect only test logs and not the summary and
other output. Available levels (both numbers and literals are accepted):
  * panic/0
  * fatal/1
  * error/2
  * warning/3 - log when an *error response* is received
  * info/4 - log when a *successful response* is received
  * debug/5 - log when a *request* is sent
`,
	Example: `  fbender dns throughput fixed -t $TARGET 100
  fbender tftp concurrency fixed -t $TARGET -o /dev/null 10
  fbender udp throughput fixed -t $TARGET -d 5m 100 200 300
  fbender http concurrency constraints -t $TARGET 20 -c "MAX(errors)<5"
  fbender dhcpv6 throughput constraints -t $TARGET 50 -c "MIN(latency)<20"
  fbender dns throughput constraints -t $TARGET 40 -c -g ^10 "MAX(errors)<5"`,
	BashCompletionFunction: `
	__fbender_handle_loglevel_flag()
	{
		COMPREPLY=($(compgen -W "panic fatal error warning info debug" -- "${cur}"))
	}

	__fbender_handle_logformat_flag()
	{
		COMPREPLY=($(compgen -W "text json" -- "${cur}"))
	}`,
}

func initIOFlags() {
	// Input
	Command.PersistentFlags().StringP("input", "i", "", "load test input data from a file (default <stdin>)")
	if err := Command.MarkPersistentFlagFilename("input"); err != nil {
		panic(err)
	}
	// Output
	logOutput := flags.NewLogOutput(logrus.StandardLogger())
	Command.PersistentFlags().VarP(logOutput, "output", "o", "log test output to a file")
	if err := Command.MarkPersistentFlagFilename("output"); err != nil {
		panic(err)
	}
	// Log Level
	logLevel := &flags.LogLevel{Logger: logrus.StandardLogger()}
	logLevelChoices := flags.ChoicesString(flags.LogLevelChoices())
	Command.PersistentFlags().VarP(logLevel, "verbosity", "v", fmt.Sprintf("verbosity level %s", logLevelChoices))
	if err := cobra.MarkFlagCustom(Command.PersistentFlags(), "verbosity", "__fbender_handle_loglevel_flag"); err != nil {
		panic(err)
	}
	// Log format
	logFormat := &flags.LogFormat{Logger: logrus.StandardLogger(), Format: "json"}
	logFormatChoices := flags.ChoicesString(flags.LogFormatChoices())
	Command.PersistentFlags().VarP(logFormat, "format", "f", fmt.Sprintf("output format %s", logFormatChoices))
	if err := cobra.MarkFlagCustom(Command.PersistentFlags(), "format", "__fbender_handle_logformat_flag"); err != nil {
		panic(err)
	}
}

func initExecutionFlags() {
	// Test duration
	Command.PersistentFlags().DurationP("duration", "d", 1*time.Minute, "single test duration")
	// Requests distribution
	distribution := flags.NewDefaultDistribution()
	distributionChoices := flags.ChoicesString(flags.DistributionChoices())
	Command.PersistentFlags().VarP(distribution, "dist", "D", fmt.Sprintf("requests distribution %s", distributionChoices))
	if err := flags.BashCompletionDistribution(Command, Command.PersistentFlags(), "dist"); err != nil {
		panic(err)
	}
	// Other settings
	Command.PersistentFlags().IntP("buffer", "b", 2048, "buffer size of the requests generator channel")
	Command.PersistentFlags().DurationP("timeout", "w", 1*time.Second, "wait timeout on requests")
	Command.PersistentFlags().DurationP("unit", "u", 1*time.Millisecond, "histogram scaling unit")
	Command.PersistentFlags().Bool("nostats", false, "disable statistics")
}

func init() {
	cobra.EnablePrefixMatching = true
	initIOFlags()
	initExecutionFlags()
	for _, subcommand := range Subcommands {
		Command.AddCommand(subcommand)
		subcommand.PersistentFlags().StringP("target", "t", "", "endpoint to load test")
		if err := subcommand.MarkPersistentFlagRequired("target"); err != nil {
			panic(err)
		}
	}
	Command.AddCommand(completionCmd)
	// Start post init functions
	core.PostInit <- struct{}{}
	core.PostInitWaitGroup.Wait()
}

// Execute runs the Command
func Execute() {
	if err := Command.Execute(); err != nil {
		os.Exit(1)
	}
}
