package cmd

import (
	"errors"
	"io/ioutil"
	"os"
	"path"

	//"path/filepath"
	//"strings"
	"github.com/epiphany-platform/cli/pkg/configuration"
	"github.com/epiphany-platform/cli/pkg/environment"
	"github.com/epiphany-platform/cli/pkg/promptui"
	"github.com/epiphany-platform/cli/pkg/util"
	"github.com/google/uuid"
	"github.com/mholt/archiver/v3"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

var (
	srcFile string
)

var environmentsImportCmd = &cobra.Command{
	Use:   "import",
	Short: "Imports an environment from specified archive",
	Long: `"import" command allows importing an environment from archive
and immediately switches to the imported environment `,
	// TODO: extend with https://pkg.go.dev/github.com/spf13/cobra#Command

	PreRun: func(cmd *cobra.Command, args []string) {
		debug("environments import called")

		err := viper.BindPFlags(cmd.Flags())
		if err != nil {
			logger.Fatal().Err(err)
		}

		srcFile = viper.GetString("from")

	},
	Run: func(cmd *cobra.Command, args []string) {
		//baseSrcFileName := filepath.Base(srcFile)
		//envID := strings.TrimSuffix(baseSrcFileName, filepath.Ext(baseSrcFileName))
		//envConfigDir := path.Join(util.GetHomeDirectory(), util.DefaultConfigurationDirectory, util.DefaultEnvironmentsSubdirectory)

		// Ask user for source file path if no file to export from is specified
		if srcFile == "" {
			srcFile, _ = promptui.PromptForString("File to export environment from")
			// Check if environment with such id is already in place
			if _, err := os.Stat(srcFile); err != nil {
				logger.Fatal().Err(err).Msg("Incorrect file path specified")
			}
		}

		// Check if environment config exists in archive before export
		// and verify its content
		var envConfig *environment.Environment
		isSrcFileValid := false
		err := archiver.Walk(srcFile, func(f archiver.File) error {
			if f.Name() == util.DefaultConfigFileName {
				configContent, err := ioutil.ReadAll(f)
				if err != nil {
					return errors.New("Unable to read environment config")
				}
				envConfig = &environment.Environment{}
				err = yaml.Unmarshal(configContent, envConfig)
				if err != nil {
					return errors.New("Cannot unmarshal config")
				}
				if envConfig.Uuid == uuid.Nil {
					return errors.New("Environment id is missing in the config")
				}
				isSrcFileValid = true
			}
			return nil
		})
		if err != nil || !isSrcFileValid {
			logger.Fatal().Err(err).Msg("Source file cannot be processed")
		}

		envConfigDir := path.Join(util.GetHomeDirectory(), util.DefaultConfigurationDirectory, util.DefaultEnvironmentsSubdirectory)

		// Check if environment with such id is already in place
		if _, err := os.Stat(path.Join(envConfigDir, envConfig.Uuid.String())); err == nil {
			logger.Fatal().Err(err).Msgf("Environment with id %s already exists", envConfig.Uuid.String())
		}

		// Unarchive specified file
		err = archiver.Unarchive(srcFile, envConfigDir)
		if err != nil {
			logger.Fatal().Err(err).Msg("Unable to import (unarchive) environment")
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
		logger.Info().Msgf("Current environment id is %s", envConfig.Uuid.String())

		// Download all Docker images for installed components
		for _, cmp := range envConfig.Installed {
			debug(cmp.String())
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
}
