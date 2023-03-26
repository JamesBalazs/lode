# Lode

[![Coverage Status](https://coveralls.io/repos/github/JamesBalazs/lode/badge.svg?branch=main&t=LIyVhQ&service=github)](https://coveralls.io/github/JamesBalazs/lode?branch=main)
[![Linux](https://svgshare.com/i/Zhy.svg)](https://svgshare.com/i/Zhy.svg)
[![macOS](https://svgshare.com/i/ZjP.svg)](https://svgshare.com/i/ZjP.svg)
[![Windows](https://svgshare.com/i/ZhY.svg)](https://svgshare.com/i/ZhY.svg)
[![made-with-Go](https://img.shields.io/badge/Made%20with-Go-1f425f.svg)](http://golang.org)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg?style=flat-square)](http://makeapullrequest.com)

Versatile load testing CLI tool written in Go, with configurable workflows to facilitate automated load testing in CI.

![Lode CLI tool example image](https://i.imgur.com/5p0CE4F.png)

**Features:**
- Concurrent
- Timing data for individual requests
- Store and replay load test results and responses
- Configurable with workflows for CI/automation use
- Open source

## Installation
Check out our [installation instructions](https://github.com/JamesBalazs/lode/wiki/Installation) for supported platforms in the Wiki.

Alternatively, head to [releases](https://github.com/JamesBalazs/lode/releases), grab the version for your platform, extract the binary, and run it. 

## Usage
### `lode test [flags] [url]`
Used to run a single load test against the specified URL.
A summary report is printed at the end.

**Supported flags:**
| Flag | Shorthand | Usage |
| --- | --- | --- |
| `--freq` | `-f` | Number of requests to make per second |
| `--delay` | `-d` | Time to wait between requests, e.g. 200ms or 1s - defaults to 1s unless --freq specified |
| `--concurrency` | `-c` | Maximum number of concurrent requests |
| `--maxRequests` | `-n` | Maximum number of requests to make - defaults to 0s (unlimited) |
| `--maxTime` | `-l` | Length of time to make requests, e.g. 20s or 1h - defaults to 0s (unlimited) |
| `--method` | `-m` | HTTP method to use - defaults to GET |
| `--timeout` | `-t` | Timeout per request, e.g. 200ms or 1s - defaults to 5s |
| `--body` | `-b` | POST/PUT body |
| `--file` | `-F` | POST/PUT body filepath |
| `--header` | `-H` | Request headers, in the form X-SomeHeader=value - separate headers with commas, or repeat the flag to add multiple headers |
| `--interactive` | `-i` | Use interactive mode, which presents a scrollable list of requests, and shows the timing, body, and headers, of the selected request |
| `--fail-fast` |  | Abort the test immediately if a non-success status code is received |
| `--ignore-failures` |  | Don't return non-zero exit code when non-success status codes are received |
| `--out` | `-O` | Filepath to write requests and timing data, if provided |
| `--outFormat` |  | Format to use when writing requests to file - valid options are `json` and `yaml`, defaults to `json` |

One of either `--delay` or `--freq` is required. If both are provided, delay will be calculated from the given frequency.

**Examples:**
- `lode test -f 20 -c 4 -l 10s http://www.google.com` make 20 req/sec to Google for 10 seconds, split across 4 threads
- `lode test -d 1h -n 24 http://www.google.com` make 1 req/hr to Google until 24 requests have been made
- `lode test -f 40 -c 8 -l 1m -n 1000 http://www.google.copm` make 40 req/sec to Google, split across 8 threads, for up to 1 minutes or until 1000 requests have been made (whichever comes first)

## Example output
```
❯ lode test --freq=20 -c 8 -l 5s http://my.example.service
Target: GET http://my.example.service
Concurrency: 8
Requests made: 100
Time taken: 5.04s
Requests per second (avg): 19.84

Response breakdown:
200: ===================>  98x
501: =>                    2x

Percentile latency breakdown:
50th: 90ms
66th: 95ms
75th: 100ms
80th: 104ms
90th: 114ms
95th: 130ms
98th: 171ms
99th: 221ms
100th: 239ms
```

### `lode time [flags] [path]`
Used to run a single request.
A breakdown of the request's timing is printed at the end.

**Supported flags:**
| Flag | Shorthand | Usage |
| --- | --- | --- |
| `--method` | `-m` | HTTP method to use - defaults to GET |
| `--timeout` | `-t` | Timeout per request, e.g. 200ms or 1s - defaults to 5s |
| `--body` | `-b` | POST/PUT body |
| `--file` | `-F` | POST/PUT body filepath |
| `--header` | `-H` | Request headers, in the form X-SomeHeader=value - separate headers with commas, or repeat the flag to add multiple headers |
| `--interactive` | `-i` | Use interactive mode, which shows the timing, body, and headers, of the request |
| `--ignore-failures` |  | Don't return non-zero exit code when non-success status codes are received |
| `--out` | `-O` | Filepath to write requests and timing data, if provided |
| `--outFormat` |  | Format to use when writing requests to file - valid options are `json` and `yaml`, defaults to `json` |

**Example:**

`lode time http://www.google.com` make 1 request to Google and print timings

**Example output:**
```
❯ lode time https://www.google.com/
Target: GET https://www.google.com/
Concurrency: 1
Requests made: 1
Time taken: 290ms
Requests per second (avg): 3.36

Timing breakdown:
<=>             DNS Lookup:        24ms
   <=>          TCP Connection:    21ms
      <=>       TLS Handshake:     184ms
         <=>    Server:            66ms
            <=> Response Transfer: 0s
<=============> Total:             296ms
```

### `lode suite [flags] [path]`
Used to run a sequence of load tests, as defined in the YAML file at the specified path.

**YAML format:**
```
tests:
  - url: https://www.google.co.uk
    method: GET
    concurrency: 4
    freq: 10
    maxrequests: 20
  - url: https://abc.xyz/
    method: GET
    concurrency: 2
    delay: 0.5s
    maxrequests: 4
    headers:
      - SomeHeader=someValue
      - OtherHeader=otherValue
```

**Example:**

`lode suite examples/suite.yaml`

**Supported keys:**
| Flag | Usage |
| --- | --- |
| `url` | URL to target |
| `freq` | Number of requests to make per second |
| `delay` | Time to wait between requests, e.g. 200ms or 1s - defaults to 1s unless --freq specified |
| `concurrency` | Maximum number of concurrent requests |
| `maxrequests` | Maximum number of requests to make - defaults to 0s (unlimited) |
| `maxtime` | Length of time to make requests, e.g. 20s or 1h - defaults to 0s (unlimited) |
| `method` | HTTP method to use - defaults to GET |
| `timeout` | Timeout per request, e.g. 200ms or 1s - defaults to 5s |
| `body` | POST/PUT body |
| `file` | POST/PUT body filepath |
| `header` | Array of request headers, in the form X-SomeHeader=value |
| `failfast` | Boolean - Abort the test immediately if a non-success status code is received |
| `ignorefailures` | Boolean - Don't return non-zero exit code when non-success status codes are received |

## Usage
### `lode replay [flags] [filepath]`
Used to load the report of a single load test from the specified file.

**Supported flags:**
| Flag | Shorthand | Usage |
| --- | --- | --- |
| `--inFormat` |  | Format of log file - valid options are `json` and `yaml`, defaults to `json` |

**Examples:**
- `lode replay ./out.json` load the log file out.json and replay the interactive report from that run
- `lode replay --inFormat yaml ./out.yaml` load the log file out.yaml, as yaml, and replay the interactive report from that run

## Example output
```
❯ lode replay ./out.json
Target: GET http://my.example.service
Concurrency: 8
Requests made: 100
Time taken: 5.04s
Requests per second (avg): 19.84

Response breakdown:
200: ===================>  98x
501: =>                    2x

Percentile latency breakdown:
50th: 90ms
66th: 95ms
75th: 100ms
80th: 104ms
90th: 114ms
95th: 130ms
98th: 171ms
99th: 221ms
100th: 239ms

[interactive report]...
```

## Planned Features
- Timing/response code assertions for CI use

## Releasing
Releases are built for multiple platforms using [goreleaser](https://github.com/goreleaser/goreleaser) in GitHub Actions.
Simply add a tag starting with `v` and the Action will release the version.

```
git checkout main
git tag -a "v1.2.3"
git push --tags
```

If you need to install locally, you need a `GITHUB_TOKEN` envar with `repo` scope. Then you can tag, push, build and release by running:
```
git tag -a "v1.2.3"
goreleaser release --rm-dist
```
