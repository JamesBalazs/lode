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
	expectedLode := &Lode{delay, client, expectedRequest, 1, 0}

	lode := New(url, method, delay, client, 1, 0)

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

	lode := New(url, method, delay, client, 1, 0)

	assert.Nil(lode)
	logMock.AssertExpectations(t)
	NewRequest = http.NewRequest
}

func TestRunDoesRequest(t *testing.T) {
	clientMock := new(mocks.Client)
	response := &http.Response{}
	clientMock.On("Do", mock.Anything).Return(response, nil)

	lode := New(url, method, delay, clientMock, 1, 0)
	lode.Run()

	clientMock.AssertExpectations(t)
}


	//
