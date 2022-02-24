package lode

import (
	"strings"
	"time"
)

type Params struct {
	Url         string
	Method      string
	Body        string
	File        string
	Freq        int
	Concurrency int
	MaxRequests int
	Delay       time.Duration
	Timeout     time.Duration
	MaxTime     time.Duration
	Headers     []string
}

func (p Params) Validate() {
	var errors []string

	if p.Url == "" {
		errors = append(errors, "url must be provided")
	}
	if p.Method == "" {
		errors = append(errors, "method must be provided")
	}
	if p.Freq == 0 && p.Delay == 0 {
		errors = append(errors, "freq or delay must be provided")
	}
	if p.Concurrency < 1 {
		errors = append(errors, "concurrency must be provided as a positive integer")
	}
	if p.Timeout == 0 {
		errors = append(errors, "timeout must be provided")
	}
	if p.MaxRequests == 0 && p.MaxTime == 0 {
		errors = append(errors, "maxrequests or maxtime must be provided")
	}
	if len(errors) != 0 {
		Logger.Panicf("Invalid test suite:\n%s\n", strings.Join(errors, "\n"))
	}
}
