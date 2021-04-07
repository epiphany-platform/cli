package cmd

import (
	"github.com/epiphany-platform/cli/internal/logger"

	"github.com/spf13/cobra"
)

// spCmd represents the sp command
var spCmd = &cobra.Command{
	Use:   "sp",
	Short: "Commands used to work with Azure Service Principal.",
	Long:  `Commands used to work with Azure Service Principal.`,
	PreRun: func(cmd *cobra.Command, args []string) {
		logger.Debug().Msg("az sp pre run called")
	},
}

func init() {
	azCmd.AddCommand(spCmd)
}
