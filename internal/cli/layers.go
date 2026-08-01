package cli

import (
	"github.com/PrashantMohite1/docklens/internal/analyzer"
	"github.com/spf13/cobra"
)

var layercmd = &cobra.Command{
	Use: "layers",

	Run: func(cmd *cobra.Command, args []string) {
		imageName := args[0]
		analyzer.Get_img_layer(imageName)

	},
}

func init() {
	imageCmd.AddCommand(layercmd)
}
