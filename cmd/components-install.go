package cmd

import (
	"errors"
	"fmt"

	"github.com/epiphany-platform/cli/pkg/configuration"
	"github.com/epiphany-platform/cli/pkg/environment"
	"github.com/epiphany-platform/cli/pkg/repository"

	"github.com/spf13/cobra"
)

// componentsInstallCmd represents the install command
var componentsInstallCmd = &cobra.Command{
	Use:   "install",
	Short: "Installs component into currently used environment",
	Long:  `TODO`,
	PreRun: func(cmd *cobra.Command, args []string) {
		debug("components install called")
	},
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			logger.
				Fatal().
				Err(errors.New(fmt.Sprintf("found %d args", len(args)))).
				Msg("incorrect number of arguments")
		}
		config, err := configuration.GetConfig()
		if err != nil {
			logger.Fatal().Err(err).Msg("get config failed")
		}
		e, err := environment.Get(config.CurrentEnvironment)
		if err != nil {
			logger.Fatal().Err(err).Msg("environments get failed")
		}

		tc, err := repository.GetRepository().GetComponentByName(args[0])
		if err != nil {
			logger.Fatal().Err(err).Msg("getting component by name failed")
		}
		c, err := tc.JustLatestVersion()
		if err != nil {
			logger.Fatal().Err(err).Msg("getting component with latest version failed")
		}

		newComponent := environment.InstalledComponentVersion{
			EnvironmentRef: e.Uuid,
			Name:           c.Name,
			Type:           c.Type,
			Version:        c.Versions[0].Version,
			Image:          c.Versions[0].Image,
			WorkDirectory:  c.Versions[0].WorkDirectory,
			Mounts:         c.Versions[0].Mounts,
		}
		for _, rc := range c.Versions[0].Commands {
			nic := environment.InstalledComponentCommand{
				Name:        rc.Name,
				Description: rc.Description,
				Command:     rc.Command,
				Envs:        rc.Envs,
				Args:        rc.Args,
			}
			newComponent.Commands = append(newComponent.Commands, nic)
		}
		err = e.Install(newComponent)
		if err != nil {
			logger.Fatal().Err(err).Msg("install component in environment failed")
		}
		fmt.Printf("Installed component %s:%s to environment %s\n", newComponent.Name, newComponent.Version, e.Name)
	},
}

func init() {
	componentsCmd.AddCommand(componentsInstallCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// componentsInstallCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// componentsInstallCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
