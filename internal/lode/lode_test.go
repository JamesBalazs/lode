package lode

import (
	"errors"
	"github.com/JamesBalazs/lode/internal/lode/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"
)

var params = Params{
	Url:         "https://www.example.com",
	Method:      "GET",
	Body:        "",
	File:        "",
	Freq:        0,
	Concurrency: 1,
	MaxRequests: 1,
	Delay:       time.Second,
	Timeout:     0,
	MaxTime:     0,
	Headers:     nil,
}

var clientMock = &mocks.Client{}

func TestNewLode_ReturnsLode(t *testing.T) {
	assert := assert.New(t)
	NewClient = func(timeout time.Duration) HttpClientInt {
		return clientMock
	}
	expectedRequest, _ := http.NewRequest(params.Method, params.Url, nil)
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
	logMock.On("Panicf", "Error creating request: %s", "could not create request")
	Logger = logMock
	NewRequest = func(string, string, io.Reader) (*http.Request, error) {
		return nil, errors.New("could not create request")
	}

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
	response := &http.Response{}
	clientMock.On("Do", mock.Anything).Return(response, nil)
	logMock := new(mocks.Log)
	Logger = logMock

	lode := New(params)
	lode.Run()

	clientMock.AssertExpectations(t)
	logMock.AssertExpectations(t)
}

func TestLode_RunErrorDoingRequest(t *testing.T) {
	clientMock.On("Do", mock.Anything).Return(&http.Response{}, errors.New("error doing request"))
	logMock := new(mocks.Log)
	logMock.On("Panicf", "Error during request: %s", "error doing request")
	Logger = logMock

	lode := New(params)
	lode.Run()

	clientMock.AssertExpectations(t)
}

func TestLode_Report(t *testing.T) {
	clientMock := new(mocks.Client)
	response := &http.Response{}
	clientMock.On("Do", mock.Anything).Return(response, nil)
	logMock := new(mocks.Log)
	Logger = logMock
	lode := New(params)
	lode.Run()
	logMock.On("Printf", mock.AnythingOfType("string")).Return() // Report after requests TODO: find a more specific way to mock this

	lode.Report()

	logMock.AssertExpectations(t)
}
