package cmd

import (
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

Information about available components are taken from https://github.com/mkyc/epiphany-wrapper-poc-repo/blob/master/v1.yaml`,
	PreRun: func(cmd *cobra.Command, args []string) {
		debug("components called")
	},
}

func init() {
	rootCmd.AddCommand(componentsCmd)
}
