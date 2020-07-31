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

func EnsureDirectory(directory string) {
	err := os.MkdirAll(directory, 0755)
	if err != nil {
		errDirectoryCreation(err, directory)
	}
}

func GetHomeDirectory() string {
	home, err := os.UserHomeDir()
	if err != nil {
		errFindingHome(err)
	}
	return home
}
