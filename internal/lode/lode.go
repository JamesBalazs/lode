package lode

import (
	"context"
	"fmt"
	"github.com/JamesBalazs/lode/internal/lode/report"
	"io"
	"log"
	"net/http"
	"net/http/httptrace"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

var Logger LoggerInt = log.New(os.Stdout, "", 0)
var NewRequest = http.NewRequest

type Lode struct {
	TargetDelay     time.Duration
	Client          HttpClientInt
	Request         *http.Request
	Concurrency     int
	MaxRequests     int
	MaxTime         time.Duration
	StartTime       time.Time
	ResponseTimings ResponseTimings
}

func New(url string, method string, delay time.Duration, client HttpClientInt, concurrency int, maxRequests int, maxTime time.Duration, body io.Reader, headers []string) *Lode {
	req, err := NewRequest(method, url, body)
	if err != nil {
		Logger.Panicf("Error creating request: %s", err.Error())
		return nil
	}

	for _, headerString := range headers {
		headerParts := strings.SplitN(headerString, "=", 2)
		req.Header[headerParts[0]] = []string{headerParts[1]}
	}

	return &Lode{
		TargetDelay: delay,
		Client:      client,
		Request:     req,
		Concurrency: concurrency,
		MaxRequests: maxRequests,
		MaxTime:     maxTime,
	}
}

func (l *Lode) Run() {
	ticker := time.NewTicker(l.TargetDelay)
	trigger := ticker.C
	stop := make(chan struct{})
	defer l.stop(ticker, stop)
	l.StartTime = time.Now()

	result := make(chan ResponseTiming, 1024)
	l.closeOnSigterm(result)

	for i := 0; i < l.Concurrency; i++ {
		go l.work(trigger, stop, result)
	}

	startTime := time.Now()
	endTime := startTime.Add(l.MaxTime).UnixNano()
	checkMaxRequests := l.MaxRequests > 0
	checkMaxTime := l.MaxTime > 0
	responseCount := 0
	for response := range result {
		responseCount++
		l.ResponseTimings = append(l.ResponseTimings, response)

		if (checkMaxRequests && responseCount >= l.MaxRequests) || (checkMaxTime && time.Now().UnixNano() >= endTime) {
			return
		}
	}
}

func (l *Lode) work(trigger <-chan time.Time, stop chan struct{}, result chan ResponseTiming) {
	ctx := context.Background()
	for {
		select {
		case <-trigger:
			timing := report.Timing{}
			trace := report.NewTrace(&timing)
			request := l.Request.WithContext(httptrace.WithClientTrace(ctx, trace))
			response, err := l.Client.Do(request)
			timing.Done = time.Now()
			if err != nil {
				Logger.Panicf("Error during request: %s", err.Error())
			}
			result <- ResponseTiming{Response: *response, Timing: timing}
		case <-stop:
			return
		}
	}
}

func (l *Lode) stop(ticker *time.Ticker, stop chan struct{}) {
	ticker.Stop()
	close(stop)
}

func (l *Lode) Report() {
	duration := time.Since(l.StartTime)
	responseCount := len(l.ResponseTimings)
	histogram := report.BuildStatusHistogram(l.ResponseTimings.Responses(), responseCount)
	latencies := report.BuildLatencyPercentiles(l.ResponseTimings.Timings())
	requestRate := float64(responseCount) / float64(duration.Seconds())

	var output string
	output += fmt.Sprintf("Target: %s %s\n", l.Request.Method, l.Request.URL)
	output += fmt.Sprintf("Concurrency: %d\n", l.Concurrency)
	output += fmt.Sprintf("Requests made: %d\n", responseCount)
	output += fmt.Sprintf("Time taken: %s\n", duration.Truncate(10*time.Millisecond))
	output += fmt.Sprintf("Requests per second (avg): %.2f\n\n", requestRate)

	if responseCount > 1 {
		output += fmt.Sprintf("Response code breakdown:\n%s\n", histogram.String())
		output += fmt.Sprintf("Percentile latency breakdown:\n%s\n", latencies.String())
	} else if responseCount == 1 {
		timing := l.ResponseTimings[0].Timing
		output += fmt.Sprintf("Timing breakdown:\n%s\n", timing.String())
	} else {
		output += "No requests made..."
	}
	Logger.Printf(output)
}

func (l *Lode) closeOnSigterm(channel chan ResponseTiming) {
	sigterm := make(chan os.Signal)
	signal.Notify(sigterm, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigterm
		close(channel)
	}()
}
