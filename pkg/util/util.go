package util

import (
	"os"
)

const (
	DefaultConfigurationDirectory       string = ".e"
	DefaultConfigFileName               string = "config.yaml"
	DefaultEnvironmentsSubdirectory     string = "environments"
	DefaultEnvironmentsTempSubdirectory string = "tmp"
	DefaultEnvironmentConfigFileName    string = "config.yaml"
	DefaultComponentRunsSubdirectory    string = "runs"
	DefaultComponentMountsSubdirectory  string = "mounts"

	GithubUrl                   = "https://raw.githubusercontent.com"
	DefaultRepository           = "epiphany-platform/modules"
	DefaultRepositoryBranch     = "develop"
	DefaultV1RepositoryFileName = "v1.yaml"
)

var (
	UsedConfigFile             string
	UsedConfigurationDirectory string
	UsedEnvironmentDirectory   string
	UsedRepositoryFile         string
	UsedTempDirectory          string
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
