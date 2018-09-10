# Usage Guide

## Test types

FBender provides two different approaches to load testing. The first,
__Throughput__, gives the tester control over the throughput (QPS), but not over
the concurrency. The second, __Concurrency__, gives the tester control over the
concurrency, but not over the throughput. In addition FBender distinguishes
between __Fixed__ and __Constraints__ based tests.

### Fixed test
Fixed tests allow user to specify the exact values of the load tests. Let's say
we want to run a DNS load test for `${TARGET}` with 50 QPS, 100 QPS and 200 QPS.

```sh
fbender dns throughput fixed -t ${TARGET} 50 100 200
```

the same can be achieved by running

```sh
fbender dns throughput fixed -t ${TARGET} 50
fbender dns throughput fixed -t ${TARGET} 100
fbender dns throughput fixed -t ${TARGET} 200
```

### Constraints test

In constraints tests user specifies a list of _constraints_ a test must meet to
be considered successful. Consecutive test values are adjusted based on the
given _growth_. For example a throughput test starting at 20 QPS that will
increase QPS by 10 after every test as long as the errors average didn't exceed
5%. This test will stop as soon as the constraints are not met.

```sh
fbender dns throughput constraints -t ${TARGET} -c "AVG(errors) < 5" 20 -g +10
```

#### Defining constraints

The constraints may be specified as a comma separated list, as well as a
consecutive `-c`, `--constraints` flags. Let's say `C${i}` for `i = 1,..,4` are
constraints. All of the below commands are equivalent:

```sh
fbender dns throughput constraints -t ${TARGET} -c "C1,C2,C3,C4" 100
fbender dns throughput constraints -t ${TARGET} -c "C1,C2" -c "C3,C4" 100
fbender dns throughput constraints -t ${TARGET} -c "C1" -c "C2,C3" -c "C4" 100
```

#### Defining growth

A __growth__ (`-g, --growth`) is used to determine a value for the next test
after performing the constraints check.

* _linear growth_ `+value` will increase a test by a constant value after every
successful test and will stop immediately after first test failure
```sh
fbender dns throughput constraints -t ${TARGET} -g +100 100 -c ${CONSTRAINTS}
# Tests: 100, 200, 300, 400, ...
```
* *percentage growth* `%value` will increase a test by a constant percentage
after every successful test and will stop immediately after first test failure
```sh
fbender dns throughput constraints -t ${TARGET} -g %100 100 -c ${CONSTRAINTS}
# Tests: 100, 200, 400, 800, ...
```
* *exponential growth* `^precision` will double the test to find a first failure
and then perform a binary search up to a given precision
```sh
fbender dns throughput constraints -t ${TARGET} -g ^20 100 -c ${CONSTRAINTS}
# Tests: 100 (OK), 200 (OK), 400 (FAIL), 300 (OK), 350 (OK), 375 (FAIL), 362 (OK)
```

#### Checking constraints
Internally each constraint consists of a __metric__, an __aggregator__, a
__comparator__ and a  __threshold__. Metrics may follow different syntaxes
depending on their needs. Pseudocode for checking metric:

```python
datapoints := fetchMetric(testStart, testDuration)
value := aggregate(datapoints)
return compare(value, threshold)
```

#### Syntax
```
Constraint ::= <Aggregator>(<Metric>) <Cmp> <Threshold>
Aggregator ::= "MIN" | "MAX" | "AVG"
Metric     ::= <string>
Cmp        ::= "<" | ">"
Threshold  ::= <float>
```

Metrics are parsed by one of the metric parsers (see `ConstraintsValue` in
`cmd/common/flags.go`). By default FBender supports only [Basic Metrics](#basic-metrics).
Check how to add your own metric parsers in [Extending FBender](#extending-fbender)
guide.

#### Basic Metrics

Basic metrics use only data gathered during the test. The available metrics are:
* __errors__ - errors percentage (the aggregator doesn't matter when using
errors metric in a constraint as the only datapoint is the overall errors
percentage)
* __latency__ - the packet latency (this may take a lot of memory so use wisely)

```bash
fbender dns throughput constraints -t ${TARGET} -c "MAX(errors) < 10" 100
# Checks if the errors during test are less than 10% of all requests
fbender dns throughput constraints -t ${TARGET} -c "AVG(latency) < 20" 100
# Checks if the average latency is less than 20ms (use -u to change unit)
```

## Common flags

### Target (required)

Target is a __required__ flag (`-t, --target`) that specifies the test target.
Formats accepted by most of the commands are:

* `IPv4`, `IPv4:port`
* `IPv6`, `[IPv6]:port`
* `hostname`, `hostname:port`

There might be other target formats and they are usually explicitly stated
in the command documentation.

### Duration

Duration (`-d, --duration`) specifies a single test duration. A duration format
is a sequence of decimal numbers, each with optional fraction and a unit suffix,
such as _"300ms"_, _"1.5h"_ or _"2h45m"._ Valid time units are _"ns"_, _"us"_
(or _"µs"_), _"ms"_, _"s"_, _"m"_, _"h"_.

### Input

Commands use input to generate requests for the load test. Unless explicitly
stated in the command documentation one request is generated per line in the
input file, skipping lines with improper format (refer to the command
documentation for format accepted by a specific protocol). The generated
requests are then reused in a __round-robin__ manner. If input flag
(`-i, --input`) is not specified `facebender` will read the requests from the
standard input.

```sh
fbender -i input.txt
```

Is equivalent to

```sh
cat input.txt | facebender
```

### Output

FBender uses `stderr` output to display test current state. All important
information is printed to `stdout`. Test logs can be redirected using the output
flag (`-o, --output`). They can also be filtered (`-v, --verbosity`) based on
the message verbosity level. Following levels are available, both numbers and
literals are accepted (`-v info` and `-v 4` are equivalent):

* _panic/0_
* _fatal/1_
* _error/2_
* _warning/3_ - log when an __error response__ is received
* _info/4_ - log when a __successful response__ is received
* _debug/5_ - log when a __request__ is sent

#### Format
User can chose a desired format (`-f, --format`) from one of the following.
Please note that this changes the output format only for the test logs.

* _text_ - human readable colored format, useful for debugging
* _json_ - very powerful and can be useful if you need to process things
further. Each line of the output contains a json message of format

##### JSON format

* _elapsed_ - elapsed time in nanoseconds (this is always equal to `end - start`)
* _end_ - request end time as unix nano
* _error_ - error message, this field is only present on failed requests
* _level_ - message verbosity level
* _msg_ - logged message (Success/Fail)
* _response_ - response converted to json
* _start_ - request start time as unix nano
* _test_ - desired qps/concurrency (depending on the protocol)
* _time_ - message log time

Example json log line of a __successful request__

```json
{
    "elapsed": 656894,
    "end": 1532360929961968467,
    "level": "info",
    "msg": "Success",
    "response": "HnRt5qTgrXkAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=",
    "start": 1532360929961311573,
    "test": 500,
    "time": "2018-07-23T08:48:49-07:00"
}
```

Example json log line of a __failed request__

```json
{
    "elapsed": 1151382976,
    "end": 1532360945472001038,
    "error": "read udp [2401:db00:3020:705f:face:0:76:0]:38665->[2401:db00:3011:5059:face:0:61:0]:50001: i/o timeout",
    "level": "warning",
    "msg": "Fail",
    "response": null,
    "start": 1532360944320618062,
    "test": 500,
    "time": "2018-07-23T08:49:05-07:00"
}
```

### Timeout

Timeout allows to set a wait timeout (`-w, --timeout`) for a single request.
Please note that this option __can be used differently__ across different
protocols. For example TFTP treats it as a timeout for a single datagram read
not for a whole file transfer. When generating statistics output the values are
clamped to `[0, timeout * 2]` range.

### Distribution

Although overall queries per second amount is constant the packets may be send
based on one of the following distributions (`-D, --distribution`)

* *uniform* - a request will be sent every `1/QPS` seconds
* *exponential* - a request will be sent based on a Poisson process, where
desired QPS corresponds to the reciprocal of the lambda parameter to an
exponential distribution.

### Statistics

To provide a short, useful output right away after the test finishes FBender
collects data in a __histogram__. However it might take a lot of memory,
multiple buckets need to be created for every `unit` in a `[0, timeout * 2]`
range. Both timeout (`-w, --timeout`) and unit (`-u, --unit`) may be customized.
Their format is a sequence of decimal numbers, each with optional fraction and a
unit suffix, such as _"300ms"_, _"1.5h"_ or _"2h45m"._ Valid time units are
_"ns"_, _"us"_ (or _"µs"_), _"ms"_, _"s"_, _"m"_, _"h"_. When memory limit is an
issue the statistics can be __disabled completely__ with the `--nostats` flag
and the JSON log output can be used later to generate them on a different
machine.

### Buffer

FBender internally uses buffers to generate the requests and process them.
Although default buffer size should be suitable for most of the standard tests,
increasing the buffer may result in better performance. Experiment with
different buffer sizes (`-b, --buffer`) if you find FBender clumsy and not
generating enough requests. Check out [Bender performance](https://github.com/pinterest/bender#performance)
for more performance hacks.

## Bash completion

### Requirements

Before attempting to run the command make sure you have a bash completion
installed and enabled. In CentOS for example you need to install:

```sh
sudo yum install bash-completion bash-completion-extras
```

And source the `bash_completion` file:

```sh
source /etc/profile.d/bash_completion.sh
```

We recommend adding the above line to your `.bashrc`.

### Enable bash completion

To enable fbender autocompletion in bash run (you may want to add this line
to your `.bashrc` to automatically run it when you open a new shell):

```sh
source <(fbender completion bash)
```

## Troubleshooting

### FBender displays help instead of running a test

Make sure you've specified the protocol, test type and all required flags and
arguments. Documentation show many proper usage example, which you may
__copy and adjust__ to your needs.

### Out of memory

Try __adjusting unit/timeout__ to match your needs and consider __disabling
statistics__. Refer to statistics documentation for more details. You may also
try decreasing the buffer size. In the worst case simply __pick a more powerful
machine__ and run the tests from a different host. Additional help may be found
at [Bender performance](https://github.com/pinterest/bender#performance)
