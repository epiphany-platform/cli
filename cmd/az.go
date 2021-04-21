package cmd

import (
	"github.com/epiphany-platform/cli/internal/logger"

	"github.com/spf13/cobra"
)

// azCmd represents the az command
var azCmd = &cobra.Command{
	Use:   "az",
	Short: "Azure Cloud related operations",
	Long:  `Commands used to work with Azure cloud.`,
	PreRun: func(cmd *cobra.Command, args []string) {
		logger.Debug().Msg("az pre run called")
	},
}

func init() {
	rootCmd.AddCommand(azCmd)
}
