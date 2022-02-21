/*
Copyright © 2021 James Balazs <j.c.balazs1@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"github.com/JamesBalazs/lode/internal/files"
	"github.com/JamesBalazs/lode/internal/lode"
	"github.com/spf13/cobra"
	"net/http"
	"time"
)

var method, body, file string
var freq, concurrency, maxRequests int
var delay, timeout, maxTime time.Duration
var headers []string

// testCmd represents the test command
var testCmd = &cobra.Command{
	Use:   "test [url]",
	Short: "Run a single load test",
	Long: `Run lode against a single URL

Supports either --delay or --freq for timing.
e.g. lode test --freq 20 https://example.com`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if freq != 0 {
			delay = time.Second / time.Duration(freq)
		}

		body := files.ReaderFromFileOrString(file, body)
		client := &http.Client{Timeout: timeout}
		lode := lode.New(args[0], method, delay, client, concurrency, maxRequests, maxTime, body, headers)
		defer lode.Report()
		lode.Run()
	},
}

func init() {
	rootCmd.AddCommand(testCmd)

	testCmd.Flags().IntVarP(&freq, "freq", "f", 0, "Number of requests to make per second")
	testCmd.Flags().DurationVarP(&delay, "delay", "d", 1*time.Second, "Time to wait between requests, e.g. 200ms or 1s - defaults to 1s unless --freq specified")
	testCmd.Flags().IntVarP(&concurrency, "concurrency", "c", 1, "Maximum number of concurrent requests")
	testCmd.Flags().IntVarP(&maxRequests, "maxRequests", "n", 0, "Maximum number of requests to make - defaults to 0s (unlimited)")
	testCmd.Flags().DurationVarP(&maxTime, "maxTime", "l", 0*time.Second, "Length of time to make requests, e.g. 20s or 1h - defaults to 0s (unlimited)")

	testCmd.Flags().StringVarP(&method, "method", "m", "GET", "HTTP method to use - defaults to GET")
	testCmd.Flags().DurationVarP(&timeout, "timeout", "t", 5*time.Second, "Timeout per request, e.g. 200ms or 1s - defaults to 5s")
	testCmd.Flags().StringVarP(&body, "body", "b", "", "POST/PUT body")
	testCmd.Flags().StringVarP(&file, "file", "F", "", "POST/PUT body filepath")
	testCmd.Flags().StringSliceVarP(&headers, "header", "H", []string{}, "Request headers, in the form X-SomeHeader=value - separate headers with commas, or repeat the flag to add multiple headers")
}
