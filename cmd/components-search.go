package cmd

import (
	"errors"
	"fmt"
	"github.com/epiphany-platform/cli/internal/logger"
	"github.com/epiphany-platform/cli/internal/repository"

	"github.com/spf13/cobra"
)

// searchCmd represents the search command
var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Searches for component by name in all repos",
	Long:  `TODO`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("there should be one positional argument")
		}
		return nil
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		logger.Debug().Msg("search called")
	},
	Run: func(cmd *cobra.Command, args []string) {
		s, err := repository.Search(args[0])
		if err != nil {
			logger.Panic().Err(err).Msg("search failed")
		}
		fmt.Print(s)
	},
}

func init() {
	componentsCmd.AddCommand(searchCmd)
}
