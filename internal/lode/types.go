package lode

import (
	"net/http"
)

type HttpClientInt interface {
	Do(*http.Request) (*http.Response, error)
}

type LoggerInt interface {
	Println(...interface{})
	Printf(string, ...interface{})
	Panicln(...interface{})
	Panicf(string, ...interface{})
}
