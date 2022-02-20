package lode

import (
	"github.com/JamesBalazs/lode/internal/lode/report"
	"net/http"
)

type ResponseTiming struct {
	Response http.Response
	Timing   report.Timing
}

type ResponseTimings []ResponseTiming

func (r ResponseTimings) Responses() (responses []http.Response) {
	for _, responseTiming := range r {
		responses = append(responses, responseTiming.Response)
	}
	return
}

func (r ResponseTimings) Timings() (timings []report.Timing) {
	for _, responseTiming := range r {
		timings = append(timings, responseTiming.Timing)
	}
	return
}
