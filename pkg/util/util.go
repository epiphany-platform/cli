package util

import (
	"github.com/epiphany-platform/cli/internal/logger"
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
	DefaultRepoDirectoryName            string = "repos"

	GithubUrl                   = "https://raw.githubusercontent.com"
	DefaultRepository           = "epiphany-platform/modules"
	DefaultRepositoryBranch     = "HEAD"
	DefaultV1RepositoryFileName = "v1.yaml"
)

var (
	UsedConfigFile             string
	UsedConfigurationDirectory string
	UsedEnvironmentDirectory   string
	UsedRepositoryFile         string
	UsedTempDirectory          string
	UsedReposDirectory         string
)

func init() {
	logger.Initialize()
}

func EnsureDirectory(directory string) {
	logger.Debug().Msgf("will try to ensure directory %s", directory)
	err := os.MkdirAll(directory, 0755)
	if err != nil {
		logger.Panic().Err(err).Msgf("directory %s creation failed", directory)
	}
	logger.Debug().Msgf("directory %s created", directory)
}

func GetHomeDirectory() string {
	logger.Debug().Msg("will try to get home directory")
	home, err := os.UserHomeDir()
	if err != nil {
		logger.Panic().Err(err).Msg("cannot determine home directory")
	}
	logger.Debug().Msgf("got user home directory: %s", home)
	return home
}
