package types

import (
	"io"
	"net/http"
	"text/template"
)

type Response struct {
	Status        string // e.g. "200 OK"
	StatusCode    int    // e.g. 200
	ContentLength int64
	Header        http.Header
	Body          string
}

type LodeInt interface {
	Run()
	Report()
}

type HttpClientInt interface {
	Do(*http.Request) (*http.Response, error)
}

type LoggerInt interface {
	Println(...interface{})
	Printf(string, ...interface{})
	Panicln(...interface{})
	Panicf(string, ...interface{})
}

type TemplateInt interface {
	Parse(text string) (*template.Template, error)
	Execute(wr io.Writer, data interface{}) error
}
