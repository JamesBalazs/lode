package files

import (
	"github.com/stretchr/testify/assert"
	"io"
	"strings"
	"testing"
)

func TestReaderFromFileOrString(t *testing.T) {
	assert := assert.New(t)
	oldOpen := Open
	defer func() { Open = oldOpen }()
	expectedReader := strings.NewReader("Some body from file")
	Open = func(name string) io.Reader {
		return expectedReader
	}

	reader := ReaderFromFileOrString("some/file/path", "")

	assert.Equal(expectedReader, reader)

	expectedBody := "Some body from string"
	expectedReader = strings.NewReader(expectedBody)

	reader = ReaderFromFileOrString("", expectedBody)

	assert.Equal(expectedReader, reader)
}
