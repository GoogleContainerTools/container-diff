package cmd

import (
	"fmt"

	"github.com/GoogleCloudPlatform/container-diff/version"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version of container-diff",
	Long:  `Print the version of container-diff.`,
	Run: func(command *cobra.Command, args []string) {
		fmt.Println(version.GetVersion())
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
