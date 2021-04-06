package cmd

import (
	"fmt"

	"github.com/epiphany-platform/cli/internal/logger"
	"github.com/epiphany-platform/cli/internal/repository"
	"github.com/spf13/cobra"
)

// repoListCmd represents the list command
var repoListCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists installed repositories",
	PreRun: func(cmd *cobra.Command, args []string) {
		logger.Debug().Msg("list called")
	},
	Run: func(cmd *cobra.Command, args []string) {
		s, err := repository.List()
		if err != nil {
			logger.Fatal().Err(err).Msg("list failed")
		}
		fmt.Print(s)
	},
}

func init() {
	reposCmd.AddCommand(repoListCmd)
}
