package cmd

import (
	"fmt"

	"github.com/epiphany-platform/cli/pkg/configuration"
	"github.com/epiphany-platform/cli/pkg/environment"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

// envInfoCmd represents the info command
var envInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Displays information about currently selected environment",
	Long:  `TODO`,
	PreRun: func(cmd *cobra.Command, args []string) {
		debug("environments info called")
	},
	Run: func(cmd *cobra.Command, args []string) {
		config, err := configuration.GetConfig()
		if err != nil {
			logger.Fatal().Err(err).Msg("get config failed")
		}
		if config.CurrentEnvironment == uuid.Nil {
			logger.Fatal().Msg("no environment selected")
		}
		environment, err := environment.Get(config.CurrentEnvironment)
		if err != nil {
			logger.Fatal().Err(err).Msg("get environments details failed")
		}
		fmt.Print(environment.String())
	},
}

func init() {
	envCmd.AddCommand(envInfoCmd)
}
