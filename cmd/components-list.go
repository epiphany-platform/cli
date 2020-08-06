/*
 * Copyright © 2020 Mateusz Kyc
 */

package cmd

import (
	"fmt"
	"github.com/mkyc/epiphany-wrapper-poc/pkg/repository"
	"github.com/spf13/cobra"
)

// componentsListCmd represents the list command
var componentsListCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists all existing components in repository",
	Long:  `TODO`,
	Run: func(cmd *cobra.Command, args []string) {
		debug("component list called")
		fmt.Println(repository.GetRepository().ComponentsString())
	},
}

func init() {
	componentsCmd.AddCommand(componentsListCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// componentsListCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// componentsListCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
