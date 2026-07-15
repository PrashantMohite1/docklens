package cli

import (
	"github.com/PrashantMohite1/docklens/internal/analyzer"
	"github.com/spf13/cobra"
)

var analyzeCmd = &cobra.Command{
	Use:   "analyze IMAGE",
	Short: "Analyze a Docker image",
	Args:  cobra.ExactArgs(1),

	Run: func(cmd *cobra.Command, args []string) {

		imageName := args[0]
		analyzer.AnalyzeImage(imageName)

	},
}

func init() {
	imageCmd.AddCommand(analyzeCmd)
}
