package cmd

import (
	"errors"
	"fmt"

	"github.com/epiphany-platform/cli/pkg/repository"
	"github.com/spf13/cobra"
)

// componentsInfoCmd represents the info command
var componentsInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Displays information about component",
	Long:  `TODO`,
	PreRun: func(cmd *cobra.Command, args []string) {
		debug("components info called")
	},
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			logger.Fatal().Err(errors.New(fmt.Sprintf("found %d args", len(args)))).Msg("too few arguments")
		}
		tc, err := repository.GetRepository().GetComponentByName(args[0])
		if err != nil {
			logger.Fatal().Err(err).Msg("getting component by name failed")
		}
		c, err := tc.JustLatestVersion()
		if err != nil {
			logger.Fatal().Err(err).Msg("getting component with latest version failed")
		}
		fmt.Print(c.String())
	},
}

func init() {
	componentsCmd.AddCommand(componentsInfoCmd)
}
