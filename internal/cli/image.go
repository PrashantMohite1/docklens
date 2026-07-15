package cli

import "github.com/spf13/cobra"

var imageCmd = &cobra.Command{
	Use:   "image",
	Short: "Image related commands",
}

func init() {
	rootCmd.AddCommand(imageCmd)
}
