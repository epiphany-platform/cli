package cmd

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/epiphany-platform/cli/pkg/environment"

	"github.com/epiphany-platform/cli/internal/logger"
	"github.com/epiphany-platform/cli/internal/repository"

	"github.com/spf13/cobra"
)

// moduleInstallCmd represents the search command
var moduleInstallCmd = &cobra.Command{
	Use:   "install",
	Short: "installs module into currently used environment",
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
		logger.Debug().Msg("module install called")
	},
	Run: func(cmd *cobra.Command, args []string) {
		a := strings.Split(args[0], "/")
		repoName := a[0]
		b := strings.Split(a[1], ":")
		moduleName := b[0]
		moduleVersion := b[1]
		v, err := repository.GetModule(repoName, moduleName, moduleVersion)
		if err != nil {
			logger.Fatal().Err(err).Msg("get module failed")
		}
		if v == nil {
			logger.Fatal().Msgf("module not found: %s", args[0])
		}
		newComponent := environment.InstalledComponentVersion{
			EnvironmentRef: currentEnvironment.Uuid,
			Name:           v.Name,
			Type:           v.Type,
			Version:        v.Version,
			Image:          v.Image,
			WorkDirectory:  v.WorkDirectory,
			Mounts:         v.Mounts,
			Shared:         v.Shared,
		}
		for _, rc := range v.Commands {
			nic := environment.InstalledComponentCommand{
				Name:        rc.Name,
				Description: rc.Description,
				Command:     rc.Command,
				Envs:        rc.Envs,
				Args:        rc.Args,
			}
			newComponent.Commands = append(newComponent.Commands, nic)
		}
		err = currentEnvironment.Install(newComponent)
		if err != nil {
			logger.Fatal().Err(err).Msg("install module in environment failed")
		}
		fmt.Printf("Installed module %s:%s to environment %s\n", newComponent.Name, newComponent.Version, currentEnvironment.Name)
	},
}

func init() {
	moduleCmd.AddCommand(moduleInstallCmd)
}
