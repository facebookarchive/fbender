# FBender
[![Build Status](https://travis-ci.org/facebookincubator/fbender.svg?branch=master)](https://travis-ci.org/facebookincubator/fbender)
[![codecov](https://codecov.io/gh/facebookincubator/fbender/branch/master/graph/badge.svg)](https://codecov.io/gh/facebookincubator/fbender)
[![Go Report Card](https://goreportcard.com/badge/github.com/facebookincubator/fbender)](https://goreportcard.com/report/github.com/facebookincubator/fbender)

FBender is a __load testing__ command line tool for generic network protocols.

As a foundation for load testing lays the [Pinterest Bender](https://github.com/pinterest/bender)
library. Similar to Bender, FBender provides two different approaches to load
testing. The first, __Throughput__, gives the tester control over the throughput
(QPS), but not over the concurrency. The second, __Concurrency__, gives the
tester control over the concurrency, but not over the throughput. You can read
more about that in the [Bender documentation](https://github.com/pinterest/bender#bender).

FBender has been designed to be easily extendable by additional protocols. Look
at the guide on how to contribute new protocols.

## Examples

In the _first example_ we will be load testing a __DNS__ server __example.com__
running on the __default port__ (53). We will perform 3 consecutive tests for
each __specified QPS__ (2000, 4000, 8000) each lasting for __1 minute__. The
queries will be generated based on the input file __queries.txt__. We will
ignore requests output.

```sh
fbender dns throughput fixed \
  --target example.com --duration 1m \
  --input queries.txt -v error \
  2000 4000 8000
```

In the _next example_ we will be load testing a __TFTP__ server __example.com__
running on the __default port__ (69). We will perform 3 consecutive tests for
each __specified number of concurrent connections__ (10, 25, 50) each lasting
for __1 minute__. The queries will be generated based on the input file
__files.txt__. We will ignore requests output.

```sh
fbender tftp concurrency fixed \
  --target example.com --duration 1m \
  --input files.txt -v error \
  10 25 50
```

The _last example_ will focus on finding the SLA for a __DHCPv6__ server
__example.com__. We want the timeouts not to exceed __5% of all requests__ in
the measure window of __1 minute__. To get the most accurate results we will be
using exponential backoff growth starting at 20 QPS with a precision of 10 QPS.
The queries will be generated based on the input file __macs.txt__. We will
ignore requests output.

```sh
fbender dhcpv6 throughput constraints \
  --target example.com --duration 1m \
  --input macs.txt -v error \
  --constraints "AVG(errors) < 5" \
  --growth ^10 20
```

## Building FBender

```sh
go get -u github.com/facebookincubator/fbender
go build github.com/facebookincubator/fbender
```

## Installing FBender

```sh
go get -u github.com/facebookincubator/fbender
go install github.com/facebookincubator/fbender
```

You may want to add the following line to your .bashrc to enable autocompletion
```sh
source <(fbender complete bash)
```

## Docs

* [General usage guide](https://github.com/facebookincubator/fbender/blob/master/docs/USAGE.md)
* [Protocol specific usage and examples](https://github.com/facebookincubator/fbender/blob/master/docs/PROTOCOLS.md)
* [Extending FBender](https://github.com/facebookincubator/fbender/blob/master/docs/EXTENDING.md)

## License

FBender is BSD licensed, as found in the LICENSE file.
