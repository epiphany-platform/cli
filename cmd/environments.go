package cmd

import (
	"github.com/spf13/cobra"
)

// envCmd represents the environments command
var envCmd = &cobra.Command{
	Use:     "environments",
	Short:   "Allows various interactions with environments",
	Long:    `TODO`,
	Aliases: []string{"env"},
	PreRun: func(cmd *cobra.Command, args []string) {
		debug("environments called")
	},
}

func init() {
	rootCmd.AddCommand(envCmd)
}
