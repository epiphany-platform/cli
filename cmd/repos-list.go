package cmd

import (
	"fmt"
	"github.com/epiphany-platform/cli/internal/repository"
	"github.com/spf13/cobra"
)

// repoListCmd represents the list command
var repoListCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists installed repositories",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("list called")
		s, err := repository.List()
		if err != nil {
			panic(err)
		}
		fmt.Println(s)
	},
}

func init() {
	reposCmd.AddCommand(repoListCmd)
}
