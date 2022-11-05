package lode

import (
	"context"
	"encoding/json"
	"github.com/JamesBalazs/lode/internal/files"
	"github.com/JamesBalazs/lode/internal/responseTimings"
	"github.com/JamesBalazs/lode/internal/types"
	"gopkg.in/yaml.v3"
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

var Logger types.LoggerInt = log.New(os.Stdout, "", 0)
var NewRequest = http.NewRequest
var NewClient = func(timeout time.Duration) types.HttpClientInt {
	return &http.Client{Timeout: timeout}
}

type Lode struct {
	Client          types.HttpClientInt
	Request         *http.Request
	Concurrency     int
	MaxRequests     int
	ExitCode        int
	TargetDelay     time.Duration
	MaxTime         time.Duration
	StartTime       time.Time
	FinishTime      time.Time
	ResponseTimings responseTimings.ResponseTimings
	FailFast        bool
	IgnoreFailures  bool
	Interactive     bool
	WriteFile       bool
	FileLogger      types.LoggerInt
	OutFormat       string
}

func New(params Params) *Lode {
	if params.Timeout == 0 {
		params.Timeout = 5 * time.Second
	}
	if params.Freq != 0 {
		params.Delay = time.Second / time.Duration(params.Freq)
	}
	params.Validate()

	body := files.ReaderFromFileOrString(params.File, params.Body)
	req, err := NewRequest(params.Method, params.Url, body)
	if err != nil {
		Logger.Panicf("Error creating request: %s", err.Error())
		return nil
	}

	for _, headerString := range params.Headers {
		headerParts := strings.SplitN(headerString, "=", 2)
		req.Header[headerParts[0]] = []string{headerParts[1]}
	}

	var fileLogger *log.Logger
	writeFile := len(params.OutFile) != 0
	if writeFile {
		fileLogger = newFileLogger(params.OutFile)
	}

	outFormat := "json"
	if params.OutFormat == "yaml" {
		outFormat = "yaml"
	}

	return &Lode{
		TargetDelay:    params.Delay,
		Client:         NewClient(params.Timeout),
		Request:        req,
		Concurrency:    params.Concurrency,
		MaxRequests:    params.MaxRequests,
		MaxTime:        params.MaxTime,
		FailFast:       params.FailFast,
		IgnoreFailures: params.IgnoreFailures,
		FileLogger:     fileLogger,
		OutFormat:      outFormat,
		WriteFile:      writeFile,
	}
}

func (l *Lode) Run() {
	ticker := time.NewTicker(l.TargetDelay)
	trigger := ticker.C
	stop := make(chan struct{})
	defer l.stop(ticker, stop)
	l.StartTime = time.Now()
	defer l.setFinishTime()

	result := make(chan responseTimings.ResponseTiming, 1024)
	l.closeOnSigterm(result)

	for i := 0; i < l.Concurrency; i++ {
		go l.work(trigger, stop, result)
	}

	startTime := time.Now()
	endTime := startTime.Add(l.MaxTime).UnixNano()
	checkMaxRequests := l.MaxRequests > 0
	checkMaxTime := l.MaxTime > 0
	responseCount := 0
	outMarshaller := l.outMarshaller()

	for response := range result {
		responseCount++
		l.ResponseTimings = append(l.ResponseTimings, response)

		if l.WriteFile {
			marshalledResponse, err := outMarshaller(response)
			if err != nil {
				panic(err)
			}

			l.FileLogger.Println(string(marshalledResponse))
		}

		if !l.IgnoreFailures && (response.Response.StatusCode < 100 || response.Response.StatusCode >= 400) {
			l.ExitCode = 1
		}

		if (checkMaxRequests && responseCount >= l.MaxRequests) || (checkMaxTime && time.Now().UnixNano() >= endTime) {
			return
		}
	}
}

func (l Lode) work(trigger <-chan time.Time, stop chan struct{}, result chan responseTimings.ResponseTiming) {
	ctx := context.Background()
	for {
		select {
		case <-trigger:
			response, timing := l.makeAndTimeRequest(ctx)
			result <- responseTimings.ResponseTiming{
				Response: response,
				Timing:   timing,
			}
		case <-stop:
			return
		}
	}
}

func (l Lode) stop(ticker *time.Ticker, stop chan struct{}) {
	ticker.Stop()
	close(stop)
}

func (l *Lode) Report() {
	report := NewTestReport(l)
	output := report.Output()

	if l.Interactive {
		output += "Requests:\n"
		Logger.Printf(output)
		prompt := newInteractivePrompt(output, l.ResponseTimings)
		_, _, err := prompt.Run()
		if err != nil {
			Logger.Panicln(err.Error())
		}
	} else {
		Logger.Printf(output)
	}
}

func (l Lode) closeOnSigterm(channel chan responseTimings.ResponseTiming) {
	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigterm
		close(channel)
	}()
}

func (l *Lode) setFinishTime() {
	l.FinishTime = time.Now()
}

func (l Lode) outMarshaller() func(v interface{}) ([]byte, error) {
	if l.OutFormat == "yaml" {
		return yaml.Marshal
	}
	return json.Marshal
}

func (l Lode) makeAndTimeRequest(ctx context.Context) (result *responseTimings.Response, timing *responseTimings.Timing) {
	var err error
	var response *http.Response
	timing = &responseTimings.Timing{}
	trace := responseTimings.NewTrace(timing)
	request := l.Request.WithContext(httptrace.WithClientTrace(ctx, trace))
	response, err = l.Client.Do(request)
	timing.Done = time.Now()
	if err != nil {
		Logger.Panicf("Error during request: %s", err.Error())
	} else if l.FailFast && (response.StatusCode < 100 || response.StatusCode >= 400) {
		Logger.Fatalf("Got non-success status code: %d", response.StatusCode)
	}

	result = &responseTimings.Response{
		Status:        response.Status,
		StatusCode:    response.StatusCode,
		ContentLength: response.ContentLength,
	}

	if l.Interactive || l.WriteFile {
		var body []byte
		body, err = io.ReadAll(response.Body)
		if err != nil {
			Logger.Panicf("Error reading body: %s", err.Error())
		}
		response.Body.Close()

		result.Header = responseTimings.Header{HttpHeader: response.Header}
		result.Body = string(body)
	}

	return
}

func (l *Lode) ExitWithCode() {
	if l.ExitCode != 0 {
		os.Exit(l.ExitCode)
	}
}
