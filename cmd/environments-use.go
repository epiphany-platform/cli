/*
 * Copyright © 2020 Mateusz Kyc
 */

package cmd

import (
	"github.com/epiphany-platform/cli/pkg/configuration"
	"github.com/epiphany-platform/cli/pkg/promptui"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

// environmentsUseCmd represents the use command
var environmentsUseCmd = &cobra.Command{
	Use:   "use",
	Short: "Allows to select environment to be used",
	Long:  `TODO`,
	PreRun: func(cmd *cobra.Command, args []string) {
		debug("environments use called")
	},
	Run: func(cmd *cobra.Command, args []string) {
		config, err := configuration.GetConfig()
		if err != nil {
			errGetConfig(err)
		}
		if config.CurrentEnvironment == uuid.Nil {
			errNilEnvironment()
		}

		var u uuid.UUID
		if len(args) == 1 {
			u = uuid.MustParse(args[0])
		} else {
			u, err = promptui.PromptForEnvironmentSelect("Environments")
			if err != nil {
				errPrompt(err)
			}
		}

		infoChosenEnvironment(u.String())
		err = config.SetUsedEnvironment(u)
		if err != nil {
			errSetEnvironment(err)
		}
	},
}

func init() {
	environmentsCmd.AddCommand(environmentsUseCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// environmentsUseCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// environmentsUseCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
