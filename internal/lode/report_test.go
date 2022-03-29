package lode

import (
	"errors"
	"github.com/JamesBalazs/lode/internal/lode/mocks"
	"github.com/JamesBalazs/lode/internal/responseTimings"
	"github.com/JamesBalazs/lode/internal/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"testing"
	"text/template"
	"time"
)

var responseTiming = responseTimings.ResponseTiming{
	Response: &responseTimings.Response{
		Status:     "200 OK",
		StatusCode: 200,
	},
	Timing: &responseTimings.Timing{
		ConnectStart: time.Unix(0, 1_000_000),
		Done:         time.Unix(0, 3_000_000),
	},
}

func TestNewTestReport(t *testing.T) {
	request, _ := http.NewRequest(params.Method, params.Url, nil)
	logMock := new(mocks.Log)
	Logger = logMock
	clientMock := new(mocks.Client)
	responseTimings := responseTimings.ResponseTimings{
		responseTiming,
		responseTiming,
	}
	lode := &Lode{
		TargetDelay:     params.Delay,
		Client:          clientMock,
		Request:         request,
		Concurrency:     2,
		MaxRequests:     100,
		MaxTime:         10,
		StartTime:       time.Time{},
		FinishTime:      time.Time{}.Add(10 * time.Second),
		ResponseTimings: responseTimings,
	}
	expectedReport := TestReport{
		Target:          "GET https://www.example.com",
		Concurrency:     2,
		Duration:        10 * time.Second,
		ResponseCount:   2,
		RequestRate:     0.2,
		ResponseTimings: responseTimings,
	}

	tr := NewTestReport(lode)

	assert.Equal(t, expectedReport, tr)
}

func TestTestReport_FirstResponse(t *testing.T) {
	tr := TestReport{
		ResponseTimings: responseTimings.ResponseTimings{
			responseTiming,
			responseTimings.ResponseTiming{},
		},
	}

	assert.Equal(t, responseTiming, tr.FirstResponse())
}

func TestTestReport_LatencyPercentiles(t *testing.T) {

}

func TestTestReport_StatusHistogram(t *testing.T) {

}

func TestTestReport_OneResponse(t *testing.T) {
	assert := assert.New(t)
	tr := TestReport{
		ResponseCount: 2,
	}
	assert.Equal(false, tr.OneResponse())

	tr.ResponseCount = 1
	assert.Equal(true, tr.OneResponse())
}

func TestTestReport_MultipleResponses(t *testing.T) {
	assert := assert.New(t)
	tr := TestReport{
		ResponseCount: 1,
	}
	assert.Equal(false, tr.MultipleResponses())

	tr.ResponseCount = 2
	assert.Equal(true, tr.MultipleResponses())
}

func TestTestReport_Output(t *testing.T) {
	assert := assert.New(t)
	tr := TestReport{
		Target:        "GET https://www.example.com",
		Concurrency:   4,
		Duration:      10 * time.Second,
		ResponseCount: 2,
		RequestRate:   0.2,
		ResponseTimings: responseTimings.ResponseTimings{
			responseTiming,
		},
	}
	output := tr.Output()
	assert.Contains(output, `Target: GET https://www.example.com
Concurrency: 4
Requests made: 2
Time taken: 10s
Requests per second (avg): 0.2

Response code breakdown:`)
	assert.Contains(output, "Percentile latency breakdown:")
	assert.NotContains(output, "Timing breakdown:")
	assert.NotContains(output, "No requests made...")

	tr.ResponseCount = 1
	output = tr.Output()
	assert.Contains(output, "Timing breakdown:")
	assert.NotContains(output, "Response code breakdown:")
	assert.NotContains(output, "Percentile latency breakdown:")
	assert.NotContains(output, "No requests made...")

	tr.ResponseCount = 0
	output = tr.Output()
	assert.Contains(output, "No requests made...")
	assert.NotContains(output, "Timing breakdown:")
	assert.NotContains(output, "Response code breakdown:")
	assert.NotContains(output, "Percentile latency breakdown:")

}

func TestTestReport_OutputErrorParsingTemplate(t *testing.T) {
	logMock := new(mocks.Log)
	Logger = logMock
	templateMock := new(mocks.Template)
	oldNewTemplate := newTemplate
	defer func() { newTemplate = oldNewTemplate }()
	newTemplate = func(name string) types.TemplateInt {
		return templateMock
	}

	templateMock.On("Parse", mock.AnythingOfType("string")).Return(&template.Template{}, errors.New("invalid template"))
	logMock.On("Panicf", "Error parsing report template: %s", "invalid template").Return().Once()

	tr := TestReport{}
	tr.Output()

	templateMock.AssertExpectations(t)
	logMock.AssertExpectations(t)
}
