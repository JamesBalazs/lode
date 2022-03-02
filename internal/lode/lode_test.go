package lode

import (
	"errors"
	"github.com/JamesBalazs/lode/internal/lode/mocks"
	"github.com/JamesBalazs/lode/internal/lode/report"
	"github.com/JamesBalazs/lode/internal/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"io"
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
	expectedLode := &Lode{
		TargetDelay:     params.Delay,
		Client:          clientMock,
		Request:         expectedRequest,
		Concurrency:     1,
		MaxRequests:     1,
		MaxTime:         0,
		StartTime:       time.Time{},
		ResponseTimings: ResponseTimings(nil),
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
	logMock.On("Panicf", "Error creating request: %s", "could not create request")

	lode := New(params)

	assert.Nil(lode)
	logMock.AssertExpectations(t)
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

func TestLode_RunDoesRequest(t *testing.T) {
	clientMock := new(mocks.Client)
	NewClient = func(timeout time.Duration) types.HttpClientInt {
		return clientMock
	}
	response := &http.Response{
		ContentLength: 3,
		Body:          io.NopCloser(strings.NewReader("abc")),
	}
	clientMock.On("Do", mock.Anything).Return(response, nil)
	logMock := new(mocks.Log)
	Logger = logMock

	lode := New(params)
	lode.Run()

	clientMock.AssertExpectations(t)
	logMock.AssertExpectations(t)
}

func TestLode_RunErrorDoingRequest(t *testing.T) {
	clientMock := new(mocks.Client)
	NewClient = func(timeout time.Duration) types.HttpClientInt {
		return clientMock
	}
	clientMock.On("Do", mock.Anything).Return(&http.Response{}, errors.New("error doing request"))
	logMock := new(mocks.Log)
	logMock.On("Panicf", "Error during request: %s", "error doing request")
	Logger = logMock

	lode := New(params)
	lode.Run()

	logMock.AssertExpectations(t)
	clientMock.AssertExpectations(t)
}

func TestLode_ReportOneRequest(t *testing.T) {
	logMock := new(mocks.Log)
	Logger = logMock
	lode := New(params)
	lode.ResponseTimings = ResponseTimings{
		ResponseTiming{Response: &types.Response{}},
	}
	logMock.On("Printf", mock.MatchedBy(func(str string) bool {
		result, _ := regexp.MatchString("Timing breakdown", str)
		return result
	})).Return()

	lode.Report()

	logMock.AssertExpectations(t)
}

func TestLode_ReportMultipleRequests(t *testing.T) {
	logMock := new(mocks.Log)
	Logger = logMock
	lode := New(params)
	lode.ResponseTimings = ResponseTimings{
		ResponseTiming{Response: &types.Response{}, Timing: &report.Timing{}},
		ResponseTiming{Response: &types.Response{}, Timing: &report.Timing{}},
	}
	logMock.On("Printf", mock.MatchedBy(func(str string) bool {
		result1, _ := regexp.MatchString("Response code breakdown", str)
		result2, _ := regexp.MatchString("Percentile latency breakdown", str)
		return result1 && result2
	})).Return()

	lode.Report()

	logMock.AssertExpectations(t)
}

func TestLode_ReportNoRequests(t *testing.T) {
	logMock := new(mocks.Log)
	Logger = logMock
	lode := New(params)
	lode.ResponseTimings = ResponseTimings{}
	logMock.On("Printf", mock.MatchedBy(func(str string) bool {
		result, _ := regexp.MatchString("No requests made...", str)
		return result
	})).Return()

	lode.Report()

	logMock.AssertExpectations(t)
}
