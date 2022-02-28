package lode

import (
	"github.com/JamesBalazs/lode/internal/files"
	"github.com/JamesBalazs/lode/internal/lode/mocks"
	"github.com/JamesBalazs/lode/internal/types"
	"github.com/stretchr/testify/assert"
	"io"
	"strings"
	"testing"
)

func TestSuiteFromFile(t *testing.T) {
	assert := assert.New(t)
	oldOpen := files.Open
	defer func() { files.Open = oldOpen }()
	files.Open = func(name string) (reader io.Reader) {
		return strings.NewReader(`tests:
  - url: https://www.google.co.uk
    method: GET
    concurrency: 4
    freq: 10
    maxrequests: 20
  - url: https://abc.xyz/
    method: GET
    concurrency: 2
    delay: 0.5s
    maxrequests: 4
    headers:
      - SomeHeader=someValue
      - OtherHeader=otherValue
`)
	}

	suite := SuiteFromFile("path")

	assert.Equal(2, len(suite.Tests))
	assert.Equal("https://www.google.co.uk", suite.Tests[0].Url)
	assert.Equal("https://abc.xyz/", suite.Tests[1].Url)
	assert.Equal("SomeHeader=someValue", suite.Tests[1].Headers[0])
	assert.Equal("OtherHeader=otherValue", suite.Tests[1].Headers[1])
}

func TestSuite_Run(t *testing.T) {
	lode1 := &mocks.Lode{}
	lode2 := &mocks.Lode{}
	suite := Suite{
		lodes: []types.LodeInt{
			lode1,
			lode2,
		},
	}
	lode1.On("Run").Return().Once()
	lode1.On("Report").Return().Once()
	lode2.On("Run").Return().Once()
	lode2.On("Report").Return().Once()

	suite.Run()

	lode1.AssertExpectations(t)
	lode2.AssertExpectations(t)
}
