package lode

import (
	"github.com/JamesBalazs/lode/internal/lode/report"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

func TestResponseTimings_Responses(t *testing.T) {
	response := http.Response{StatusCode: 200}
	timing := report.Timing{Done: time.Now()}
	responseTiming := ResponseTiming{Response: response, Timing: timing}
	responseTimings := ResponseTimings{
		responseTiming,
		responseTiming,
	}

	assert.Equal(t, []http.Response{response, response}, responseTimings.Responses())
}

func TestResponseTimings_Timings(t *testing.T) {
	response := http.Response{StatusCode: 200}
	timing := report.Timing{Done: time.Now()}
	responseTiming := ResponseTiming{Response: response, Timing: timing}
	responseTimings := ResponseTimings{
		responseTiming,
		responseTiming,
	}

	assert.Equal(t, []report.Timing{timing, timing}, responseTimings.Timings())
}

func TestResponseTimings_GetLongestDuration(t *testing.T) {
	response := http.Response{StatusCode: 200}
	responseTimings := ResponseTimings{
		ResponseTiming{Response: response, Timing: report.Timing{
			ConnectStart: time.Unix(0, 1_000_000),
			Done:         time.Unix(0, 3_000_000),
		}},
		ResponseTiming{Response: response, Timing: report.Timing{
			ConnectStart: time.Unix(0, 5_000_000),
			Done:         time.Unix(0, 10_000_000),
		}},
		ResponseTiming{Response: response, Timing: report.Timing{
			ConnectStart: time.Unix(0, 2_000_000),
			Done:         time.Unix(0, 6_000_000),
		}},
	}

	assert.Equal(t, 5*time.Millisecond, responseTimings.GetLongestDuration())
}
