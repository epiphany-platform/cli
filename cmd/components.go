package cmd

import (
	"github.com/epiphany-platform/cli/internal/logger"
	"github.com/spf13/cobra"
)

// componentsCmd represents the components command
var componentsCmd = &cobra.Command{
	Use:   "components",
	Short: "Allows to inspect and install available components",
	Long: `This command provides way to:
 - list available components, 
 - install new component to environment
 - get information about component

Information about available components are taken from https://github.com/epiphany-platform/modules/blob/develop/v1.yaml`,
	Aliases: []string{"cmp"},
	PreRun: func(cmd *cobra.Command, args []string) {
		logger.Debug().Msg("components called")
	},
}

func init() {
	rootCmd.AddCommand(componentsCmd)
}
