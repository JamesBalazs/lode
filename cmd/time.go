/*
Copyright © 2022 James Balazs <j.c.balazs1@gmail.com>

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
	"net/http"
	"time"

	"github.com/spf13/cobra"
)

// timeCmd represents the time command
var timeCmd = &cobra.Command{
	Use:   "time",
	Short: "Make a single request",
	Long: `Sends a single request and prints a handy timing breakdown.

e.g. lode time --timeout 3s -m GET https://example.com`,
	Run: func(cmd *cobra.Command, args []string) {
		body := files.ReaderFromFileOrString(file, body)
		client := &http.Client{Timeout: timeout}
		lode := lode.New(args[0], method, time.Duration(1), client, 1, 1, 0, body)

		defer lode.Report()
		lode.Run()
	},
}

func init() {
	rootCmd.AddCommand(timeCmd)

	timeCmd.Flags().StringVarP(&method, "method", "m", "GET", "HTTP method to use - defaults to GET")
	timeCmd.Flags().DurationVarP(&timeout, "timeout", "t", 5*time.Second, "Timeout per request, e.g. 200ms or 1s - defaults to 5s")
	timeCmd.Flags().StringVarP(&body, "body", "b", "", "POST/PUT body")
	timeCmd.Flags().StringVarP(&file, "file", "F", "", "POST/PUT body filepath")

}
