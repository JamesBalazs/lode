package mocks

import "github.com/stretchr/testify/mock"

type Select struct {
	mock.Mock
}

func (l *Select) Run() (int, string, error) {
	args := l.Called()
	return args.Get(0).(int), args.Get(1).(string), args.Error(2)
}
