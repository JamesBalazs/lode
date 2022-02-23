package lode

import (
	"github.com/JamesBalazs/lode/internal/files"
	"gopkg.in/yaml.v3"
	"log"
)

type Suite struct {
	Tests []Params
	lodes []LodeInt
}

func SuiteFromFile(path string) Suite {
	reader := files.Open(path)
	r := yaml.NewDecoder(reader)
	var suite Suite
	if err := r.Decode(&suite); err != nil {
		log.Panicf("Error unmarshalling yaml: %s", err.Error())
	}

	for _, params := range suite.Tests {
		suite.lodes = append(suite.lodes, New(params))
	}

	return suite
}

func (s *Suite) Run() {
	for _, lode := range s.lodes {
		lode.Run()
		lode.Report()
	}
}
