package janitor

import (
	"path"

	"github.com/epiphany-platform/cli/internal/logger"
	"github.com/epiphany-platform/cli/internal/util"
)

func InitializeFilesStructure(directory string) {
	logger.Debug().Msg("InitializeFilesStructure()")
	setUsedConfigPaths(directory)
}

//setUsedConfigPaths to provided values
func setUsedConfigPaths(configDir string) {
	logger.Debug().Msgf("will try to set config directory to %s", configDir)
	if util.UsedConfigurationDirectory == "" {
		util.UsedConfigurationDirectory = configDir
		util.EnsureDirectory(util.UsedConfigurationDirectory)
	} else {
		logger.Debug().Msgf("util.UsedConfigurationDirectory is already %s", util.UsedConfigurationDirectory)
	}

	logger.Debug().Msg("will try to set used config file")
	if util.UsedConfigFile == "" {
		util.UsedConfigFile = path.Join(configDir, util.DefaultConfigFileName)
	} else {
		logger.Debug().Msgf("util.UsedConfigFile is already %s", util.UsedConfigFile)
	}

	logger.Debug().Msg("will try to set used environments directory")
	if util.UsedEnvironmentDirectory == "" {
		util.UsedEnvironmentDirectory = path.Join(configDir, util.DefaultEnvironmentsSubdirectory)
		util.EnsureDirectory(util.UsedEnvironmentDirectory)
	} else {
		logger.Debug().Msgf("util.UsedEnvironmentDirectory is already %s", util.UsedEnvironmentDirectory)
	}

	logger.Debug().Msg("will try to set used temporary directory")
	if util.UsedTempDirectory == "" {
		util.UsedTempDirectory = path.Join(configDir, util.DefaultEnvironmentsTempSubdirectory)
		util.EnsureDirectory(util.UsedTempDirectory)
	} else {
		logger.Debug().Msgf("util.UsedTempDirectory is already %s", util.UsedTempDirectory)
	}

	logger.Debug().Msg("will try to set repo config file path")
	if util.UsedRepositoryFile == "" {
		util.UsedRepositoryFile = path.Join(configDir, util.DefaultV1RepositoryFileName)
	} else {
		logger.Debug().Msgf("util.UsedRepositoryFile is already %s", util.UsedRepositoryFile)
	}

	logger.Debug().Msg("will try to set repos directory")
	if util.UsedReposDirectory == "" {
		util.UsedReposDirectory = path.Join(configDir, util.DefaultRepoDirectoryName)
		util.EnsureDirectory(util.UsedReposDirectory)
	} else {
		logger.Debug().Msgf("util.UsedReposDirectory is already %s", util.UsedTempDirectory)
	}
}
