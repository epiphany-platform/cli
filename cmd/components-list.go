package cmd

import (
	"fmt"
	"github.com/epiphany-platform/cli/internal/logger"

	"github.com/epiphany-platform/cli/pkg/repository"
	"github.com/spf13/cobra"
)

// componentsListCmd represents the list command
var componentsListCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists all existing components in repository",
	Long:  `TODO`,
	PreRun: func(cmd *cobra.Command, args []string) {
		logger.Debug().Msg("component list called")
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print(repository.GetRepository().ComponentsString())
	},
}

func init() {
	componentsCmd.AddCommand(componentsListCmd)
}
