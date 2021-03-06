package cmd

import (
	"os"

	"github.com/epiphany-platform/cli/internal/logger"
	"github.com/epiphany-platform/cli/pkg/environment"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	envIdStr string
	dstDir   string
)

// envExportCmd represents envs export command
var envExportCmd = &cobra.Command{
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
			if config.CurrentEnvironment == uuid.Nil { // TODO this will never occur due to new janitor environment initialization
				logger.Fatal().Msg("Environment has to be selected if id is not specified")
			} else {
				envId = config.CurrentEnvironment
			}
		} else {
			envId = uuid.MustParse(envIdStr)
			// Check if passed environment id is valid
			exists, err := environment.IsExisting(envId)
			if err != nil {
				logger.Fatal().Err(err).Msgf("Environment existence check failed (environment id: %s)", envId.String())
			} else if !exists {
				logger.Fatal().Msgf("Environment not found (environment id: %s)", envId.String())
			}
		}

		// Export an environment
		env, err := environment.Get(envId)
		if err != nil {
			logger.Fatal().Err(err).Msgf("Unable to get an environment by id (environment id: %s)", envId.String())
		}

		err = env.Export(dstDir)
		if err != nil {
			logger.Fatal().Err(err).Msgf("Unable to export environment (environment id: %s)", envId.String())
		}

		logger.Info().Msgf("Export operation finished correctly (environment id: %s)", envId.String())
	},
}

func init() {
	envCmd.AddCommand(envExportCmd)

	//TODO decide if we need this parameter at all
	envExportCmd.Flags().StringP("id", "i", "", "id of the environment to export, default is current environment")
	envExportCmd.Flags().StringP("destination", "d", "", "destination directory to store exported archive, default is current directory")
	_ = envExportCmd.MarkFlagDirname("destination")
}
