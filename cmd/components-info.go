/*
 * Copyright © 2020 Mateusz Kyc
 */

package cmd

import (
	"bytes"
	"fmt"
	"github.com/mkyc/epiphany-wrapper-poc/pkg/repository"
	"github.com/spf13/cobra"
)

// componentsInfoCmd represents the info command
var componentsInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Displays information about component",
	Long:  `TODO`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("components info called")
		if len(args) != 1 {
			fmt.Println("incorrect number of args") //TODO fixme
		}
		printComponentLatestVersionInfo(args[0])
	},
}

func printComponentLatestVersionInfo(componentName string) { //TODO implement component.LatestVersionString()
	tc, err := repository.GetRepository().GetComponentByName(componentName)
	if err != nil {
		panic(fmt.Sprintf("getting component by name failed: %v\n", err)) //TODO err
	}
	c, err := tc.JustLatestVersion()
	if err != nil {
		panic(fmt.Sprintf("getting component with latest version failed: %v\n", err)) //TODO err
	}
	var b bytes.Buffer
	b.WriteString(fmt.Sprintf("Component:\n Name: %s\n Type: %s\n Version: %s\n Image: %s\n Commands:\n", c.Name, c.Type, c.Versions[0].Version, c.Versions[0].Image))
	for _, t := range c.Versions[0].Commands {
		b.WriteString(fmt.Sprintf("  Name: %s\n  Description: %s\n", t.Name, t.Description))
	}
	fmt.Println(b.String())
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
