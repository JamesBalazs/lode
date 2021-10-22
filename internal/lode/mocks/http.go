package mocks

import (
	"github.com/stretchr/testify/mock"
	"net/http"
)

type Client struct {
	mock.Mock
}

func (c *Client) Do(req *http.Request) (http.Response, error) {
  args := c.Called(req)
  return args.Get(0).(http.Response), args.Error(1)
}

