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

const url = "https://www.example.com"
const method = "GET"
const delay = time.Second

var client = &mocks.Client{}

func TestNewLode_ReturnsLode(t *testing.T) {
	assert := assert.New(t)
	expectedRequest, _ := http.NewRequest(method, url, nil)
	expectedLode := &Lode{
		TargetDelay:     delay,
		Client:          client,
		Request:         expectedRequest,
		Concurrency:     1,
		MaxRequests:     1,
		MaxTime:         0,
		StartTime:       time.Time{},
		ResponseTimings: ResponseTimings(nil),
	}

	lode := New(url, method, delay, client, 1, 1, 0, nil, nil)

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

	lode := New(url, method, delay, client, 1, 1, 0, nil, nil)

	assert.Nil(lode)
	logMock.AssertExpectations(t)
	NewRequest = http.NewRequest
}

func TestNewLode_SetsBody(t *testing.T) {
	body := strings.NewReader("{\"example\":\"value\"}")
	expectedBody := io.NopCloser(body)

	lode := New(url, method, delay, client, 1, 1, 0, body, nil)

	assert.Equal(t, expectedBody, lode.Request.Body)
}

func TestNewLode_SetsHeaders(t *testing.T) {
	headers := []string{"Content-Type=application/json", "X-Something=value"}
	expectedHeader := http.Header{"Content-Type": {"application/json"}, "X-Something": {"value"}}

	lode := New(url, method, delay, client, 1, 1, 0, nil, headers)

	assert.Equal(t, expectedHeader, lode.Request.Header)
}

func TestLode_RunDoesRequest(t *testing.T) {
	clientMock := new(mocks.Client)
	response := &http.Response{}
	clientMock.On("Do", mock.Anything).Return(response, nil)
	logMock := new(mocks.Log)
	Logger = logMock

	lode := New(url, method, delay, clientMock, 1, 1, 0, nil, nil)
	lode.Run()

	clientMock.AssertExpectations(t)
	logMock.AssertExpectations(t)
}

func TestLode_RunErrorDoingRequest(t *testing.T) {
	clientMock := new(mocks.Client)
	clientMock.On("Do", mock.Anything).Return(&http.Response{}, errors.New("error doing request"))
	logMock := new(mocks.Log)
	logMock.On("Panicf", "Error during request: %s", "error doing request")
	Logger = logMock

	lode := New(url, method, delay, clientMock, 1, 1, 0, nil, nil)
	lode.Run()

	clientMock.AssertExpectations(t)
}

func TestLode_Report(t *testing.T) {
	clientMock := new(mocks.Client)
	response := &http.Response{}
	clientMock.On("Do", mock.Anything).Return(response, nil)
	logMock := new(mocks.Log)
	Logger = logMock
	lode := New(url, method, delay, clientMock, 1, 1, 0, nil, nil)
	lode.Run()
	logMock.On("Printf", mock.AnythingOfType("string")).Return() // Report after requests TODO: find a more specific way to mock this

	lode.Report()

	logMock.AssertExpectations(t)
}
