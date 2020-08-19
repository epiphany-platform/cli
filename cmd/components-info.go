/*
 * Copyright © 2020 Mateusz Kyc
 */

package cmd

import (
	"errors"
	"fmt"
	"github.com/mkyc/epiphany-wrapper-poc/pkg/repository"
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
			errTooFewArguments(errors.New(fmt.Sprintf("found %d args", len(args))))
		}
		tc, err := repository.GetRepository().GetComponentByName(args[0])
		if err != nil {
			errGetComponentByName(err)
		}
		c, err := tc.JustLatestVersion()
		if err != nil {
			errGetComponentWithLatestVersion(err)
		}
		fmt.Print(c.String())
	},
}

func init() {
	componentsCmd.AddCommand(componentsInfoCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// componentsInfoCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// componentsInfoCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
