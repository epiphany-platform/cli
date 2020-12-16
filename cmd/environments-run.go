package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"text/template"

	"github.com/epiphany-platform/cli/pkg/configuration"
	"github.com/epiphany-platform/cli/pkg/environment"

	"github.com/spf13/cobra"
)

// environmentsRunCmd represents the run command
var environmentsRunCmd = &cobra.Command{ //TODO consider what are options to create integration tests here. For me it seams that it would be testing of docker
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
			e, err := environment.Get(config.CurrentEnvironment)
			if err != nil {
				logger.Fatal().Err(err).Msg("get environments details failed")
			}
			c, err := e.GetComponentByName(args[0])
			if err != nil {
				logger.Fatal().Err(err).Msg("getting component by name failed")
			}
			err = c.Run(args[1], environmentsProcessor(config))
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
	environmentsCmd.AddCommand(environmentsRunCmd)
}

func environmentsProcessor(config *configuration.Config) func(envs map[string]string) map[string]string {
	return func(envs map[string]string) map[string]string {
		result := make(map[string]string)
		for k, v := range envs {
			if strings.HasPrefix(v, "#") {
				logger.Debug().Msgf("key %s has value %s", k, v)
				parts := strings.Split(v, "#")
				logger.Debug().Msgf("and value has parts: %#v", parts)
				if parts[1] == "Config" {
					t, err := template.New(k).Parse(parts[2])
					if err != nil {
						logger.Error().Err(err)
						break
					}
					logger.Debug().Msgf("parsed template: %#v", t)
					var b bytes.Buffer
					err = t.Execute(&b, config)
					if err != nil {
						logger.Error().Err(err)
						break
					}
					r := b.String()
					if r == "" {
						logger.Error().Err(errors.New("there was no value obtained"))
						break
					}
					logger.Debug().Msgf("result value: %#v", b.String())
					result[k] = b.String()
				}
			} else {
				result[k] = v
			}
		}

		return result
	}
}
