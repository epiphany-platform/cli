/*
 * Copyright © 2020 Mateusz Kyc
 */

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
	Run: func(cmd *cobra.Command, args []string) {
		debug("components called")
	},
}

func init() {
	rootCmd.AddCommand(componentsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// componentsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// componentsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
