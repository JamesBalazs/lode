package report

import (
	"fmt"
	"github.com/montanaflynn/stats"
	"sort"
)

var latencyPercentiles = []float64{50, 66, 75, 80, 90, 95, 98, 99, 100}

type LatencyPercentiles struct {
	Data map[int]int
}

func BuildLatencyPercentiles(timings []Timing) (histogram LatencyPercentiles) {
	histogram = LatencyPercentiles{Data: make(map[int]int)}
	timingsCount := len(timings)
	durations := make([]float64, timingsCount)
	for i, timing := range timings {
		durations[i] = float64(timing.TotalDuration().Milliseconds())
	}

	for _, percentile := range latencyPercentiles {
		percentileLatency, _ := stats.Percentile(durations, percentile)
		histogram.Data[int(percentile)] = int(percentileLatency)
	}

	return
}

func (t LatencyPercentiles) String() (string string) {
	sort.Float64s(latencyPercentiles)
	for _, percentile := range latencyPercentiles {
		string = string + fmt.Sprintf("%dth: %dms\n", int(percentile), t.Data[int(percentile)])
	}
	return
}
