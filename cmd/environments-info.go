package cmd

import (
	"fmt"

	"github.com/epiphany-platform/cli/internal/logger"
	"github.com/spf13/cobra"
)

// envInfoCmd represents the info command
var envInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Displays information about currently selected environment",
	Long:  `TODO`,
	PreRun: func(cmd *cobra.Command, args []string) {
		logger.Debug().Msg("environments info called")
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print(currentEnvironment.String())
	},
}

func init() {
	envCmd.AddCommand(envInfoCmd)
}
