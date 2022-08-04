package lode

import (
	"encoding/json"
	"errors"
	"github.com/JamesBalazs/lode/internal/lode/mocks"
	"github.com/JamesBalazs/lode/internal/responseTimings"
	"github.com/JamesBalazs/lode/internal/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"
	"testing"
	"time"
)

var params = Params{
	Url:         "https://www.example.com",
	Method:      "GET",
	Body:        "",
	File:        "",
	Freq:        1,
	Concurrency: 1,
	MaxRequests: 1,
	Delay:       time.Second,
	Timeout:     time.Second,
	MaxTime:     0,
	Headers:     nil,
}

func TestNewLode_ReturnsLode(t *testing.T) {
	assert := assert.New(t)
	expectedRequest, _ := http.NewRequest(params.Method, params.Url, nil)
	logMock := new(mocks.Log)
	Logger = logMock
	clientMock := new(mocks.Client)
	NewClient = func(timeout time.Duration) types.HttpClientInt {
		return clientMock
	}
	NewRequest = func(method, url string, body io.Reader) (*http.Request, error) {
		return expectedRequest, nil
	}
	var nilLogger *log.Logger
	expectedLode := &Lode{
		TargetDelay:     params.Delay,
		Client:          clientMock,
		Request:         expectedRequest,
		Concurrency:     1,
		MaxRequests:     1,
		MaxTime:         0,
		StartTime:       time.Time{},
		ResponseTimings: responseTimings.ResponseTimings(nil),
		OutFormat:       "json",
		FileLogger:      nilLogger,
	}

	lode := New(params)

	assert.Equal(expectedLode, lode)
}

func TestNewLode_ErrorCreatingRequest(t *testing.T) {
	assert := assert.New(t)
	logMock := new(mocks.Log)
	Logger = logMock
	clientMock := new(mocks.Client)
	NewClient = func(timeout time.Duration) types.HttpClientInt {
		return clientMock
	}
	NewRequest = func(string, string, io.Reader) (*http.Request, error) {
		return nil, errors.New("could not create request")
	}
	logMock.On("Panicf", "Error creating request: %s", "could not create request").Once()

	lode := New(params)

	assert.Nil(lode)
	logMock.AssertExpectations(t)
	NewClient = func(timeout time.Duration) types.HttpClientInt {
		return &http.Client{Timeout: timeout}
	}
	NewRequest = http.NewRequest
}

func TestNewLode_SetsBody(t *testing.T) {
	params.Body = "{\"example\":\"value\"}"
	expectedBody := io.NopCloser(strings.NewReader(params.Body))

	lode := New(params)

	assert.Equal(t, expectedBody, lode.Request.Body)
	params.Body = ""
}

func TestNewLode_SetsHeaders(t *testing.T) {
	params.Headers = []string{"Content-Type=application/json", "X-Something=value"}
	expectedHeader := http.Header{"Content-Type": {"application/json"}, "X-Something": {"value"}}

	lode := New(params)

	assert.Equal(t, expectedHeader, lode.Request.Header)
}

func TestNewLode_DefaultTimeout(t *testing.T) {
	params.Timeout = 0
	expectedTimeout := 5 * time.Second

	lode := New(params)
	client := lode.Client.(*http.Client)

	assert.Equal(t, expectedTimeout, client.Timeout)
}

func TestNewLode_FileLogger(t *testing.T) {
	params.OutFile = "/tmp/outfile.txt"

	lode := New(params)

	assert.NotNil(t, lode.FileLogger)
	params.OutFile = ""
}

func TestNewLode_OutfileYaml(t *testing.T) {
	params.OutFile = "/tmp/outfile.txt"
	params.OutFormat = "yaml"
	expectedOutFormat := "yaml"

	lode := New(params)

	assert.Equal(t, lode.OutFormat, expectedOutFormat)
	params.OutFile = ""
	params.OutFormat = ""
}

func TestLode_RunDoesRequest(t *testing.T) {
	clientMock := new(mocks.Client)
	NewClient = func(timeout time.Duration) types.HttpClientInt {
		return clientMock
	}
	response := &http.Response{
		StatusCode:    200,
		ContentLength: 3,
		Body:          io.NopCloser(strings.NewReader("abc")),
	}
	clientMock.On("Do", mock.Anything).Return(response, nil).Once()
	logMock := new(mocks.Log)
	Logger = logMock

	lode := New(params)
	lode.Run()

	clientMock.AssertExpectations(t)
	logMock.AssertExpectations(t)
}

func TestLode_RunInteractiveStoresBodyAndHeaders(t *testing.T) {
	assert := assert.New(t)
	clientMock := new(mocks.Client)
	NewClient = func(timeout time.Duration) types.HttpClientInt {
		return clientMock
	}
	body := "someBody"
	header := http.Header{"Set-Cookie": {`abc="def"`}}
	response := &http.Response{
		StatusCode:    200,
		ContentLength: 3,
		Body:          io.NopCloser(strings.NewReader(body)),
		Header:        header,
	}
	clientMock.On("Do", mock.Anything).Return(response, nil).Once()
	logMock := new(mocks.Log)
	Logger = logMock
	lode := New(params)
	lode.Interactive = true

	lode.Run()

	clientMock.AssertExpectations(t)
	logMock.AssertExpectations(t)
	assert.Equal(body, lode.ResponseTimings[0].Response.Body)
	assert.Equal(responseTimings.Header{HttpHeader: header}, lode.ResponseTimings[0].Response.Header)
}

func TestLode_RunErrorDoingRequest(t *testing.T) {
	clientMock := new(mocks.Client)
	NewClient = func(timeout time.Duration) types.HttpClientInt {
		return clientMock
	}
	clientMock.On("Do", mock.Anything).Return(&http.Response{}, errors.New("error doing request"))
	logMock := new(mocks.Log)
	logMock.On("Panicf", "Error during request: %s", "error doing request").Once()
	Logger = logMock

	lode := New(params)
	lode.Run()

	logMock.AssertExpectations(t)
	clientMock.AssertExpectations(t)
}

func TestLode_RunFailFast(t *testing.T) {
	clientMock := new(mocks.Client)
	NewClient = func(timeout time.Duration) types.HttpClientInt {
		return clientMock
	}
	response := &http.Response{
		StatusCode:    400,
		ContentLength: 3,
		Body:          io.NopCloser(strings.NewReader("abc")),
	}
	clientMock.On("Do", mock.Anything).Return(response, nil).Once()
	logMock := new(mocks.Log)
	logMock.On("Fatalf", "Got non-success status code: %d", 400).Once()
	Logger = logMock

	oldFailFast := params.FailFast
	defer func() { params.FailFast = oldFailFast }()
	params.FailFast = true
	lode := New(params)
	lode.Run()

	clientMock.AssertExpectations(t)
	logMock.AssertExpectations(t)
}

func TestLode_RunNonZeroExitCode(t *testing.T) {
	clientMock := new(mocks.Client)
	NewClient = func(timeout time.Duration) types.HttpClientInt {
		return clientMock
	}
	response := &http.Response{
		StatusCode:    400,
		ContentLength: 3,
		Body:          io.NopCloser(strings.NewReader("abc")),
	}
	clientMock.On("Do", mock.Anything).Return(response, nil).Once()
	logMock := new(mocks.Log)
	Logger = logMock

	lode := New(params)
	lode.Run()

	clientMock.AssertExpectations(t)
	logMock.AssertExpectations(t)
	assert.Equal(t, 1, lode.ExitCode)
}

func TestLode_RunIgnoreFailures(t *testing.T) {
	clientMock := new(mocks.Client)
	NewClient = func(timeout time.Duration) types.HttpClientInt {
		return clientMock
	}
	response := &http.Response{
		StatusCode:    400,
		ContentLength: 3,
		Body:          io.NopCloser(strings.NewReader("abc")),
	}
	clientMock.On("Do", mock.Anything).Return(response, nil).Once()
	logMock := new(mocks.Log)
	Logger = logMock

	oldIgnoreFailures := params.IgnoreFailures
	defer func() { params.IgnoreFailures = oldIgnoreFailures }()
	params.IgnoreFailures = true
	lode := New(params)
	lode.Run()

	clientMock.AssertExpectations(t)
	logMock.AssertExpectations(t)
	assert.Equal(t, 0, lode.ExitCode)
}

func TestLode_RunOutfile(t *testing.T) {
	clientMock := new(mocks.Client)
	NewClient = func(timeout time.Duration) types.HttpClientInt {
		return clientMock
	}
	response := &http.Response{
		StatusCode:    200,
		ContentLength: 3,
		Body:          io.NopCloser(strings.NewReader("abc")),
	}
	clientMock.On("Do", mock.Anything).Return(response, nil).Once()
	logMock := new(mocks.Log)
	fileMock := new(mocks.Log)
	fileMock.On("Println", mock.MatchedBy(func(str string) bool {
		var respTiming responseTimings.ResponseTiming
		err := json.Unmarshal([]byte(str), &respTiming)
		return err == nil && respTiming.Response != &responseTimings.Response{}
	})).Return().Once()

	Logger = logMock
	lode := New(params)
	lode.WriteFile = true
	lode.FileLogger = fileMock

	lode.Run()

	clientMock.AssertExpectations(t)
	logMock.AssertExpectations(t)
	fileMock.AssertExpectations(t)
}

func TestLode_ReportOneRequest(t *testing.T) {
	logMock := new(mocks.Log)
	Logger = logMock
	lode := New(params)
	lode.ResponseTimings = responseTimings.ResponseTimings{
		responseTimings.ResponseTiming{Response: &responseTimings.Response{}},
	}
	logMock.On("Printf", mock.MatchedBy(func(str string) bool {
		result, _ := regexp.MatchString("Timing breakdown", str)
		return result
	})).Once()

	lode.Report()

	logMock.AssertExpectations(t)
}

func TestLode_ReportMultipleRequests(t *testing.T) {
	logMock := new(mocks.Log)
	Logger = logMock
	lode := New(params)
	lode.ResponseTimings = responseTimings.ResponseTimings{
		responseTimings.ResponseTiming{Response: &responseTimings.Response{}, Timing: &responseTimings.Timing{}},
		responseTimings.ResponseTiming{Response: &responseTimings.Response{}, Timing: &responseTimings.Timing{}},
	}
	logMock.On("Printf", mock.MatchedBy(func(str string) bool {
		result1, _ := regexp.MatchString("Response code breakdown", str)
		result2, _ := regexp.MatchString("Percentile latency breakdown", str)
		return result1 && result2
	})).Once()

	lode.Report()

	logMock.AssertExpectations(t)
}

func TestLode_ReportNoRequests(t *testing.T) {
	logMock := new(mocks.Log)
	Logger = logMock
	lode := New(params)
	lode.ResponseTimings = responseTimings.ResponseTimings{}
	logMock.On("Printf", mock.MatchedBy(func(str string) bool {
		result, _ := regexp.MatchString("No requests made...", str)
		return result
	})).Once()

	lode.Report()

	logMock.AssertExpectations(t)
}

func TestLode_ReportInteractive(t *testing.T) {
	logMock := new(mocks.Log)
	Logger = logMock
	lode := New(params)
	lode.Interactive = true
	promptuiSelectMock := mocks.Select{}
	oldNewInteractivePrompt := newInteractivePrompt
	newInteractivePrompt = func(label string, responseTimings responseTimings.ResponseTimings) types.PromptSelectInt {
		return &promptuiSelectMock
	}
	defer func() { newInteractivePrompt = oldNewInteractivePrompt }()
	lode.ResponseTimings = responseTimings.ResponseTimings{
		responseTimings.ResponseTiming{Response: &responseTimings.Response{}, Timing: &responseTimings.Timing{}},
	}
	logMock.On("Printf", mock.MatchedBy(func(str string) bool {
		result1, _ := regexp.MatchString("Response code breakdown", str)
		result2, _ := regexp.MatchString("Percentile latency breakdown", str)
		return result1 && result2
	})).Once()
	promptuiSelectMock.On("Run").Return(0, "", nil).Once()

	lode.Report()

	promptuiSelectMock.AssertExpectations(t)
	logMock.AssertExpectations(t)
}

func TestLode_ReportInteractiveError(t *testing.T) {
	logMock := new(mocks.Log)
	Logger = logMock
	lode := New(params)
	lode.Interactive = true
	promptuiSelectMock := mocks.Select{}
	oldNewInteractivePrompt := newInteractivePrompt
	newInteractivePrompt = func(label string, responseTimings responseTimings.ResponseTimings) types.PromptSelectInt {
		return &promptuiSelectMock
	}
	defer func() { newInteractivePrompt = oldNewInteractivePrompt }()
	lode.ResponseTimings = responseTimings.ResponseTimings{
		responseTimings.ResponseTiming{Response: &responseTimings.Response{}, Timing: &responseTimings.Timing{}},
	}
	logMock.On("Printf", mock.MatchedBy(func(str string) bool {
		result1, _ := regexp.MatchString("Response code breakdown", str)
		result2, _ := regexp.MatchString("Percentile latency breakdown", str)
		return result1 && result2
	})).Once()
	promptuiSelectMock.On("Run").Return(0, "", errors.New("promptui error")).Once()
	logMock.On("Panicln", "promptui error").Once()

	lode.Report()

	promptuiSelectMock.AssertExpectations(t)
	logMock.AssertExpectations(t)
}
