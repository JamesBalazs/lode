package mocks

import (
	"github.com/stretchr/testify/mock"
)

type Log struct {
	mock.Mock
}

func (l *Log) Println(v ...interface{}) {
	l.Called(v)
}

func (l *Log) Printf(str string, v ...interface{}) {
	strings := append([]interface{}{str}, v...)
	l.Called(strings...)
}

func (l *Log) Panicf(str string, v ...interface{}) {
	strings := append([]interface{}{str}, v...)
	l.Called(strings...)
}
