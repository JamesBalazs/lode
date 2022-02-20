package report

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

func TestBuildHistogram(t *testing.T) {
	responses := []http.Response{
		{StatusCode: 200},
		{StatusCode: 200},
		{StatusCode: 200},
		{StatusCode: 400},
		{StatusCode: 503},
		{StatusCode: 503},
	}

	histogram := BuildStatusHistogram(responses, len(responses))

	assert.Equal(t, 3, histogram.Data[200])
	assert.Equal(t, 1, histogram.Data[400])
	assert.Equal(t, 2, histogram.Data[503])
}

func TestStatusHistogram_String(t *testing.T) {
	histogram := StatusHistogram{
		TotalCount: 10,
		keys:       []int{200, 503, 400},
		Data: map[int]int{
			200: 7,
			400: 1,
			503: 2,
		},
	}
	assert.Equal(t,
		`200: ==============>       7x
400: ==>                   1x
503: ====>                 2x
`,
		histogram.String())
}

func TestStatusHistogram_Add(t *testing.T) {
	histogram := StatusHistogram{
		TotalCount: 10,
		keys:       []int{200, 503, 400},
		Data: map[int]int{
			200: 7,
			400: 1,
			503: 2,
		},
	}

	histogram.Add(503)
	histogram.Add(500)

	assert.Equal(t, 3, histogram.Data[503])
	assert.Equal(t, 1, histogram.Data[500])
	assert.Equal(t, []int{200, 503, 400, 500}, histogram.keys)
}

func TestBuildLatencyPercentiles(t *testing.T) {
	timings := []Timing{
		{ConnectStart: time.Unix(0, 0), Done: time.Unix(0, 700_000_000)},
		{ConnectStart: time.Unix(0, 0), Done: time.Unix(0, 300_000_000)},
		{ConnectStart: time.Unix(0, 0), Done: time.Unix(0, 500_000_000)},
		{ConnectStart: time.Unix(0, 0), Done: time.Unix(0, 600_000_000)},
	}
	expectedHistogram := LatencyPercentiles{
		Data: map[int]int{
			50:  500,
			66:  550,
			75:  600,
			80:  650,
			90:  650,
			95:  650,
			98:  650,
			99:  650,
			100: 700,
		},
	}

	histogram := BuildLatencyPercentiles(timings)

	assert.Equal(t, expectedHistogram, histogram)
}

func TestLatencyPercentiles_String(t *testing.T) {
	histogram := LatencyPercentiles{
		Data: map[int]int{
			50:  500,
			66:  550,
			75:  600,
			80:  650,
			90:  650,
			95:  650,
			98:  650,
			99:  650,
			100: 700,
		},
	}

	assert.Equal(t,
		`50th: 500ms
66th: 550ms
75th: 600ms
80th: 650ms
90th: 650ms
95th: 650ms
98th: 650ms
99th: 650ms
100th: 700ms
`,
		histogram.String())
}
