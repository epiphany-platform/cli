/*
 * Copyright © 2020 Mateusz Kyc
 */

package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

// environmentsCmd represents the environments command
var environmentsCmd = &cobra.Command{
	Use:   "environments",
	Short: "Allows various interactions with environments",
	Long:  `TODO`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("environments called")
	},
}

func init() {
	rootCmd.AddCommand(environmentsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// environmentsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// environmentsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
