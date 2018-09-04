# Extending FBender

## Adding new protocol

### Implement protocol executor

If the protocol executor implementation can be used other than in FBender try
contributing to [Bender](https://github.com/pinterest/bender) and use the
upstreamed version in the tester. If you think that it's not beneficial to send
a Pull Request with it to Bender create a subdirectory in `protocols` directory.

### Create tester

Create a subdirectory in `tester` with the name of your protocol and a single
`tester.go` file. The tester file should contain a definition of `struct Tester`
and optional definition of `interface Options`. Tester needs to implement:

```go
// Tester is used to setup the test for a specific endpoint.
type Tester interface {
	// Before is called once, before any tests.
	Before(options interface{}) error
	// After is called once, after all tests (or after some of them if a test fails).
	// This should be used to cleanup everything that was set up in the Before.
	After(options interface{})
	// BeforeEach is called before every test.
	BeforeEach(options interface{}) error
	// AfterEach is called after every test, even if the test fails. This should
	// be used to cleanup everything that was set up in the BeforeEach.
	AfterEach(options interface{})
	// RequestExecutor is called every time a test is to be ran to get an executor.
	RequestExecutor(options interface{}) (bender.RequestExecutor, error)
}
```

For compatibility with different testers options are passed as an `interface{}`,
when asserting options type your tester should return `tester.ErrInvalidOptions`
in case of failure.

### Create a command

Create a subdirectory in `cmd` with the name of your protocol and a files named
`cmd.go` and `${PROTOCOL}.go`. File `cmd.go` should contain a definition of a
command and an optional `init` function adding additional flags or subcommands.
You can copy the following template and replace `${VARIABLES}` with appropriate
values.

```go
var template = &core.CommandTemplate{
	Name: "${PROTOCOL}",
	Short: "Test ${PROTOCOL}",
	Long: `
Input format: "${INPUT_FORMAT}"
  ${INPUT_EXAMPLE_1}
  ${INPUT_EXAMPLE_2}`,
	Fixed: `  fbender ${PROTOCOL} {test} fixed -t $TARGET 10 20
  ${ANOTHER_FIXED_TEST_EXAMPLE}`,
	Constraints: `  fbender ${PROTOCOL} {test} constraints -t $TARGET -c "AVG(latency)<10" 20
 ${ANOTHER_CONSTRAINTS_TEST_EXAMPLE}`,
}

var Command = core.NewTestCommand(template, params)
```

The `${PROTOCOL}.go` should implement all protocol specific features. The
simplest one could look like this:

```go
func params(cmd *cobra.Command, o *options.Options) (*runner.Params, error) {
	// create input based request generator
	requests, err := input.NewRequestGenerator(o.Input, inputTransformer)
	if err != nil {
		return nil, err
	}
	// create tester for your protocol
	tester := &protocol.Tester{
		Target: o.Target,
	}
	return &runner.Params{Tester: tester, RequestGenerator: requests}, nil
}

// inputTransformer accepts any string as an input and returns it as a request
func inputTransformer(input string) (interface{}, error) {
	return input, nil
}
```

Take a look at already implemented protocols for better overview.

### Register a subcommand

Add import line and a subcommand to the main `cmd/cmd.go` file.

```go
import (
	// ...
	"github.com/facebookincubator/fbender/cmd/dhcpv4"
	"github.com/facebookincubator/fbender/cmd/dhcpv6"
	"github.com/facebookincubator/fbender/cmd/dns"
	"github.com/facebookincubator/fbender/cmd/http"
	// ... , add you command package here (alphabetically)
	"github.com/facebookincubator/fbender/cmd/tftp"
	"github.com/facebookincubator/fbender/cmd/udp"
)

var Subcommands = []*cobra.Command{
	dhcpv4.Command,
	dhcpv6.Command,
	dns.Command,
	http.Command,
	// ... , add you command here (alphabetically)
	tftp.Command,
	udp.Command,
}
```

## Adding internal features

To adjust FBender to your needs without the necessity to create your own forks
of the repository and always keep up to date with the newest version we
recommend creating custom `main` function and "patching" your features into the
main command. For example at Facebook we use it to add load tests for internal
services and integrate our metric system for constraints tests.

```go
package main

import (
	"github.com/facebookincubator/fbender/cmd"
	"github.com/facebookincubator/fbender/cmd/core"

	"fbender/internal/cmd"
	"fbender/internal/metric"
)

func main() {
		// Add internal commands
		cmd.Command.AddCommand(cmd.Command)
		// Patch internal metrics
		core.ConstraintsValue.Parsers = append(core.ConstraintsValue.Parsers, metric.MetricParser)
		// Execute command
		cmd.Execute()
}
```
