package mocks

import (
	"github.com/stretchr/testify/mock"
)

type Lode struct {
	mock.Mock
}

func (l *Lode) Run() {
	l.Called()
}

func (l *Lode) Report() {
	l.Called()
}

func (l *Lode) ExitWithCode() {
	l.Called()
}
