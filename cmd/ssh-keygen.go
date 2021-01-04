package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// keygenCmd represents the keygen command
var keygenCmd = &cobra.Command{
	Use:   "keygen",
	Short: "TODO",
	Long:  `TODO`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("keygen called")
	},
}

func init() {
	sshCmd.AddCommand(keygenCmd)
}
