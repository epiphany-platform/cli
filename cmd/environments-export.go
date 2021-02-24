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
	envIDStr string
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

		var envID uuid.UUID

		// Default environment and destination directory are current ones
		// Check if environment is default
		if envIDStr == "" {
			config, err := configuration.GetConfig()
			if err != nil {
				logger.Fatal().Err(err).Msg("Unable to get environment config")
			}
			envID = config.CurrentEnvironment
		} else {
			envID = uuid.MustParse(envIDStr)
			// Check if passed environment id is valid
			isEnvValid, err := environment.IsValid(envID)
			if err != nil {
				logger.Fatal().Err(err).Msgf("Environment %s validation failed", envID.String())
			} else if !isEnvValid {
				logger.Fatal().Msgf("Environment %s is not found", envID.String())
			}
		}

		// Export an environment
		env, err := environment.Get(envID)
		if err != nil {
			logger.Fatal().Err(err).Msgf("Unable to get an environment by id %s", envID.String())
		}

		err = env.Export(dstDir)
		if err != nil {
			logger.Fatal().Err(err).Msgf("Unable to export environment with id %s", envID.String())
		}

		logger.Info().Msgf("Environment with id %s was exported", envID.String())
	},
}

func init() {
	environmentsCmd.AddCommand(environmentsExportCmd)

	environmentsExportCmd.Flags().StringP("id", "i", "", "id of the environment to export, default is current environment")
	environmentsExportCmd.Flags().StringP("destination", "d", "", "destination directory to store exported archive, default is current directory")
	environmentsExportCmd.MarkFlagDirname("destination")
}
