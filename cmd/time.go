package cmd

import (
	"github.com/JamesBalazs/lode/internal/lode"
	"time"

	"github.com/spf13/cobra"
)

// timeCmd represents the time command
var timeCmd = &cobra.Command{
	Use:   "time",
	Short: "Make a single request",
	Long: `Sends a single request and prints a handy timing breakdown.

e.g. lode time --timeout 3s -m GET https://example.com`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		params.Url = args[0]
		params.Concurrency = 1
		params.Delay = 1 * time.Second
		params.MaxRequests = 1
		lode := lode.New(params)
		lode.Interactive = interactive
		defer lode.ExitWithCode()
		defer lode.Report()
		lode.Run()
	},
}

func init() {
	rootCmd.AddCommand(timeCmd)

	timeCmd.Flags().StringVarP(&params.Method, "method", "m", "GET", "HTTP method to use - defaults to GET")
	timeCmd.Flags().DurationVarP(&params.Timeout, "timeout", "t", 5*time.Second, "Timeout per request, e.g. 200ms or 1s - defaults to 5s")
	timeCmd.Flags().StringVarP(&params.Body, "body", "b", "", "POST/PUT body")
	timeCmd.Flags().StringVarP(&params.File, "file", "F", "", "POST/PUT body filepath")
	timeCmd.Flags().StringSliceVarP(&params.Headers, "header", "H", []string{}, "Request headers, in the form X-SomeHeader=value - separate headers with commas, or repeat the flag to add multiple headers")

	timeCmd.Flags().BoolVarP(&interactive, "interactive", "i", false, "Interactive list of responses and timing data")
	timeCmd.Flags().BoolVar(&params.IgnoreFailures, "ignore-failures", false, "Don't return non-zero exit code when non-success status codes are received")

	timeCmd.Flags().StringVarP(&params.Outfile, "out", "O", "", "Filepath to write requests and timing data, if provided")
	timeCmd.Flags().StringVar(&params.OutFormat, "outFormat", "json", "Format to use when writing requests to file - valid options are json and yaml")
}
