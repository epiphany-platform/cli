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
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("components info called")
		if len(args) != 1 {
			fmt.Println("incorrect number of args") //TODO fixme
		}
		printComponentLatestVersionInfo(args[0])
	},
}

func printComponentLatestVersionInfo(componentName string) {
	c, err := repository.GetComponentWithLatestVersion(componentName) //TODO implement component.LatestVersionString()
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
