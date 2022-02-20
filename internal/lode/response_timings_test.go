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
			ConnectStart: time.Unix(0, 1),
			Done:         time.Unix(0, 3),
		}},
		ResponseTiming{Response: response, Timing: report.Timing{
			ConnectStart: time.Unix(0, 5),
			Done:         time.Unix(0, 10),
		}},
		ResponseTiming{Response: response, Timing: report.Timing{
			ConnectStart: time.Unix(0, 2),
			Done:         time.Unix(0, 6),
		}},
	}

	assert.Equal(t, time.Duration(5), responseTimings.GetLongestDuration())
}
