/*
 * Copyright © 2020 Mateusz Kyc
 */

package cmd

import (
	"bytes"
	"fmt"
	"github.com/mkyc/epiphany-wrapper-poc/pkg/util"
	"os"

	"github.com/spf13/cobra"
)

// infoCmd represents the info command
var infoCmd = &cobra.Command{
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
	c, err := util.GetComponentWithLatestVersion(componentName) //TODO implement component.LatestVersionString()
	if err != nil {
		fmt.Printf("getting latest component failed: %v\n", err) //TODO err?
		os.Exit(1)
	}
	var b bytes.Buffer
	b.WriteString(fmt.Sprintf("Component:\n Name: %s\n Type: %s\n Version: %s\n Image: %s\n Commands:\n", c.Name, c.Type, c.Versions[0].Version, c.Versions[0].Image))
	for _, t := range c.Versions[0].Commands {
		b.WriteString(fmt.Sprintf("  Name: %s\n  Description: %s\n", t.Name, t.Description))
	}
	fmt.Println(b.String())
}

func init() {
	componentsCmd.AddCommand(infoCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// infoCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// infoCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
