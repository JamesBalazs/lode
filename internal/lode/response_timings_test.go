package lode

import (
	"github.com/JamesBalazs/lode/internal/lode/report"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

func TestResponseTimings_Responses(t *testing.T) {
	assert := assert.New(t)

	response := http.Response{StatusCode: 200}
	timing := report.Timing{Done: time.Now()}
	responseTiming := ResponseTiming{Response: response, Timing: timing}
	responseTimings := ResponseTimings{
		responseTiming,
		responseTiming,
	}

	assert.Equal([]http.Response{response, response}, responseTimings.Responses())
}

func TestResponseTimings_Timings(t *testing.T) {
	assert := assert.New(t)

	response := http.Response{StatusCode: 200}
	timing := report.Timing{Done: time.Now()}
	responseTiming := ResponseTiming{Response: response, Timing: timing}
	responseTimings := ResponseTimings{
		responseTiming,
		responseTiming,
	}

	assert.Equal([]report.Timing{timing, timing}, responseTimings.Timings())
}
