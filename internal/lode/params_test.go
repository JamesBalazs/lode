package lode

import (
	"github.com/JamesBalazs/lode/internal/lode/mocks"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func TestParams_Validate(t *testing.T) {
	oldLogger := Logger
	defer func() { Logger = oldLogger }()
	logMock := new(mocks.Log)
	Logger = logMock

	invalidSuite := "Invalid test suite:\n%s\n"

	oldParam := Params{
		Url:         "https://www.google.com",
		Method:      "GET",
		Freq:        1,
		Concurrency: 1,
		MaxRequests: 1,
		Delay:       1 * time.Second,
		MaxTime:     1 * time.Second,
		Timeout:     1 * time.Second,
	}
	param := oldParam

	param.Validate()

	logMock.AssertNotCalled(t, "Panicf", mock.AnythingOfType("string"))

	param.Url = ""
	logMock.On("Panicf", invalidSuite, "url must be provided").Return().Once()
	param.Validate()
	logMock.AssertExpectations(t)
	param.Url = oldParam.Url

	param.Method = ""
	logMock.On("Panicf", invalidSuite, "method must be provided").Return().Once()
	param.Validate()
	logMock.AssertExpectations(t)
	param.Method = oldParam.Method

	param.Freq = 0
	param.Validate()
	logMock.AssertNotCalled(t, "Panicf", invalidSuite, "freq or delay must be provided")
	param.Freq = oldParam.Freq

	param.Delay = 0
	param.Validate()
	logMock.AssertNotCalled(t, "Panicf", invalidSuite, "freq or delay must be provided")
	param.Delay = oldParam.Delay

	param.Freq, param.Delay = 0, 0
	logMock.On("Panicf", invalidSuite, "freq or delay must be provided").Return().Once()
	param.Validate()
	logMock.AssertExpectations(t)
	param.Freq, param.Delay = oldParam.Freq, oldParam.Delay

	param.Concurrency = 0
	logMock.On("Panicf", invalidSuite, "concurrency must be provided as a positive integer").Return().Once()
	param.Validate()
	logMock.AssertExpectations(t)
	param.Concurrency = oldParam.Concurrency

	param.Timeout = 0
	logMock.On("Panicf", invalidSuite, "timeout must be provided").Return().Once()
	param.Validate()
	logMock.AssertExpectations(t)
	param.Timeout = oldParam.Timeout

	param.MaxRequests = 0
	param.Validate()
	logMock.AssertNotCalled(t, "Panicf", invalidSuite, "maxrequests or maxtime must be provided")
	param.MaxRequests = oldParam.MaxRequests

	param.MaxTime = 0
	param.Validate()
	logMock.AssertNotCalled(t, "Panicf", invalidSuite, "maxrequests or maxtime must be provided")
	param.MaxTime = oldParam.MaxTime

	param.MaxRequests, param.MaxTime = 0, 0
	logMock.On("Panicf", invalidSuite, "maxrequests or maxtime must be provided").Return().Once()
	param.Validate()
	logMock.AssertExpectations(t)
	param.MaxRequests, param.MaxTime = oldParam.MaxRequests, oldParam.MaxTime

	param.OutFormat = "json"
	logMock.On("Panicf", invalidSuite, "outFormat must be used with outFile").Return().Once()
	param.Validate()
	logMock.AssertExpectations(t)
	param.OutFormat = oldParam.OutFormat

	param.OutFile, param.OutFormat = "/tmp/out.txt", "invalid"
	logMock.On("Panicf", invalidSuite, "invalid outFormat - valid options are json and yaml").Return().Once()
	param.Validate()
	logMock.AssertExpectations(t)
	param.OutFile, param.OutFormat = oldParam.OutFile, oldParam.OutFormat
}
