package mocks

import (
	"fmt"
	"github.com/stretchr/testify/mock"
)

type Log struct {
	mock.Mock
}

func (l *Log) Println(v ...interface{}) {
	str := fmt.Sprint(v...)
	l.Called(str)
}

func (l *Log) Printf(str string, v ...interface{}) {
	strings := append([]interface{}{str}, v...)
	l.Called(strings...)
}

func (l *Log) Panicln(v ...interface{}) {
	str := fmt.Sprint(v...)
	l.Called(str)
}

func (l *Log) Panicf(str string, v ...interface{}) {
	strings := append([]interface{}{str}, v...)
	l.Called(strings...)
}
