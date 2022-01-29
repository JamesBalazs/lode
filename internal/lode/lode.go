package lode

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

var Logger LoggerInt = log.New(os.Stdout, "", log.LstdFlags)
var NewRequest = http.NewRequest

type Lode struct {
	TargetDelay   time.Duration
	Client        HttpClientInt
	Request       *http.Request
	Concurrency   int
	MaxRequests   int
	MaxTime       time.Duration
	Responses     []http.Response
	ResponseCount int
}

func New(url string, method string, delay time.Duration, client HttpClientInt, concurrency int, maxRequests int, maxTime time.Duration) *Lode {
	req, err := NewRequest(method, url, nil)
	if err != nil {
		Logger.Panicf("Error creating request: %s", err.Error())
		return nil
	}

	return &Lode{
		TargetDelay: delay,
		Client:      client,
		Request:     req,
		Concurrency: concurrency,
		MaxRequests: maxRequests,
		MaxTime:     maxTime,
	}
}

func (l *Lode) Run() {
	ticker := time.NewTicker(l.TargetDelay)
	trigger := ticker.C
	stop := make(chan struct{})
	defer l.stop(ticker, stop)
	defer l.report()
	result := make(chan http.Response, 1024)

	for i := 0; i < l.Concurrency; i++ {
		go l.work(trigger, stop, result)
	}

	startTime := time.Now()
	endTime := startTime.Add(l.MaxTime).UnixNano()
	checkMaxRequests := l.MaxRequests > 0
	checkMaxTime := l.MaxTime > 0
	for response := range result {
		l.ResponseCount++
		l.Responses = append(l.Responses, response)

		if (checkMaxRequests && l.ResponseCount >= l.MaxRequests) || (checkMaxTime && time.Now().UnixNano() >= endTime) {
			return
		}
	}
}

func (l *Lode) work(trigger <-chan time.Time, stop chan struct{}, result chan http.Response) {
	for {
		select {
		case <-trigger:
			response, err := l.Client.Do(l.Request)
			if err != nil {
				log.Panicf("Error during request: %s", err.Error())
			}
			result <- *response
		case <-stop:
			return
		}
	}
}

func (l *Lode) stop(ticker *time.Ticker, stop chan struct{}) {
	ticker.Stop()
	close(stop)
}

func (l *Lode) report() {
	responseCodeDistribution := map[int]int{}
	for _, response := range l.Responses {
		responseCodeDistribution[response.StatusCode]++
	}
	histogram := ""
	for statusCode, count := range responseCodeDistribution {
		var percentage = float32(count) / float32(l.ResponseCount)
		bar := strings.Repeat("=", int(percentage*20)) + ">"
		histogram = histogram + fmt.Sprintf("%d: %-21s %dx\n", statusCode, bar, count)
	}
	fmt.Printf("Target: %s %s\n", l.Request.Method, l.Request.URL)
	fmt.Printf("Requests made: %d\n", l.ResponseCount)
	fmt.Printf("Response Breakdown:\n%s\n", histogram)
}
