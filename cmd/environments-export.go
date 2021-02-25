package cmd

import (
	"os"

	"github.com/epiphany-platform/cli/pkg/configuration"
	"github.com/epiphany-platform/cli/pkg/environment"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	envIdStr string
	dstDir   string
)

var environmentsExportCmd = &cobra.Command{
	Use:        "export",
	SuggestFor: []string{"expor", "exprt"},
	Short:      "Exports an environment as a zip archive",
	Long: `"export" command allows exporting any environment 
as a zip archive into the specified directory 
or into the current working directory by default`,
	Example: `Export current environment into current working directory: e environments export
Export environment into home directory: e environments export --id ba03a2ba-8fa0-4c15-ac07-894af3dbb364 --destination ~`,

	PreRun: func(cmd *cobra.Command, args []string) {
		logger.Debug().Msg("environments export called")

		err := viper.BindPFlags(cmd.Flags())
		if err != nil {
			logger.Fatal().Err(err).Msg("Command flags are specified incorrectly")
		}

		envIdStr = viper.GetString("id")
		dstDir = viper.GetString("destination")
	},

	Run: func(cmd *cobra.Command, args []string) {

		// Check if destination directory is default
		if dstDir == "" {
			path, err := os.Getwd()
			if err != nil {
				logger.Fatal().Err(err).Msg("Unable to get working directory")
			}
			dstDir = path
		}

		var envId uuid.UUID

		// Default environment and destination directory are current ones
		// Check if environment is default
		if envIdStr == "" {
			config, err := configuration.GetConfig()
			if err != nil {
				logger.Fatal().Err(err).Msg("Unable to get environment config")
			}
			if config.CurrentEnvironment == uuid.Nil {
				logger.Fatal().Msg("Environment has to be selected if id is not specified")
			} else {
				envId = config.CurrentEnvironment
			}
		} else {
			envId = uuid.MustParse(envIdStr)
			// Check if passed environment id is valid
			isEnvValid, err := environment.IsExisting(envId)
			if err != nil {
				logger.Fatal().Err(err).Msgf("Environment %s validation failed", envId.String())
			} else if !isEnvValid {
				logger.Fatal().Msgf("Environment %s is not found", envId.String())
			}
		}

		// Export an environment
		env, err := environment.Get(envId)
		if err != nil {
			logger.Fatal().Err(err).Msgf("Unable to get an environment by id %s", envId.String())
		}

		err = env.Export(dstDir)
		if err != nil {
			logger.Fatal().Err(err).Msgf("Unable to export environment with id %s", envId.String())
		}

		logger.Info().Msgf("Environment with id %s was exported", envId.String())
	},
}

func init() {
	environmentsCmd.AddCommand(environmentsExportCmd)

	environmentsExportCmd.Flags().StringP("id", "i", "", "id of the environment to export, default is current environment")
	environmentsExportCmd.Flags().StringP("destination", "d", "", "destination directory to store exported archive, default is current directory")
	environmentsExportCmd.MarkFlagDirname("destination")
}
