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

	testCmd.Flags().StringVarP(&params.OutFile, "out", "O", "", "Filepath to write requests and timing data, if provided")
	testCmd.Flags().StringVar(&params.OutFormat, "outFormat", "json", "Format to use when writing requests to file - valid options are json and yaml")
}
