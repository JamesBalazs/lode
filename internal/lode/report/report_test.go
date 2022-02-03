package report

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestBuildHistogram(t *testing.T) {
	responses := []http.Response{
		{StatusCode: 200},
		{StatusCode: 200},
		{StatusCode: 200},
		{StatusCode: 400},
		{StatusCode: 503},
		{StatusCode: 503},
	}

	histogram := BuildStatusHistogram(responses, len(responses))

	assert.Equal(t, 3, histogram.Data[200])
	assert.Equal(t, 1, histogram.Data[400])
	assert.Equal(t, 2, histogram.Data[503])
}

func TestStatusHistogram_String(t *testing.T) {
	histogram := StatusHistogram{
		TotalCount: 10,
		keys:       []int{200, 503, 400},
		Data: map[int]int{
			200: 7,
			400: 1,
			503: 2,
		},
	}
	assert.Equal(t,
		`200: ==============>       7x
400: ==>                   1x
503: ====>                 2x
`,
		histogram.String())
}

func TestStatusHistogram_Add(t *testing.T) {
	histogram := StatusHistogram{
		TotalCount: 10,
		keys:       []int{200, 503, 400},
		Data: map[int]int{
			200: 7,
			400: 1,
			503: 2,
		},
	}

	histogram.Add(503)
	histogram.Add(500)

	assert.Equal(t, 3, histogram.Data[503])
	assert.Equal(t, 1, histogram.Data[500])
	assert.Equal(t, []int{200, 503, 400, 500}, histogram.keys)
}
