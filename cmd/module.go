package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// moduleCmd represents the module command
var moduleCmd = &cobra.Command{
	Use:   "module",
	Short: "Manages modules",
	Long:  `TODO`,
	PreRun: func(cmd *cobra.Command, args []string) {
		fmt.Println("module called")
	},
}

func init() {
	rootCmd.AddCommand(moduleCmd)
}
