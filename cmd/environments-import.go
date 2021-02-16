package cmd

import (
	"os"

	"github.com/epiphany-platform/cli/pkg/configuration"
	"github.com/epiphany-platform/cli/pkg/environment"
	"github.com/epiphany-platform/cli/pkg/promptui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var srcFile string

var environmentsImportCmd = &cobra.Command{
	Use:        "import",
	SuggestFor: []string{"impor", "imprt"},
	Short:      "Imports a zip compressed environment",
	Long: `"import" command allows importing an environment from a zip archive
and immediately switches to the imported environment`,
	Example: "e environments import --from ba03a2ba-8fa0-4c15-ac07-894af3dbb364.zip",

	PreRun: func(cmd *cobra.Command, args []string) {
		logger.Debug().Msg("environments import called")

		err := viper.BindPFlags(cmd.Flags())
		if err != nil {
			logger.Fatal().Err(err)
		}

		srcFile = viper.GetString("from")

	},

	Run: func(cmd *cobra.Command, args []string) {
		// Ask user for source file path if no file to export from is specified
		if srcFile == "" {
			srcFile, _ = promptui.PromptForString("File to export environment from")
			// Check if environment with such id is already in place
			if _, err := os.Stat(srcFile); err != nil {
				logger.Fatal().Err(err).Msg("Incorrect file path specified")
			}
		}

		// Import environment
		envConfig, err := environment.Import(srcFile)
		if err != nil {
			logger.Fatal().Err(err).Msg("Unable to import environment from specified file")
		}

		// Switch to the imported environment
		config, err := configuration.GetConfig()
		if err != nil {
			logger.Fatal().Err(err).Msg("Get config failed")
		}
		err = config.SetUsedEnvironment(envConfig.Uuid)
		if err != nil {
			logger.Fatal().Err(err).Msg("Setting used environment failed")
		}
		logger.Info().Msgf("Switched to the imported environment with id %s", envConfig.Uuid.String())

		// Download all Docker images for installed components
		for _, cmp := range envConfig.Installed {
			err = cmp.Download()
			if err != nil {
				logger.Fatal().Err(err)
			}
		}
	},
}

func init() {
	environmentsCmd.AddCommand(environmentsImportCmd)

	environmentsImportCmd.Flags().StringP("from", "f", "", "File to import from")
	environmentsImportCmd.MarkFlagFilename("from")
}
