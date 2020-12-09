package cmd

import (
	"github.com/spf13/cobra"
)

// spCmd represents the sp command
var spCmd = &cobra.Command{
	Use:   "sp",
	Short: "Commands used to work with Azure Service Principal.",
	Long:  `Commands used to work with Azure Service Principal.`,
	PreRun: func(cmd *cobra.Command, args []string) {
		debug("az sp pre run called")
	},
}

func init() {
	azCmd.AddCommand(spCmd)
}
