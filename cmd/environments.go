package cmd

import (
	"github.com/spf13/cobra"
)

// environmentsCmd represents the environments command
var environmentsCmd = &cobra.Command{
	Use:   "environments",
	Short: "Allows various interactions with environments",
	Long:  `TODO`,
	PreRun: func(cmd *cobra.Command, args []string) {
		debug("environments called")
	},
}

func init() {
	rootCmd.AddCommand(environmentsCmd)
}
