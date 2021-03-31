package cmd

import (
	"errors"
	"github.com/epiphany-platform/cli/internal/repository"
	"regexp"

	"github.com/spf13/cobra"
)

// installCmd represents the install command
var installCmd = &cobra.Command{
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
	Run: func(cmd *cobra.Command, args []string) {
		logger.Debug().Msg("install called")
		err := repository.Install(args[0])
		if err != nil {
			logger.Panic().Err(err).Msg("install failed")
		}
	},
}

func init() {
	reposCmd.AddCommand(installCmd)
}
