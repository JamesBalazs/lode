package lode

import (
	"log"
	"net/http"
	"os"
	"time"
)

var Logger LoggerInt = log.New(os.Stdout, "", log.LstdFlags)
var NewRequest = http.NewRequest

type Lode struct {
	Url string
	Method string
	TargetDelay time.Duration
	Client      HttpClientInt
	Request     *http.Request
}

func NewLode(url string, method string, delay time.Duration, client HttpClientInt) *Lode {
	req, err := NewRequest(method, url, nil)
	if err != nil {
		Logger.Panicf("Error creating request: %s", err.Error())
		return nil
	}

	return &Lode{
		Url: url,
		Method: method,
		TargetDelay: delay,
		Client: client,
		Request: req,
	}
}

func (l *Lode) Run() {
	response, err := l.Client.Do(l.Request)
	if err != nil {
		log.Panicf("Error during request: %s", err.Error())
	}
	log.Printf("Got status: %s\n", response.Status)
}