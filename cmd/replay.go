package cmd

import (
	"github.com/JamesBalazs/lode/internal/lode"
	"github.com/spf13/cobra"
)

var inFormat string

// replayCmd represents the replay command
var replayCmd = &cobra.Command{
	Use:   "replay",
	Short: "Replay a log file that was written with --out",
	Long: `Load a log file and display the timings, response bodies, headers etc. in interactive form

e.g. lode replay --inFormat yaml ./out.yaml`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		runData := lode.RunDataFromFile(args[0], inFormat)
		report := runData.ToInteractiveTestReport()
		lode.RunReport(report)
	},
}

func init() {
	rootCmd.AddCommand(replayCmd)

	replayCmd.Flags().StringVar(&inFormat, "inFormat", "json", "Format of requests in file - valid options are json and yaml")
}
