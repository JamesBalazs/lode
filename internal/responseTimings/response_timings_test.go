package responseTimings

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

func TestResponseTimings_Responses(t *testing.T) {
	response := &Response{StatusCode: 200}
	timing := &Timing{Done: time.Now()}
	responseTiming := ResponseTiming{Response: response, Timing: timing}
	responseTimings := ResponseTimings{
		responseTiming,
		responseTiming,
	}

	assert.Equal(t, []*Response{response, response}, responseTimings.Responses())
}

func TestResponseTimings_Timings(t *testing.T) {
	response := &Response{StatusCode: 200}
	timing := &Timing{Done: time.Now()}
	responseTiming := ResponseTiming{Response: response, Timing: timing}
	responseTimings := ResponseTimings{
		responseTiming,
		responseTiming,
	}

	assert.Equal(t, []*Timing{timing, timing}, responseTimings.Timings())
}

func TestResponseTimings_GetLongestDuration(t *testing.T) {
	response := &Response{StatusCode: 200}
	responseTimings := ResponseTimings{
		ResponseTiming{Response: response, Timing: &Timing{
			ConnectStart: time.Unix(0, 1_000_000),
			Done:         time.Unix(0, 3_000_000),
		}},
		ResponseTiming{Response: response, Timing: &Timing{
			ConnectStart: time.Unix(0, 5_000_000),
			Done:         time.Unix(0, 10_000_000),
		}},
		ResponseTiming{Response: response, Timing: &Timing{
			ConnectStart: time.Unix(0, 2_000_000),
			Done:         time.Unix(0, 6_000_000),
		}},
	}

	assert.Equal(t, 5*time.Millisecond, responseTimings.GetLongestDuration())
}

func TestHeaderString(t *testing.T) {
	header := Header{HttpHeader: http.Header{
		"Set-Cookie": {`abc="def"`},
		"A-Header":   {`someValue`},
		"B-Header":   {`otherValue`},
	}}

	assert.Equal(t, "\033[36mA-Header\033[0m: someValue\n"+
		"\033[36mB-Header\033[0m: otherValue\n"+
		"\033[36mSet-Cookie\033[0m: abc=\"def\"\n",
		header.String())
}
