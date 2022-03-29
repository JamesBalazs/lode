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
	"github.com/JamesBalazs/lode/internal/lode"
	"github.com/spf13/cobra"
	"time"
)

var interactive bool

// testCmd represents the test command
var testCmd = &cobra.Command{
	Use:   "test [url]",
	Short: "Run a single load test",
	Long: `Run lode against a single URL

Supports either --delay or --freq for timing.
e.g. lode test --freq 20 https://example.com`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		params.Url = args[0]
		lode := lode.New(params)
		lode.Interactive = interactive
		defer lode.ExitWithCode()
		defer lode.Report()
		lode.Run()
	},
}

func init() {
	rootCmd.AddCommand(testCmd)

	testCmd.Flags().IntVarP(&params.Freq, "freq", "f", 0, "Number of requests to make per second")
	testCmd.Flags().DurationVarP(&params.Delay, "delay", "d", 1*time.Second, "Time to wait between requests, e.g. 200ms or 1s - defaults to 1s unless --freq specified")
	testCmd.Flags().IntVarP(&params.Concurrency, "concurrency", "c", 1, "Maximum number of concurrent requests")
	testCmd.Flags().IntVarP(&params.MaxRequests, "maxRequests", "n", 0, "Maximum number of requests to make - defaults to 0s (unlimited)")
	testCmd.Flags().DurationVarP(&params.MaxTime, "maxTime", "l", 0*time.Second, "Length of time to make requests, e.g. 20s or 1h - defaults to 0s (unlimited)")

	testCmd.Flags().StringVarP(&params.Method, "method", "m", "GET", "HTTP method to use - defaults to GET")
	testCmd.Flags().DurationVarP(&params.Timeout, "timeout", "t", 5*time.Second, "Timeout per request, e.g. 200ms or 1s - defaults to 5s")
	testCmd.Flags().StringVarP(&params.Body, "body", "b", "", "POST/PUT body")
	testCmd.Flags().StringVarP(&params.File, "file", "F", "", "POST/PUT body filepath")
	testCmd.Flags().StringSliceVarP(&params.Headers, "header", "H", []string{}, "Request headers, in the form X-SomeHeader=value - separate headers with commas, or repeat the flag to add multiple headers")

	testCmd.Flags().BoolVar(&params.FailFast, "fail-fast", false, "Abort the test immediately if a non-success status code is received")
	testCmd.Flags().BoolVar(&params.IgnoreFailures, "ignore-failures", false, "Don't return non-zero exit code when non-success status codes are received")
	testCmd.Flags().BoolVarP(&interactive, "interactive", "i", false, "Interactive list of responses and timing data")
}
