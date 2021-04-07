package cmd

import (
	"errors"
	"regexp"

	"github.com/epiphany-platform/cli/internal/logger"
	"github.com/epiphany-platform/cli/internal/repository"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	force  bool
	branch string
)

// reposInstallCmd represents the install command
var reposInstallCmd = &cobra.Command{
	Use:   "install",
	Short: "installs new repository",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("'install' command needs exactly one positional argument")
		}
		match, _ := regexp.MatchString(".*/.*", args[0])
		if !match {
			return errors.New("argument needs to have 'user-name/repo-name' format")
		}
		return nil
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		logger.Debug().Msg("install called")

		err := viper.BindPFlags(cmd.Flags())
		if err != nil {
			logger.Fatal().Err(err)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		err := repository.Install(args[0], force, branch)
		if err != nil {
			logger.Error().Err(err).Msg("install failed")
		}
	},
}

func init() {
	reposCmd.AddCommand(reposInstallCmd)

	reposInstallCmd.Flags().BoolVar(&force, "force", false, "force repo install even if file already exists.")
	reposInstallCmd.Flags().StringVar(&branch, "branch", "", "provide branch other than default HEAD")
}
