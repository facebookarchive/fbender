# Protocols

## DHCPv6

Due to DHCPv6 nature we can only test relay. This means taking ownership of
port 547 on the machine we're running the load test, which requires __root
access__. DHCPv6 load test performs a __4-way handshake__ (Solicit, Advertise,
Request, Reply) and results in an error if any of these fail.

### Input Format

```
DeviceMAC
```

#### Example input

```
01:23:45:67:89:ab
E3:63:BD:7B:D2:2C
c8:6c:2c:47:96:fd
```

### Examples

In the _first example_ we will be load-testing a DHCPv6 (relay) server running on
__example.com__ on the default __port 547__. We will perform 3 consecutive tests
for each specified QPS (__200__, __400__, __800__) each lasting for __1 minute__.
The solicits will be generated based on the input file __macs.txt__. The results
will be stored as __json__ in __results.json__.

```sh
fbender dhcpv6 throughput fixed \
    --target example.com --duration 1m \
    --input macs.txt --output results.json --format json \
    200 400 800
```

In the _second example_ we will try to automatically find a maximum QPS a server
can handle before it's __average latency exceeds 50ms__.  We will use the
same server (__example.com:547__), starting our test at __400__ increasing by
__100__ with every test. To get more measure points we will run tests for
__10 minutes__ and check for __average latency__ not to exceed 50ms. Just like
previously we will generate solicits from __macs.txt__ and save the results as
json to __results.json__.

```sh
fbender dhcpv6 throughput constraints \
    --target example.com --duration 10m \
    --input macs.txt --output results.json --format json \
    400 --growth +100 --constraints "AVG(latency) < 50"
```

## DNS

### Custom Flags

In addition to other standard flags DNS allows you prefix all queries with 
randomly generated values to avoid cached responses. When (`-r, --randomize` ) 
the queries will be prefixed as follows (the _QType_ remains unchanged)

```
time.hex.domain
```

In bash this could have been achieved with the following command:

```sh
$(date +%s).$(openssl rand -hex 16).domain
```

### Input Format

```
Domain QType
```

#### Example input

```
example.com AAAA
another.example.com A
yet.another.example.com MX
example.com TXT
```

Also, DNS can be load tested over tcp, rather than just the standard udp
interface, by specifying (`-p, --protocol tcp`).

### Examples

In the _first example_ we will be load-testing a  server running on
__example.com__ on the default __port 53__. We will perform 3 consecutive tests
for each specified QPS (__2000__, __4000__, __8000__) each lasting for
__1 minute__. The queries will be generated based on the input file
__queries.txt__, and to avoid hitting the cache we will __randomize__ them with
a random hex of 16-character length. The results will be stored as __json__ in
__results.json__.

```sh
fbender throughput fixed \
    --target example.com--duration 1m --randomize 16 \
    --input queries.txt --output results.json --format json \
    2000 4000 8000
```

In the _second example_ we will try to automatically find a maximum QPS a server
can handle before its __average latency exceeds 20ms__.  We will use the
same server (__example.com:53__), starting our test at __4000__ increasing by
__1000__ with every test. To get more measure points we will run tests for
__10 minutes__ and only take into an account the __minimum data point__. Just
like previously we will generate queries from __queries.txt__, __randomize__
them and save the results as json to __results.json__.

```sh
fbender  throughput constraints \
    --target example.com:53 --duration 10m --randomize\
    --input queries.txt --output results.json --format json \
    4000 --growth +1000 --constraints "AVG(latency) < 20"
```

## HTTP

### Custom Flags

In addition to standard flags HTTP provides a flag to enable ssl (`-s, --ssl`).
If enabled all requests will be sent over HTTPS instead of default HTTP.

### Input Format

Depending on the method lines may have one of two formats. When generating
requests a protocol, target and url will be joined together removing excess `/`.
Let's say the target is __example.com__, the __ssl__ flag is turned on and the
input is __GET /index.html__ than the actual request will be
`GET https://example.com/index.html.`

```
GET RelativeURL
POST RelativeURL FormData
```

#### Example input

```
GET index.html
GET /
POST echo message=Hello
POST hello/ lang=en&name=Mikolaj
```

### Examples

In the _first example_ we will be load-testing a HTTP server running on
__example.com__ on the default __port 443__ (with __SSL__ enabled). We will
perform 3 consecutive tests for each specified QPS (__2000__, __4000__, __8000__)
each lasting for __1 minute__. The queries will be generated based on the input
file __queries.txt__. The results will be stored as __json__ in __results.json__.

```sh
fbender http throughput fixed \
    --target example.com --duration 1m --ssl \
    --input queries.txt --output results.json --format json \
    2000 4000 8000
```

In the _second example_ we will try to automatically find a maximum QPS a server
can handle before it's __average latency exceeds 200ms__.  We will use the same
server (__example.com:443__), starting our test at __400__ increasing by
__100__ with every test. To get more measure points we will run tests for
__10 minutes__ and only take into an account the __minimum data point__. Just
like previously we will generate queries from __queries.txt__ and save the
results as json to __results.json__.

```sh
fbender http throughput constraints \
    --target example.com --duration 10m --ssl \
    --input queries.txt --output results.json --format json \
    400 --growth +100 --constraints "AVG(latency) < 500.0"
```

The _next example_ we will be load-testing endpoint which downloads a large file
(around 5GB), so instead of using throughput test we will test concurrency. We
will use the server __example.com:443__. We will perform 3 consecutive tests for
each specified number of concurrent connections (__10__, __20__, __50__) each
lasting for __1 minute__. The queries will be generated based on the input file
__queries.txt__. The results will be stored as __json__ in __results.json__.

```sh
fbender http concurrency fixed \
    --target example.com --duration 1m --ssl \
    --input queries.txt --output results.json --format json \
    10 20 50
```

In the _last example_ we will try to automatically find a maximum concurrent
connections a server can handle before it's __errors rate exceeds 10%__.  We
will use the same server (__example.com:443__), starting our test at __20__ and
growth of __exponential backoff__ with precision of 10. To get more measure
points we will run tests for __10 minutes__ and only take into an account the.
Just like previously we will generate queries from __queries.txt__ and save the
results as json to __results.json__.

```sh
fbender http concurrency constraints \
    --target example.com --duration 10m --ssl \
    --input queries.txt --output results.json --format json \
    20 --growth ^5 --constraints "MAX(errors) < 10.0"
```

## TFTP

Due to a large amount of packets required to download a single file over TFTP
instead of queries per second we will test number of concurrent connections.
The specified __timeout__ applies to a single datagram in a tftp transfer rather
than to the whole session.

### Custom Flags
In addition to standard flags TFTP provides a flag to set the tftp block size
(`-s`, `--blocksize`).

### Input Format

Mode can be one of `octet` or `netascii`

```
Filename Mode
```

#### Example input

```
/my/file octet
/my/otherfile octet
/anotherfile netascii
```

### Examples

In the _first example_ we will be load-testing a TFTP server running on
__example.com__ on the default __port 69__. We will perform 3 consecutive tests
for each specified concurrent connections (__10__, __25__, __50__) each lasting
for __1 minute__. The queries will be generated based on the input file
__files.txt__. The results will be stored as __json__ in __results.json__.

```sh
fbender tftp concurrency fixed \
    --target example.com:69 --duration 1m \
    --input files.txt --output results.json --format json \
    10 25 50
```

In the _second example_ we will try to automatically find a maximum QPS a server
can handle before it's __errors rate exceeds 10%__.  We will use the same server
(__example.com:69__), starting our test at __25__ increasing by __10__ with
every test. To get more measure points we will run tests for __10 minutes__
and only take into an account the __minimum data point__. Just like previously
we will generate queries from __files.txt__ and save the results as json to
__results.json__.

```sh
fbender tftp concurrency constraints \
    --target example.com --duration 10m \
    --input files.txt --output results.json --format json \
    25 --growth +10 --constraints "MAX(errors) < 10.0"
```

## UDP

### Target

Target format accepted by the udp tests is `ipv4`, `ipv6`, `hostname`. The ports
are taken from the input file.

### Input Format

For udp load-testing the datagram content has to be base64 encoded to allow for
easy line by line parsing.

```
DestinationPort Base64EncodedPayload
```

#### Example input

```
2545 TG9yZW0=
7346 aXBzdW0gZG9sb3Igc2l0
5012 YW1ldCBpbg==
```

### Examples

In the _first example_ we will be load-testing a server running on
__example.com__. We will perform 3 consecutive tests for each specified QPS
(__2000__, __4000__, __8000__) each lasting for __1 minute__. The queries will
be generated based on the input file __payloads.txt__. The results will be
stored as __json__ in __results.json__.

```sh
fbender udp throughput fixed \
    --target example.com --duration 1m \
    --input payloads.txt --output results.json --format json \
    2000 4000 8000
```

In the _second example_ we will try to automatically find a maximum QPS a server
can handle before it's __errors rate exceeds 10%__.  We will use the
same server (__example.com__), starting our test at __4000__ increasing by
__1000__ with every test. To get more measure points we will run tests for
__10 minutes__ and only take into an account the __minimum data point__. Just
like previously we will generate requests from __payloads.txt__ and save the
results as json to __results.json__.

```sh
fbender udp throughput constraints \
    --target example.com --duration 10m \
    --input payloads.txt --output results.json --format json \
    4000 --growth +1000 --constraints "MAX(errors) < 10.0"
```
