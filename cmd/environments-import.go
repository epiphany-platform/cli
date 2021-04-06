package cmd

import (
	"os"

	"github.com/epiphany-platform/cli/internal/logger"
	"github.com/epiphany-platform/cli/pkg/environment"
	"github.com/epiphany-platform/cli/pkg/promptui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var srcFile string

// envImportCmd represents envs import command
var envImportCmd = &cobra.Command{
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
			logger.Fatal().Err(err).Msg("Command flags are specified incorrectly")
		}

		srcFile = viper.GetString("from")

	},

	Run: func(cmd *cobra.Command, args []string) {
		// Ask user for source file path if no file to import from is specified
		if srcFile == "" {
			srcFile, _ = promptui.PromptForString("File to import environment from")
		}

		// Check if source file exists
		if _, err := os.Stat(srcFile); err != nil {
			logger.Fatal().Err(err).Msg("Incorrect file path specified")
		}

		// Import environment
		envId, err := environment.Import(srcFile)
		if err != nil {
			logger.Fatal().Err(err).Msg("Unable to import environment from specified file")
		}
		logger.Info().Msgf("Environment with id %s was imported", envId.String())

		// Switch to the imported environment
		err = config.SetUsedEnvironment(envId)
		if err != nil {
			logger.Fatal().Err(err).Msg("Setting used environment failed")
		}
		logger.Info().Msgf("Switched to the imported environment with id %s", envId.String())
	},
}

func init() {
	envCmd.AddCommand(envImportCmd)

	envImportCmd.Flags().StringP("from", "f", "", "File to import from")
	_ = envImportCmd.MarkFlagFilename("from")
}
