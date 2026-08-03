package cli

import (
	"github.com/PrashantMohite1/docklens/internal/analyzer"
	"github.com/spf13/cobra"
)

var verifycmd = &cobra.Command{
	Use:  "verify [command]",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		imageName := args[0]
		imgpath, _ := cmd.Flags().GetString("path")

		analyzer.Run_in_container(imageName, imgpath)
	},
}

func init() {
	imageCmd.AddCommand(verifycmd)
	verifycmd.Flags().StringP("path", "p", "", "Command to run in container")

}
