/*
 * Copyright © 2020 Mateusz Kyc
 */

package cmd

import (
	"fmt"
	"github.com/mkyc/epiphany-wrapper-poc/pkg/configuration"
	"github.com/mkyc/epiphany-wrapper-poc/pkg/environment"
	"github.com/mkyc/epiphany-wrapper-poc/pkg/repository"

	"github.com/spf13/cobra"
)

// componentsInstallCmd represents the install command
var componentsInstallCmd = &cobra.Command{
	Use:   "install",
	Short: "Installs component into currently used environment",
	Long:  `TODO`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("install called")
		if len(args) != 1 {
			panic(fmt.Sprintf("incorrect number of arguments: %v", len(args)))
		}
		config, err := configuration.GetConfig()
		if err != nil {
			panic(fmt.Sprintf("get config failed: %v\n", err)) //TODO err
		}
		e, err := environment.Get(config.CurrentEnvironment)
		if err != nil {
			panic(fmt.Sprintf("environemtns get failed: %v\n", err)) //TODO err
		}

		tc, err := repository.GetRepository().GetComponentByName(args[0])
		if err != nil {
			panic(fmt.Sprintf("get component by name failed: %v\n", err)) //TODO err
		}
		c, err := tc.JustLatestVersion()
		if err != nil {
			panic(fmt.Sprintf("get component latest version failed: %v\n", err)) //TODO err
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
			panic(fmt.Sprintf("install component in environment failed: %v\n", err)) //TODO err
		}
		fmt.Printf("Installed component %s %s to environment %s\n", newComponent.Name, newComponent.Version, e.Name)
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
