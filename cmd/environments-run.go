/*
 * Copyright © 2020 Mateusz Kyc
 */

package cmd

import (
	"errors"
	"fmt"
	"github.com/epiphany-platform/cli/pkg/configuration"
	"github.com/epiphany-platform/cli/pkg/environment"

	"github.com/spf13/cobra"
)

// environmentsRunCmd represents the run command
var environmentsRunCmd = &cobra.Command{ //TODO consider what are options to create integration tests here. For me it seams that it would be testing of docker
	Use:   "run",
	Short: "Runs installed component command in environment",
	Long:  `TODO`,
	PreRun: func(cmd *cobra.Command, args []string) {
		debug("environments run called")
	},
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 2 {
			config, err := configuration.GetConfig()
			if err != nil {
				errGetConfig(err)
			}
			e, err := environment.Get(config.CurrentEnvironment)
			if err != nil {
				errGetEnvironmentDetails(err)
			}
			c, err := e.GetComponentByName(args[0])
			if err != nil {
				errGetComponentByName(err)
			}
			err = c.Run(args[1])
			if err != nil {
				errRunCommand(err)
			}
			infoRunFinished(args[0], args[1])
		} else {
			errIncorrectNumberOfArguments(errors.New(fmt.Sprintf("found %d args", len(args))))
		}
	},
}

func init() {
	environmentsCmd.AddCommand(environmentsRunCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// environmentsRunCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// environmentsRunCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
