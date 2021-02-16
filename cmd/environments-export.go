package cmd

import (
	"os"

	"github.com/epiphany-platform/cli/pkg/configuration"
	"github.com/epiphany-platform/cli/pkg/environment"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	envID  string
	dstDir string
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

		envID = viper.GetString("id")
		dstDir = viper.GetString("destination")
	},

	Run: func(cmd *cobra.Command, args []string) {
		isCurrentEnvUsed := false

		// Default environment and destination directory are current ones
		// Check if environment is default
		if envID == "" {
			isCurrentEnvUsed = true
			config, err := configuration.GetConfig()
			if err != nil {
				logger.Fatal().Err(err).Msg("Unable to get environment config")
			}
			envID = config.CurrentEnvironment.String()
		}

		// Check if destination directory is default
		if dstDir == "" {
			path, err := os.Getwd()
			if err != nil {
				logger.Fatal().Err(err).Msg("Unable to get working directory")
			}
			dstDir = path
		}

		// Check if passed environment id is valid
		if !isCurrentEnvUsed {
			isEnvValid, err := environment.IsValid(envID)
			if err != nil {
				logger.Fatal().Err(err).Msgf("Environment %s validation failed", envID)
			} else if !isEnvValid {
				logger.Fatal().Msgf("Environment %s is not found", envID)
			}
		}

		// Export an environment
		err := environment.Export(envID, dstDir)
		if err != nil {
			logger.Fatal().Err(err).Msgf("Unable to export environment with id %s", envID)
		}
	},
}

func init() {
	environmentsCmd.AddCommand(environmentsExportCmd)

	environmentsExportCmd.Flags().StringP("id", "i", "", "id of the environment to export, default is current environment")
	environmentsExportCmd.Flags().StringP("destination", "d", "", "destination directory to store exported archive, default is current directory")
	environmentsExportCmd.MarkFlagDirname("destination")
}
