package cmd

import (
	"fmt"
	"os"
	"path"

	// TODO: why with "github.com" - paths are different in docker image and locally?
	"github.com/epiphany-platform/cli/pkg/configuration"
	"github.com/epiphany-platform/cli/pkg/environment"
	"github.com/epiphany-platform/cli/pkg/util"
	"github.com/mholt/archiver/v3"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	envID   string
	destDir string
)

var environmentsExportCmd = &cobra.Command{
	Use:   "export",
	Short: "Exports environment as an archive",
	Long: `"export" command allows exporting any environment
	as an archive into the specified directory or
	into the current working directory by default`,
	// TODO: extend with https://pkg.go.dev/github.com/spf13/cobra#Command

	PreRun: func(cmd *cobra.Command, args []string) {
		debug("environments export called")

		err := viper.BindPFlags(cmd.Flags())
		if err != nil {
			logger.Fatal().Err(err)
		}

		envID = viper.GetString("id")
		destDir = viper.GetString("destination")
	},
	Run: func(cmd *cobra.Command, args []string) {
		isCurrentEnvUsed := false

		// Default environment and destination directory are current ones
		// Check if environment is default
		if envID == "" {
			isCurrentEnvUsed = true
			config, err := configuration.GetConfig()
			if err != nil {
				// TODO: difference between logger.Fatal() and debug
				logger.Fatal().Err(err).Msg("Unable to get environment config")
			}
			envID = config.CurrentEnvironment.String()
		}

		// Check if destination directory is default
		if destDir == "" {
			path, err := os.Getwd()
			if err != nil {
				logger.Fatal().Err(err).Msg("Unable to get working directory")
			}
			destDir = path
		}

		// Check if passed environment id is valid
		if !isCurrentEnvUsed {
			environments, err := environment.GetAll()
			if err != nil {
				logger.Fatal().Err(err).Msg("Unable to list all environments")
			}
			isEnvValid := false
			for _, e := range environments {
				if e.Uuid.String() == envID {
					isEnvValid = true
					debug("Found environment to export: %s", e.Uuid.String())
					break
				}
			}
			if !isEnvValid {
				logger.Fatal().Err(err).Msg(fmt.Sprintf("Environment %s is not found", envID))
			}
		}

		envPath := path.Join(util.UsedEnvironmentDirectory, envID)

		// Final archive name is envID + .zip extension
		// TODO: implement filtering logs out, at least runs/*.log
		err := archiver.Archive([]string{envPath}, path.Join(destDir, envID)+".zip")
		if err != nil {
			logger.Fatal().Err(err).Msg(fmt.Sprintf("Unable to archive environment directory: %s", envPath))
		}

		// TODO: any prompt action needs to be implemented here?
	},
}

func init() {
	environmentsCmd.AddCommand(environmentsExportCmd)

	environmentsExportCmd.Flags().StringP("id", "i", "", "id of the environment to export")
	environmentsExportCmd.Flags().StringP("destination", "d", "", "destination directory to store exported archive")
}
