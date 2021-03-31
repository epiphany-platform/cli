package cmd

import (
	"github.com/epiphany-platform/cli/internal/logger"

	"github.com/spf13/cobra"
)

// sshKeygenCmd represents the keygen command
var sshKeygenCmd = &cobra.Command{
	Use:   "keygen",
	Short: "Commands related to ssh keygen operations.",
	Long:  `Commands related to ssh keygen operations.`,
	PreRun: func(cmd *cobra.Command, args []string) {
		logger.Debug().Msg("keygen called")
	},
}

func init() {
	sshCmd.AddCommand(sshKeygenCmd)
}
