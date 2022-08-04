package lode

import (
	report2 "github.com/JamesBalazs/lode/internal/report"
	"github.com/JamesBalazs/lode/internal/responseTimings"
	"github.com/JamesBalazs/lode/internal/types"
	"github.com/manifoldco/promptui"
	"log"
	"math"
	"os"
	"strings"
	"text/template"
	"time"
)

var reportTemplate = &promptui.SelectTemplates{
	Label:    "{{ . }}?",
	Active:   "\U0000276F {{ .Response.Status | cyan }} (Duration {{ .Timing.TotalDuration | red }})",
	Inactive: "  {{ .Response.Status | cyan }} (Duration {{ .Timing.TotalDuration | red }})",
	Details: `
Request details:
{{ "Status:" | faint }}	{{ .Response.Status }}
{{ "Code:" | faint }}	{{ .Response.StatusCode }}
{{ "Timing breakdown:" | faint }}
{{ .Timing.String }}

Request headers:
{{ .Response.Header }}
Request body:
{{ .Response.Body }}`,
}

var newInteractivePrompt = func(label string, responseTimings responseTimings.ResponseTimings) types.PromptSelectInt {
	return &promptui.Select{
		Label:     label,
		Items:     responseTimings,
		Templates: reportTemplate,
	}
}

var newTemplate = func(name string) types.TemplateInt {
	return template.New(name)
}

var newFileLogger = func(path string) *log.Logger {
	file, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	return log.New(file, "", 0)
}

type TestReport struct {
	Target          string
	Concurrency     int
	Duration        time.Duration
	ResponseCount   int
	RequestRate     float64
	ResponseTimings responseTimings.ResponseTimings
	Interactive     bool
}

func NewTestReport(lode *Lode) TestReport {
	duration := lode.FinishTime.Sub(lode.StartTime).Truncate(responseTimings.TimingResolution)
	responseCount := len(lode.ResponseTimings)

	return TestReport{
		Target:          strings.Join([]string{lode.Request.Method, lode.Request.URL.String()}, " "),
		Concurrency:     lode.Concurrency,
		Duration:        duration,
		ResponseCount:   responseCount,
		RequestRate:     math.Round((float64(responseCount)/duration.Seconds())*100) / 100,
		ResponseTimings: lode.ResponseTimings,
		Interactive:     lode.Interactive,
	}
}

func (t TestReport) StatusHistogram() report2.StatusHistogram {
	return report2.BuildStatusHistogram(t.ResponseTimings.Responses(), t.ResponseCount)
}

func (t TestReport) LatencyPercentiles() report2.LatencyPercentiles {
	return report2.BuildLatencyPercentiles(t.ResponseTimings.Timings())
}

func (t TestReport) FirstResponse() responseTimings.ResponseTiming {
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
{{ if or .MultipleResponses .Interactive }}
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
