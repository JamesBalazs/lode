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
