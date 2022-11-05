package lode

import (
	"encoding/json"
	"github.com/JamesBalazs/lode/internal/files"
	"github.com/JamesBalazs/lode/internal/responseTimings"
	"gopkg.in/yaml.v3"
	"time"
)

type RunDataV1 struct {
	Version         string
	Target          string
	Concurrency     int
	Duration        time.Duration
	ResponseCount   int
	RequestRate     float64
	ResponseTimings responseTimings.ResponseTimings
}

func (runData RunDataV1) ToInteractiveTestReport() TestReport {
	return TestReport{
		Target:          runData.Target,
		Concurrency:     runData.Concurrency,
		Duration:        runData.Duration,
		ResponseCount:   runData.ResponseCount,
		RequestRate:     runData.RequestRate,
		ResponseTimings: runData.ResponseTimings,
		Interactive:     true,
	}
}

func RunDataFromFile(path string, format string) (runData RunDataV1) {
	reader := files.Open(path)

	var decoder files.Decoder
	switch format {
	case "json":
		decoder = json.NewDecoder(reader)
	case "yaml":
		decoder = yaml.NewDecoder(reader)
	default:
		panic("invalid format")
	}

	if err := decoder.Decode(&runData); err != nil {
		panic(err)
	}
	return
}
