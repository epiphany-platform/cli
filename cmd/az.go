package cmd

import (
	"github.com/spf13/cobra"
)

// azCmd represents the az command
var azCmd = &cobra.Command{
	Use:   "az",
	Short: "Commands used to work with Azure cloud.",
	Long:  `Commands used to work with Azure cloud.`,
	PreRun: func(cmd *cobra.Command, args []string) {
		debug("az pre run called")
	},
}

func init() {
	rootCmd.AddCommand(azCmd)
}
