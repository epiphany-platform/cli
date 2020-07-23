/*
 * Copyright © 2020 Mateusz Kyc
 */

package util

import (
	"fmt"
	"os"
)

const (
	DefaultConfigurationDirectory    string = ".e"
	DefaultConfigFileName            string = "config.yaml"
	DefaultEnvironmentsSubdirectory  string = "environments"
	DefaultEnvironmentConfigFileName string = "config.yaml"
)

func EnsureDirectory(directory string) {
	err := os.MkdirAll(directory, 0755)
	if err != nil {
		panic(fmt.Sprintf("directory creation failed: %v\n", err)) //TODO err
	}
}

func GetHomeDirectory() string {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(fmt.Sprintf("finding home dir failed: %v\n", err)) //TODO err
	}
	return home
}
