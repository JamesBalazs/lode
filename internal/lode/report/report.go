package report

import (
	"fmt"
	"github.com/montanaflynn/stats"
	"net/http"
	"sort"
	"strings"
)

var latencyPercentiles = []float64{50, 66, 75, 80, 90, 95, 98, 99, 100}

type StatusHistogram struct {
	Data       map[int]int
	TotalCount int
	keys       []int
}

func BuildStatusHistogram(responses []http.Response, totalResponses int) (histogram StatusHistogram) {
	histogram = StatusHistogram{Data: make(map[int]int)}
	histogram.TotalCount = totalResponses
	for _, response := range responses {
		histogram.Add(response.StatusCode)
	}
	return
}

func (s *StatusHistogram) Add(statusCode int) {
	if s.Data[statusCode] == 0 {
		s.keys = append(s.keys, statusCode)
	}
	s.Data[statusCode]++
}

func (s *StatusHistogram) String() (string string) {
	sort.Ints(s.keys)
	for _, statusCode := range s.keys {
		statusCount := s.Data[statusCode]
		var percentage = float32(statusCount) / float32(s.TotalCount)
		bar := strings.Repeat("=", int(percentage*20)) + ">"
		string = string + fmt.Sprintf("%d: %-21s %dx\n", statusCode, bar, statusCount)
	}
	return
}

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

func (t *LatencyPercentiles) String() (string string) {
	sort.Float64s(latencyPercentiles)
	for _, percentile := range latencyPercentiles {
		string = string + fmt.Sprintf("%dth: %dms\n", int(percentile), t.Data[int(percentile)])
	}
	return
}
