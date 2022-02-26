package lode

import (
	"github.com/JamesBalazs/lode/internal/lode/report"
	"github.com/JamesBalazs/lode/internal/types"
	"math"
	"strings"
	"text/template"
	"time"
)

var newTemplate = func(name string) types.TemplateInt {
	return template.New(name)
}

type TestReport struct {
	Target          string
	Concurrency     int
	Duration        time.Duration
	ResponseCount   int
	RequestRate     float64
	ResponseTimings ResponseTimings
}

func NewTestReport(lode *Lode) TestReport {
	duration := lode.FinishTime.Sub(lode.StartTime).Truncate(report.TimingResolution)
	responseCount := len(lode.ResponseTimings)

	return TestReport{
		Target:          strings.Join([]string{lode.Request.Method, lode.Request.URL.String()}, " "),
		Concurrency:     lode.Concurrency,
		Duration:        duration,
		ResponseCount:   responseCount,
		RequestRate:     math.Round((float64(responseCount)/duration.Seconds())*100) / 100,
		ResponseTimings: lode.ResponseTimings,
	}
}

func (t TestReport) StatusHistogram() report.StatusHistogram {
	return report.BuildStatusHistogram(t.ResponseTimings.Responses(), t.ResponseCount)
}

func (t TestReport) LatencyPercentiles() report.LatencyPercentiles {
	return report.BuildLatencyPercentiles(t.ResponseTimings.Timings())
}

func (t TestReport) FirstResponse() ResponseTiming {
	return t.ResponseTimings[0]
}

func (t TestReport) MultipleResponses() bool {
	return t.ResponseCount > 1
}

func (t TestReport) OneResponse() bool {
	return t.ResponseCount == 1
}

func (t TestReport) Output() string {
	templateString := `Target: {{ .Target }}
Concurrency: {{ .Concurrency }}
Requests made: {{ .ResponseCount }}
Time taken: {{ .Duration }}
Requests per second (avg): {{ .RequestRate }}
{{ if .MultipleResponses }}
Response code breakdown:
{{ .StatusHistogram }}
Percentile latency breakdown:
{{ .LatencyPercentiles }}
{{ else if .OneResponse }}
Timing breakdown:
{{ .FirstResponse.Timing }}
{{ else }}
No requests made...
{{ end }}`

	var err error
	tmpl := newTemplate("report")
	tmpl, err = tmpl.Parse(templateString)
	if err != nil {
		Logger.Panicf("Error parsing report template: %s", err.Error())
		return ""
	}
	builder := strings.Builder{}
	err = tmpl.Execute(&builder, t)
	if err != nil {
		Logger.Panicf("Error executing report template: %s", err.Error())
		return ""
	}
	return builder.String()
}
