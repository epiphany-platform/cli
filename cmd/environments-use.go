package cmd

import (
	"github.com/epiphany-platform/cli/internal/logger"
	"github.com/epiphany-platform/cli/pkg/configuration"
	"github.com/epiphany-platform/cli/pkg/promptui"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

// envUseCmd represents the use command
var envUseCmd = &cobra.Command{
	Use:   "use",
	Short: "Allows to select environment to be used",
	Long:  `TODO`,
	PreRun: func(cmd *cobra.Command, args []string) {
		logger.Debug().Msg("environments use called")
	},
	Run: func(cmd *cobra.Command, args []string) {
		config, err := configuration.GetConfig()
		if err != nil {
			logger.Fatal().Err(err).Msg("get config failed")
		}
		if config.CurrentEnvironment == uuid.Nil {
			logger.Fatal().Msg("no environment selected")
		}

		var u uuid.UUID
		if len(args) == 1 {
			u = uuid.MustParse(args[0])
		} else {
			u, err = promptui.PromptForEnvironmentSelect("Environments")
			if err != nil {
				logger.Fatal().Err(err).Msg("prompt failed")
			}
		}

		logger.Info().Msgf("Chosen environment UUID is %s", u.String())
		err = config.SetUsedEnvironment(u)
		if err != nil {
			logger.Fatal().Err(err).Msg("setting used environment failed")
		}
	},
}

func init() {
	envCmd.AddCommand(envUseCmd)
}
