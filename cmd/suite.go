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
	"github.com/JamesBalazs/lode/internal/lode"

	"github.com/spf13/cobra"
)

var dryRun bool

// suiteCmd represents the suite command
var suiteCmd = &cobra.Command{
	Use:   "suite",
	Short: "Run a test suite from a YAML file",
	Long: `Run a series of predefined load tests from a YAML file, for CI/automation use.

e.g. lode suite examples/suite.yaml

Example YAML format:
tests:
  - url: https://www.google.co.uk
    method: GET
    concurrency: 4
    freq: 10
    maxrequests: 20
  - url: https://abc.xyz/
    method: GET
    concurrency: 2
    delay: 0.5s
    maxrequests: 4
    headers:
      - SomeHeader=someValue
      - OtherHeader=otherValue`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		suite := lode.SuiteFromFile(args[0])
		if !dryRun {
			suite.Run()
		}
	},
}

func init() {
	rootCmd.AddCommand(suiteCmd)

	suiteCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Validate YAML file without running the test suite")
}
