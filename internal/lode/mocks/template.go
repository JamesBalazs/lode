package mocks

import (
	"github.com/stretchr/testify/mock"
	"io"
	"text/template"
)

type Template struct {
	mock.Mock
}

func (l *Template) Parse(text string) (*template.Template, error) {
	args := l.Called(text)
	return args.Get(0).(*template.Template), args.Error(1)
}

func (l *Template) Execute(wr io.Writer, data interface{}) error {
	args := l.Called(wr, data)
	return args.Error(0)
}
