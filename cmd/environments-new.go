package cmd

import (
	"github.com/epiphany-platform/cli/internal/logger"
	"github.com/epiphany-platform/cli/pkg/promptui"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var newEnvName string

// envNewCmd represents the new command
var envNewCmd = &cobra.Command{
	Use:   "new",
	Short: "Creates new environment",
	Long:  `TODO`,
	PreRun: func(cmd *cobra.Command, args []string) {
		logger.Debug().Msg("environments new called")

		err := viper.BindPFlags(cmd.Flags())
		if err != nil {
			logger.Fatal().Err(err).Msg("BindPFlags failed")
		}

		newEnvName = viper.GetString("name")
	},
	Run: func(cmd *cobra.Command, args []string) {
		if newEnvName == "" {
			if len(args) == 1 {
				newEnvName = args[0]
			}
		}
		if newEnvName == "" {
			ne, err := promptui.PromptForString("Environment name")
			if err != nil {
				logger.Fatal().Err(err).Msg("prompt failed")
			}
			newEnvName = ne
		}

		logger.Debug().Msgf("new environment name is: %s", newEnvName)

		envId, err := config.CreateNewEnvironment(newEnvName)
		if err != nil {
			logger.Fatal().Err(err).Msg("create new environment failed")
		} else {
			logger.Info().Msgf("Created an environment with id %s", envId.String())
		}
	},
}

func init() {
	envCmd.AddCommand(envNewCmd)

	envNewCmd.Flags().String("name", "", "name of new environment to create")
}
