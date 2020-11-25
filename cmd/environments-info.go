package cmd

import (
	"fmt"

	"github.com/epiphany-platform/cli/pkg/configuration"
	"github.com/epiphany-platform/cli/pkg/environment"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

// environmentsInfoCmd represents the info command
var environmentsInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Displays information about currently selected environment",
	Long:  `TODO`,
	PreRun: func(cmd *cobra.Command, args []string) {
		debug("environments info called")
	},
	Run: func(cmd *cobra.Command, args []string) {
		config, err := configuration.GetConfig()
		if err != nil {
			errGetConfig(err)
		}
		if config.CurrentEnvironment == uuid.Nil {
			errNilEnvironment()
		}
		environment, err := environment.Get(config.CurrentEnvironment)
		if err != nil {
			errGetEnvironmentDetails(err)
		}
		fmt.Print(environment.String())
	},
}

func init() {
	environmentsCmd.AddCommand(environmentsInfoCmd)
}
