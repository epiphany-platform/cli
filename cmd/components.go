/*
 * Copyright © 2020 Mateusz Kyc
 */

package cmd

import (
	"fmt"
	"github.com/mkyc/epiphany-wrapper-poc/pkg/repository"
	"github.com/spf13/cobra"
)

// componentsCmd represents the components command
var componentsCmd = &cobra.Command{
	Use:   "components",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("components called")
		if len(args) == 1 {
			printComponentLatestVersionInfo(args[0])
		}
		if len(args) == 2 {
			c, err := repository.GetRepository().GetComponentByName(args[0])
			if err != nil {
				panic(fmt.Sprintf("get component failed: %v\n", err)) //TODO err
			}
			err = c.Run(args[1])
			if err != nil {
				panic(fmt.Sprintf("run command failed: %v\n", err)) //TODO err
			}
			fmt.Printf("running command completed!")
		}
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
