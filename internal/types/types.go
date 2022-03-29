package types

import (
	"io"
	"net/http"
	"text/template"
)

type LodeInt interface {
	Run()
	Report()
	ExitWithCode()
}

type HttpClientInt interface {
	Do(*http.Request) (*http.Response, error)
}

type LoggerInt interface {
	Println(...interface{})
	Printf(string, ...interface{})
	Panicln(...interface{})
	Panicf(string, ...interface{})
	Fatalf(string, ...interface{})
}

type TemplateInt interface {
	Parse(text string) (*template.Template, error)
	Execute(wr io.Writer, data interface{}) error
}

type PromptSelectInt interface {
	Run() (int, string, error)
}
