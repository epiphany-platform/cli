package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// reposCmd represents the repos command
var reposCmd = &cobra.Command{
	Use:   "repos",
	Short: "Commands related to repos management",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("repos called")
	},
}

func init() {
	rootCmd.AddCommand(reposCmd)
}
