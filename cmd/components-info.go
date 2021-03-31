package cmd

import (
	"fmt"

	"github.com/epiphany-platform/cli/internal/logger"
	"github.com/epiphany-platform/cli/pkg/repository"
	"github.com/spf13/cobra"
)

// componentsInfoCmd represents the info command
var componentsInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Displays information about component",
	Long:  `TODO`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("found %d args", len(args))
		}
		return nil
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		logger.Debug().Msg("components info called")
	},
	Run: func(cmd *cobra.Command, args []string) {
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
