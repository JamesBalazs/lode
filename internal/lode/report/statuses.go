package report

import (
	"fmt"
	"net/http"
	"sort"
	"strings"
)

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

func (s StatusHistogram) String() (string string) {
	sort.Ints(s.keys)
	for _, statusCode := range s.keys {
		statusCount := s.Data[statusCode]
		var percentage = float32(statusCount) / float32(s.TotalCount)
		bar := strings.Repeat("=", int(percentage*20)) + ">"
		string = string + fmt.Sprintf("%d: %-21s %dx\n", statusCode, bar, statusCount)
	}
	return
}
