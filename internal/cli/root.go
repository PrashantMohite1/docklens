package cli

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:   "docklens",
	Short: "Docker inspection toolkit",
	Long:  "DockLens helps inspect and analyze Docker resources.",
}

func Execute() error {
	return rootCmd.Execute()
}
