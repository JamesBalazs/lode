# Lode

Versatile load testing CLI tool written in Go, with configurable workflows to facilitate automated load testing in CI. 

**Features:**
- Lightweight
- Portable
- Concurrent
- Configurable
- Open source

## Example output
```
❯ lode test --freq=20 -c 8 -l 5s http://my.example.service
Target: GET http://my.example.service
Concurrency: 8
Requests made: 100
Time taken: 5.04s
Requests per second (avg): 19.84

Response Breakdown:
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

## Usage
### `lode test [flags] [url]`
Used to run a single load test against the specified URL.
A summary report is printed at the end.

**Supported flags:**
| Flag | Shorthand | Usage |
| --- | --- | --- |
| `--freq` | `-f` | Number of requests to make per second |
| `--delay` | `-d` | Time to wait between requests, e.g. 200ms or 1s - defaults to 1s unless --freq specified |
| `--timeout` | `-t` | Timeout per request, e.g. 200ms or 1s - defaults to 5s |
| `--method` | `-m` | HTTP method to use - defaults to GET |
| `--concurrency` | `-c` | Maximum number of concurrent requests |
| `--maxRequests` | `-n` | Maximum number of requests to make - defaults to 0s (unlimited) |
| `--maxTime` | `-l` | Length of time to make requests, e.g. 20s or 1h - defaults to 0s (unlimited) |

One of either `--delay` or `--freq` is required. If both are provided, delay will be calculated from the given frequency.

**Examples:**
- `lode test -f 20 -c 4 -l 10s http://www.google.com` make 20 req/sec to Google for 10 seconds, split across 4 threads
- `lode test -d 1h -n 24 http://www.google.com` make 1 req/hr to Google until 24 requests have been made
- `lode test -f 40 -c 8 -l 1m -n 1000 http://www.google.copm` make 40 req/sec to Google, split across 8 threads, for up to 1 minutes or until 1000 requests have been made (whichever comes first) 

### `lode workflow [flags] [path]` (not yet implemented)
Used to run a sequence of load tests, as defined in the YAML file at the specified path.

### `lode time [flags] [path]` (not yet implemented)
Used to run a single request, and print the timings of each stage of the request.


**Expected YAML format:**
```
jobs:
 - URL: http://www.google.com/
   Method: HEAD
   Concurrency: 8
   Freq: 40
   MaxTime: 10m
 - URL: http://www.example.com/
   Method: HEAD
   Freq: 1
   MaxTime: 10m
```

## Planned Features
- Workflows for CI/automation use
- Log responses to a file
- Better analysis for recorded response timing
