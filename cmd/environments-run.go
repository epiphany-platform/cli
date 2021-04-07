package cmd

import (
	"errors"

	"github.com/epiphany-platform/cli/internal/logger"
	"github.com/epiphany-platform/cli/pkg/environment"
	"github.com/epiphany-platform/cli/pkg/processor"

	"github.com/spf13/cobra"
)

// envRunCmd represents the run command
var envRunCmd = &cobra.Command{ //TODO consider what are options to create integration tests here. For me it seams that it would be testing of docker
	Use:   "run",
	Short: "Runs installed component command in environment",
	Long:  `TODO`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 2 {
			return errors.New("incorrect number of arguments")
		}
		return nil
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		logger.Debug().Msg("environments run called")
	},
	Run: func(cmd *cobra.Command, args []string) {
		env, err := environment.Get(config.CurrentEnvironment)
		if err != nil {
			logger.Fatal().Err(err).Msg("get environments details failed")
		}
		c, err := env.GetComponentByName(args[0])
		if err != nil {
			logger.Fatal().Err(err).Msg("getting component by name failed")
		}
		err = c.Run(args[1], processor.TemplateProcessor(config, env))
		if err != nil {
			logger.Fatal().Err(err).Msg("run command failed")
		}
		logger.Info().Msgf("running %s %s finished", args[0], args[1])
	},
}

func init() {
	envCmd.AddCommand(envRunCmd)
}
