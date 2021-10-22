package mocks

import (
	"github.com/stretchr/testify/mock"
)

type Log struct {
	mock.Mock
}

func (l *Log) Println(v ...interface{}) {
	l.Called(v)
	return
}

func (l *Log) Printf(str string, v ...interface{}) {
	strings := append([]interface{}{str}, v...)
	l.Called(strings...)
	return
}

func (l *Log) Panicf(str string, v ...interface{}) {
	strings := append([]interface{}{str}, v...)
	l.Called(strings...)
	return
}

