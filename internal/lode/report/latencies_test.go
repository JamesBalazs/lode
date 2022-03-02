package report

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestBuildLatencyPercentiles(t *testing.T) {
	timings := []*Timing{
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
