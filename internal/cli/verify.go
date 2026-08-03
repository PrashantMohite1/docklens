package cli

import (
	"github.com/PrashantMohite1/docklens/internal/analyzer"
	"github.com/spf13/cobra"
)

var verifycmd = &cobra.Command{
	Use: "verify",

	Run: func(cmd *cobra.Command, args []string) {
		imageName := args[0]
		command := args[1]
		analyzer.Run_in_container(imageName, command)
	},
}

func init() {
	imageCmd.AddCommand(verifycmd)
}
