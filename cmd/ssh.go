package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// sshCmd represents the ssh command
var sshCmd = &cobra.Command{
	Use:   "ssh",
	Short: "TODO",
	Long:  `TODO`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("ssh called")
	},
}

func init() {
	rootCmd.AddCommand(sshCmd)
}
