/*
 * Copyright © 2020 Mateusz Kyc
 */

package cmd

import (
	"fmt"
	"github.com/mkyc/epiphany-wrapper-poc/pkg/util"
	"github.com/spf13/viper"
	"os"

	"github.com/spf13/cobra"
)

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("new called") //TODO debug
		name, err := util.PromptForString("Environment name")
		if err != nil {
			fmt.Printf("environment new failed: %v\n", err) //TODO warn?
			os.Exit(1)
		}
		fmt.Printf("name is: %s\n", name) //TODO debug
		err = util.CreateNewEnvironment(viper.ConfigFileUsed(), name)
		if err != nil {
			fmt.Printf("create new environment failed: %v\n", err) //TODO warn?
			os.Exit(1)
		}
	},
}

func init() {
	environmentsCmd.AddCommand(newCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// newCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// newCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
