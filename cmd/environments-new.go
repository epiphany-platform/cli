package cmd

import (
	"github.com/epiphany-platform/cli/pkg/configuration"
	"github.com/epiphany-platform/cli/pkg/promptui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var newEnvName string

// environmentsNewCmd represents the new command
var environmentsNewCmd = &cobra.Command{
	Use:   "new",
	Short: "Creates new environment",
	Long:  `TODO`,
	PreRun: func(cmd *cobra.Command, args []string) {
		debug("environments new called")

		err := viper.BindPFlags(cmd.Flags())
		if err != nil {
			logger.Fatal().Err(err)
		}

		newEnvName = viper.GetString("name")
	},
	Run: func(cmd *cobra.Command, args []string) {
		config, err := configuration.GetConfig()
		if err != nil {
			logger.Fatal().Err(err).Msg("get config failed")
		}

		if newEnvName == "" {
			if len(args) == 1 {
				newEnvName = args[0]
			}
		}
		if newEnvName == "" {
			newEnvName, err = promptui.PromptForString("Environment name")
			if err != nil {
				logger.Fatal().Err(err).Msg("prompt failed")
			}
		}

		debug("new environment name is: %s", newEnvName)

		err = config.CreateNewEnvironment(newEnvName)
		if err != nil {
			logger.Fatal().Err(err).Msg("create new environment failed")
		}
	},
}

func init() {
	environmentsCmd.AddCommand(environmentsNewCmd)

	environmentsNewCmd.Flags().String("name", "", "name of new environment to create")
}
