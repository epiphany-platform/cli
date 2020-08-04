/*
 * Copyright © 2020 Mateusz Kyc
 */

package util

import (
	"os"
)

const (
	DefaultConfigurationDirectory      string = ".e"
	DefaultConfigFileName              string = "config.yaml"
	DefaultEnvironmentsSubdirectory    string = "environments"
	DefaultEnvironmentConfigFileName   string = "config.yaml"
	DefaultComponentRunsSubdirectory   string = "runs"
	DefaultComponentMountsSubdirectory string = "mounts"
)

var (
	UsedConfigFile             string
	UsedConfigurationDirectory string
)

func EnsureDirectory(directory string) {
	debug("will try to ensure directory %s", directory)
	err := os.MkdirAll(directory, 0755)
	if err != nil {
		errDirectoryCreation(err, directory)
	}
	debug("directory %s created", directory)
}

func GetHomeDirectory() string {
	debug("will try to get home directory")
	home, err := os.UserHomeDir()
	if err != nil {
		errFindingHome(err)
	}
	debug("got user home directory: %s", home)
	return home
}
