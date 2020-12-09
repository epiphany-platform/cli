package cmd

import (
	"github.com/epiphany-platform/cli/pkg/configuration"
	"github.com/epiphany-platform/cli/pkg/promptui"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

// environmentsUseCmd represents the use command
var environmentsUseCmd = &cobra.Command{
	Use:   "use",
	Short: "Allows to select environment to be used",
	Long:  `TODO`,
	PreRun: func(cmd *cobra.Command, args []string) {
		debug("environments use called")
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
	environmentsCmd.AddCommand(environmentsUseCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// environmentsUseCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// environmentsUseCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
