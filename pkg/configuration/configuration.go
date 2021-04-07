package configuration

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/epiphany-platform/cli/internal/logger"
	"github.com/epiphany-platform/cli/internal/util"
	"github.com/epiphany-platform/cli/pkg/az"
	"github.com/epiphany-platform/cli/pkg/environment"

	"github.com/google/uuid"
	"gopkg.in/yaml.v2"
)

type Kind string

const (
	KindConfig Kind = "Config"
)

func init() {
	logger.Initialize()
}

type AzureConfig struct {
	Credentials az.Credentials `yaml:"credentials"`
}

type Config struct {
	Version            string      `yaml:"version"`
	Kind               Kind        `yaml:"kind"`
	CurrentEnvironment uuid.UUID   `yaml:"current-environment"`
	AzureConfig        AzureConfig `yaml:"azure-config,omitempty"`
}

//CreateNewEnvironment in Config
func (c *Config) CreateNewEnvironment(name string) (uuid.UUID, error) {
	logger.Debug().Msgf("will try to create environment %s", name)
	env, err := environment.Create(name)
	if err != nil {
		logger.Panic().Err(err).Msg("creation of new environment failed")
	}
	util.EnsureDirectory(path.Join(
		util.UsedEnvironmentDirectory,
		env.Uuid.String(),
		"/shared", //TODO to consts
	))
	c.CurrentEnvironment = env.Uuid
	logger.Debug().Msgf("will try to save updated config %+v", c)
	return env.Uuid, c.Save()
}

//SetUsedEnvironment to another value
func (c *Config) SetUsedEnvironment(u uuid.UUID) error {
	// Check if passed environment id is valid
	isEnvValid, err := environment.IsExisting(u) // TODO think if it should be here
	if err != nil {
		return err
	} else if !isEnvValid {
		return fmt.Errorf("environment %s not found", u.String())
	}

	logger.Debug().Msgf("changing used environment to %s", u.String())
	c.CurrentEnvironment = u
	logger.Debug().Msgf("will try to save updated config %+v", c)
	return c.Save()
}

//Save Config to usedConfigFile
func (c *Config) Save() error {
	logger.Debug().Msgf("will try to marshal config %+v", c)
	data, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	logger.Debug().Msgf("will try to write marshaled data to file %s", util.UsedConfigFile)
	err = ioutil.WriteFile(util.UsedConfigFile, data, 0644)
	if err != nil {
		return err
	}
	return nil
}

func (c *Config) AddAzureCredentials(credentials az.Credentials) {
	c.AzureConfig.Credentials = credentials
}

//GetConfig sets usedConfigFile and usedConfigurationDirectory to default values and returns (existing or just initialized) Config
func GetConfig() (*Config, error) {
	logger.Debug().Msg("will try to get config file")
	return SetConfigDirectory(path.Join(util.GetHomeDirectory(), util.DefaultConfigurationDirectory))
}

//SetConfigDirectory sets variable usedConfigurationDirectory and returns (existing or just initialized) Config
func SetConfigDirectory(configDir string) (*Config, error) {
	return setUsedConfigPaths(configDir, path.Join(configDir, util.DefaultConfigFileName))
}

//setUsedConfigPaths to provided values
func setUsedConfigPaths(configDir string, configFile string) (*Config, error) {
	logger.Debug().Msgf("will try to set config directory to %s", configDir)
	if util.UsedConfigurationDirectory == "" {
		util.UsedConfigurationDirectory = configDir
		util.EnsureDirectory(util.UsedConfigurationDirectory)
	} else {
		logger.Debug().Msgf("util.UsedConfigurationDirectory is already %s", util.UsedConfigurationDirectory)
	}

	logger.Debug().Msg("will try to set used config file")
	if util.UsedConfigFile == "" {
		util.UsedConfigFile = configFile
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

	logger.Debug().Msg("will try to make or get configuration")
	return makeOrGetConfig()
}

//makeOrGetConfig initializes new config file or reads existing one and returns Config
func makeOrGetConfig() (*Config, error) {
	if _, err := os.Stat(util.UsedConfigFile); os.IsNotExist(err) {
		logger.Debug().Msg("there is no config file, will try to initialize one")
		config := &Config{
			Version: "v1",
			Kind:    KindConfig,
		}
		err = config.Save()
		if err != nil {
			logger.Panic().Err(err).Msg("failed to save")
		}
		return config, nil
	}
	logger.Debug().Msgf("will try to load existing config file from %s", util.UsedConfigFile)
	config := &Config{}
	file, err := os.Open(util.UsedConfigFile)
	if err != nil {
		logger.Error().Err(err).Msgf("failed to open file %s", util.UsedConfigFile)
		return nil, err
	}
	defer file.Close()
	d := yaml.NewDecoder(file)
	logger.Trace().Msgf("will try to decode file %s to yaml", util.UsedConfigFile)
	if err := d.Decode(&config); err != nil {
		logger.Error().Err(err).Msgf("failed to decode file %s to yaml", util.UsedConfigFile)
		return nil, err
	}
	return config, nil
}
