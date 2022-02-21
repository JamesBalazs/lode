package files

import (
	"io"
	"log"
	"os"
	"strings"
)

func ReaderFromFileOrString(file string, body string) (reader io.Reader) {
	if len(file) > 0 {
		var err error
		reader, err = os.Open(file)
		if err != nil {
			log.Panicf("Error reading body from file: %s", err.Error())
		}
	} else {
		reader = strings.NewReader(body)
	}
	return
}
