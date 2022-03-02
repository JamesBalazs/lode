package responseTimings

import (
	"fmt"
	"github.com/manifoldco/promptui"
	"net/http"
	"sort"
	"time"
)

type ResponseTiming struct {
	Response *Response
	Timing   *Timing
}

type ResponseTimings []ResponseTiming

func (r ResponseTimings) Responses() (responses []*Response) {
	for _, responseTiming := range r {
		responses = append(responses, responseTiming.Response)
	}
	return
}

func (r ResponseTimings) Timings() (timings []*Timing) {
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

type Response struct {
	Status        string // e.g. "200 OK"
	StatusCode    int    // e.g. 200
	ContentLength int64
	Header        Header
	Body          string
}

type Header struct {
	HttpHeader http.Header
}

func (header Header) String() (result string) {
	keys := make([]string, len(header.HttpHeader))
	for key := range header.HttpHeader {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for i := range keys {
		key := keys[i]
		if key != "" {
			headerName := promptui.Styler(promptui.FGCyan)(keys[i])
			headerValue := header.HttpHeader.Get(keys[i])
			result += fmt.Sprintf("%s: %s\n", headerName, headerValue)
		}
	}
	return
}
