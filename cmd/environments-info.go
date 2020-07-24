/*
 * Copyright © 2020 Mateusz Kyc
 */

package cmd

import (
	"fmt"
	"github.com/mkyc/epiphany-wrapper-poc/pkg/configuration"
	"github.com/mkyc/epiphany-wrapper-poc/pkg/environment"

	"github.com/spf13/cobra"
)

// environmentsInfoCmd represents the info command
var environmentsInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("enviroments info called")
		config, err := configuration.GetConfig()
		if err != nil {
			panic(fmt.Sprintf("get config failed: %v\n", err)) //TODO err
		}
		environment, err := environment.Get(config.CurrentEnvironment)
		if err != nil {
			panic(fmt.Sprintf("environemtns details failed: %v\n", err)) //TODO err
		}
		fmt.Println(environment.String())
	},
}

func init() {
	environmentsCmd.AddCommand(environmentsInfoCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// environmentsInfoCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// environmentsInfoCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
