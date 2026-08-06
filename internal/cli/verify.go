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
		file, _ := cmd.Flags().GetString("file")

		dir, _ := cmd.Flags().GetString("dir")

		if file != "" {
			analyzer.Verify_file_sha256_in_container(imageName, file)
		}

		if dir != "" {
			analyzer.Verify_dir_in_container(imageName, dir)
		}

	},
}

func init() {
	imageCmd.AddCommand(verifycmd)
	verifycmd.Flags().StringP("file", "f", "", "Path to file for SHA256 hash verification")
	verifycmd.Flags().StringP("dir", "d", "", "Directory patch local : container ")

}
