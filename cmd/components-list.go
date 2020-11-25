package cmd

import (
	"fmt"

	"github.com/epiphany-platform/cli/pkg/repository"
	"github.com/spf13/cobra"
)

// componentsListCmd represents the list command
var componentsListCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists all existing components in repository",
	Long:  `TODO`,
	PreRun: func(cmd *cobra.Command, args []string) {
		debug("component list called")
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print(repository.GetRepository().ComponentsString())
	},
}

func init() {
	componentsCmd.AddCommand(componentsListCmd)
}
