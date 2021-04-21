package cmd

import (
	"errors"

	"github.com/epiphany-platform/cli/internal/logger"
	"github.com/epiphany-platform/cli/pkg/environment"
	"github.com/epiphany-platform/cli/pkg/promptui"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

var uu uuid.UUID

// envUseCmd represents the use command
var envUseCmd = &cobra.Command{
	Use:   "use",
	Short: "Allows to select environment to be used",
	Long:  `TODO`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) > 1 {
			return errors.New("'use' command gets 1 optional argument with UUID of environment to use")
		}
		if len(args) == 1 {
			u, err := uuid.Parse(args[0])
			if err != nil {
				return err
			}
			exists, err := environment.IsExisting(u)
			if err != nil {
				return err
			}
			if !exists {
				return errors.New("environment with UUID: " + u.String() + " not found")
			}
		}
		return nil
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		if len(args) == 1 {
			logger.Debug().Msg("environments use called")
			uu = uuid.MustParse(args[0])
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			u, err := promptui.PromptForEnvironmentSelect("Environments")
			if err != nil {
				logger.Fatal().Err(err).Msg("prompt failed")
			}
			uu = u
		}

		logger.Info().Msgf("Chosen environment UUID is %s", uu.String())
		err := config.SetUsedEnvironment(uu)
		if err != nil {
			logger.Fatal().Err(err).Msg("setting used environment failed")
		}
	},
}

func init() {
	envCmd.AddCommand(envUseCmd)
}
