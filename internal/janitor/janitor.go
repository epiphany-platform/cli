package janitor

import (
	"os"
	"path"
	"time"

	"github.com/epiphany-platform/cli/internal/logger"
	"github.com/epiphany-platform/cli/internal/repository"
	"github.com/epiphany-platform/cli/internal/util"
	"github.com/epiphany-platform/cli/pkg/configuration"

	"github.com/google/uuid"
)

func InitializeStructure(directory string) error {
	logger.Debug().Msg("InitializeStructure()")
	logger.Trace().Msg("will setUsedConfigPaths(directory)")
	setUsedConfigPaths(directory)
	logger.Trace().Msg("will ensureConfig()")
	err := ensureConfig()
	if err != nil {
		logger.Error().Err(err).Msg("ensureConfig() failed in InitializeStructure(directory string)")
		return err
	}
	logger.Trace().Msg("will ensureEnvironment()")
	err = ensureEnvironment()
	if err != nil {
		logger.Error().Err(err).Msg("ensureEnvironment() failed in InitializeStructure(directory string)")
		return err
	}
	logger.Trace().Msg("will ensureRepository()")
	err = ensureRepository()
	if err != nil {
		logger.Error().Err(err).Msg("ensureRepository() failed in InitializeStructure(directory string)")
		return err
	}
	return nil
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

//ensureConfig initializes new config if one does not exists
func ensureConfig() error {
	if _, err := os.Stat(util.UsedConfigFile); os.IsNotExist(err) {
		logger.Debug().Msg("there is no config file, will try to initialize one")
		config := &configuration.Config{
			Version: "v1",
			Kind:    configuration.KindConfig,
		}
		err = config.Save()
		if err != nil {
			logger.Error().Err(err).Msg("failed to save")
			return err
		}
	}
	return nil
}

// ensureEnvironment checks that config file do not have Nil environment and if yes, initializes it.
func ensureEnvironment() error {
	config, err := configuration.GetConfig()
	if err != nil {
		return err
	}
	if config.CurrentEnvironment == uuid.Nil {
		_, err = config.CreateNewEnvironment(time.Now().Format("060102-1504"))
		if err != nil {
			return err
		}
	}
	return nil
}

// ensureRepository tries to install default repository (but not forcibly)
func ensureRepository() error {
	return repository.Init()
}
