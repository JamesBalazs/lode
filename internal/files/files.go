package files

import (
	"io"
	"log"
	"os"
	"strings"
)

var open = func(name string) (reader io.Reader) {
	var err error
	reader, err = os.Open(name)
	if err != nil {
		log.Panicf("Error reading from file: %s", err.Error())
	}
	return
}

func ReaderFromFileOrString(file string, body string) (reader io.Reader) {
	if len(file) > 0 {
		reader = open(file)
	} else {
		reader = strings.NewReader(body)
	}
	return
}
