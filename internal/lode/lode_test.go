package lode

import (
	"errors"
	"github.com/JamesBalazs/lode/internal/lode/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"io"
	"net/http"
	"testing"
	"time"
)

const url = "https://www.example.com"
const method = "GET"
const delay = time.Second

var client = &mocks.Client{}

func TestNewReturnsLode(t *testing.T) {
	assert := assert.New(t)
	expectedRequest, _ := http.NewRequest(method, url, nil)
	expectedLode := &Lode{delay, client, expectedRequest, 1, 1, 0, []http.Response(nil)}

	lode := New(url, method, delay, client, 1, 1, 0)

	assert.Equal(expectedLode, lode)
}

func TestNewErrorCreatingRequest(t *testing.T) {
	assert := assert.New(t)
	logMock := new(mocks.Log)
	logMock.On("Panicf", "Error creating request: %s", "could not create request")
	Logger = logMock
	NewRequest = func(string, string, io.Reader) (*http.Request, error) {
		return nil, errors.New("could not create request")
	}

	lode := New(url, method, delay, client, 1, 1, 0)

	assert.Nil(lode)
	logMock.AssertExpectations(t)
	NewRequest = http.NewRequest
}

func TestRunDoesRequest(t *testing.T) {
	clientMock := new(mocks.Client)
	response := &http.Response{}
	clientMock.On("Do", mock.Anything).Return(response, nil)
	logMock := new(mocks.Log)
	logMock.On("Printf", mock.AnythingOfType("string")).Return() // report after requests TODO: find a more specific way to mock this
	Logger = logMock

	lode := New(url, method, delay, clientMock, 1, 1, 0)
	lode.Run()

	clientMock.AssertExpectations(t)
	logMock.AssertExpectations(t)
}

func TestRunErrorDoingRequest(t *testing.T) {
	clientMock := new(mocks.Client)
	clientMock.On("Do", mock.Anything).Return(&http.Response{}, errors.New("error doing request"))
	logMock := new(mocks.Log)
	logMock.On("Panicf", "Error during request: %s", "error doing request")
	logMock.On("Printf", mock.AnythingOfType("string")).Return() // report after requests TODO: find a more specific way to mock this
	Logger = logMock

	lode := New(url, method, delay, clientMock, 1, 1, 0)
	lode.Run()

	clientMock.AssertExpectations(t)

}

//
