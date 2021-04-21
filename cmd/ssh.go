package cmd

import (
	"github.com/epiphany-platform/cli/internal/logger"

	"github.com/spf13/cobra"
)

// sshCmd represents the ssh command
var sshCmd = &cobra.Command{
	Use:   "ssh",
	Short: "Ssh and keys operations",
	Long:  `Commands related to ssh and keys operations.`,
	PreRun: func(cmd *cobra.Command, args []string) {
		logger.Debug().Msg("ssh called")
	},
}

func init() {
	rootCmd.AddCommand(sshCmd)
}
