package cmd

import (
	"github.com/epiphany-platform/cli/internal/logger"

	"github.com/spf13/cobra"
)

// reposCmd represents the repos command
var reposCmd = &cobra.Command{
	Use:   "repos",
	Short: "Commands related to repos management",
	PreRun: func(cmd *cobra.Command, args []string) {
		logger.Debug().Msg("repos called")
	},
}

func init() {
	rootCmd.AddCommand(reposCmd)
}
