/*
 * Copyright © 2020 Mateusz Kyc
 */

package cmd

import (
	"fmt"
	"github.com/mkyc/epiphany-wrapper-poc/pkg/configuration"
	"github.com/mkyc/epiphany-wrapper-poc/pkg/promptui"
	"github.com/spf13/cobra"
)

// environmentsNewCmd represents the new command
var environmentsNewCmd = &cobra.Command{
	Use:   "new",
	Short: "Creates new environment",
	Long:  `TODO`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("new called") //TODO debug
		config, err := configuration.GetConfig()
		if err != nil {
			panic(fmt.Sprintf("get config failed: %v\n", err)) //TODO err
		}
		var name string
		if len(args) == 1 {
			name = args[0]
		} else {
			name, err = promptui.PromptForString("Environment name")
			if err != nil {
				panic(fmt.Sprintf("prompt for new environment failed: %v\n", err)) //TODO err
			}
		}
		fmt.Printf("new created environment is: %s\n", name) //TODO debug
		err = config.CreateNewEnvironment(name)
		if err != nil {
			panic(fmt.Sprintf("create new environemtn failed: %v\n", err)) //TODO err
		}
	},
}

func init() {
	environmentsCmd.AddCommand(environmentsNewCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// environmentsNewCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// environmentsNewCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
