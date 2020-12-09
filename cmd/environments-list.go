package cmd

import (
	"fmt"
	"github.com/epiphany-platform/cli/pkg/configuration"
	"github.com/epiphany-platform/cli/pkg/environment"

	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "TODO",
	Long:  `TODO`,
	PreRun: func(cmd *cobra.Command, args []string) {
		debug("environments list pre run called")
	},
	Run: func(cmd *cobra.Command, args []string) {
		debug("list called")
		config, err := configuration.GetConfig()
		if err != nil {
			errGetConfig(err)
		}
		environments, err := environment.GetAll()
		if err != nil {
			errGetEnvironments(err)
		}
		for _, e := range environments {
			if e.Uuid.String() == config.CurrentEnvironment.String() {
				fmt.Printf("* %s\n", e.Name)
			} else {
				fmt.Printf("  %s\n", e.Name)
			}
		}
	},
}

func init() {
	environmentsCmd.AddCommand(listCmd)
}
