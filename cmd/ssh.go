package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// sshCmd represents the ssh command
var sshCmd = &cobra.Command{
	Use:   "ssh",
	Short: "Ssh and keys operations",
	Long:  `Commands related to ssh and keys operations.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("ssh called")
	},
}

func init() {
	rootCmd.AddCommand(sshCmd)
}
