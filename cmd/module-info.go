package cmd

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/rs/zerolog"
	"gopkg.in/yaml.v2"

	"github.com/epiphany-platform/cli/internal/logger"
	"github.com/epiphany-platform/cli/internal/repository"

	"github.com/spf13/cobra"
)

// moduleInfoCmd represents the search command
var moduleInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "shows ifo of named module",
	Long:  `TODO`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("there should be one positional argument")
		}
		r := regexp.MustCompile("^[0-9a-zA-Z-_]+/[0-9a-zA-Z-_]+:[0-9a-zA-Z-_.]+$") // TODO ensure github user and repo formats
		if !r.MatchString(args[0]) {
			return fmt.Errorf("module name argument incorrectly formatted")
		}
		return nil
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		logger.Debug().Msg("module info called")
	},
	Run: func(cmd *cobra.Command, args []string) {
		a := strings.Split(args[0], "/")
		repoName := a[0]
		b := strings.Split(a[1], ":")
		moduleName := b[0]
		moduleVersion := b[1]
		v, err := repository.GetModule(repoName, moduleName, moduleVersion)
		if err != nil {
			logger.Error().Err(err).Msg("info failed")
		}
		if v != nil {
			if zerolog.GlobalLevel() == zerolog.TraceLevel {
				l, _ := yaml.Marshal(v)
				logger.Trace().Msgf("will return: %s", string(l))
			}
			fmt.Print(v.String())
		} else {
			fmt.Println("module not found")
		}

	},
}

func init() {
	moduleCmd.AddCommand(moduleInfoCmd)
}
