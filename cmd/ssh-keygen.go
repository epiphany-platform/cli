package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// sshKeygenCmd represents the keygen command
var sshKeygenCmd = &cobra.Command{
	Use:   "keygen",
	Short: "Commands related to ssh keygen operations.",
	Long:  `Commands related to ssh keygen operations.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("keygen called")
	},
}

func init() {
	sshCmd.AddCommand(sshKeygenCmd)
}
