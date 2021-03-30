package cmd

import (
	"errors"
	"fmt"
	"github.com/epiphany-platform/cli/pkg/configuration"
	"github.com/epiphany-platform/cli/pkg/environment"
	"github.com/epiphany-platform/cli/pkg/processor"

	"github.com/spf13/cobra"
)

// envRunCmd represents the run command
var envRunCmd = &cobra.Command{ //TODO consider what are options to create integration tests here. For me it seams that it would be testing of docker
	Use:   "run",
	Short: "Runs installed component command in environment",
	Long:  `TODO`,
	PreRun: func(cmd *cobra.Command, args []string) {
		debug("environments run called")
	},
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 2 {
			config, err := configuration.GetConfig()
			if err != nil {
				logger.Fatal().Err(err).Msg("get config failed")
			}
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
		} else {
			logger.
				Fatal().
				Err(errors.New(fmt.Sprintf("found %d args", len(args)))).
				Msg("incorrect number of arguments")
		}
	},
}

func init() {
	envCmd.AddCommand(envRunCmd)
}
