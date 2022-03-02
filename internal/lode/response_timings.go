package lode

import (
	"github.com/JamesBalazs/lode/internal/lode/report"
	"github.com/JamesBalazs/lode/internal/types"
	"time"
)

type ResponseTiming struct {
	Response *types.Response
	Timing   *report.Timing
}

type ResponseTimings []ResponseTiming

func (r ResponseTimings) Responses() (responses []*types.Response) {
	for _, responseTiming := range r {
		responses = append(responses, responseTiming.Response)
	}
	return
}

func (r ResponseTimings) Timings() (timings []*report.Timing) {
	for _, responseTiming := range r {
		timings = append(timings, responseTiming.Timing)
	}
	return
}

func (r ResponseTimings) GetLongestDuration() (duration time.Duration) {
	timings := r.Timings()
	duration = timings[0].TotalDuration()
	for _, currentTiming := range timings {
		currentDuration := currentTiming.TotalDuration()
		if currentDuration > duration {
			duration = currentDuration
		}
	}

	return
}
