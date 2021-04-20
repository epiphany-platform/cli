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
		logger.Error().Err(err).Msg("creation of new environment failed")
		return uuid.Nil, err
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

//GetConfig returns existing Config or fails if there is no config file or file is incorrect
func GetConfig() (*Config, error) {
	logger.Debug().Msgf("will try to load existing config file from %s", util.UsedConfigFile)
	config := &Config{}
	file, err := os.Open(util.UsedConfigFile)
	if err != nil {
		logger.Error().Err(err).Msgf("failed to open file %s", util.UsedConfigFile)
		return nil, err
	}
	defer func(f *os.File) {
		_ = f.Close()
	}(file)
	d := yaml.NewDecoder(file)
	logger.Trace().Msgf("will try to decode file %s to yaml", util.UsedConfigFile)
	if err := d.Decode(&config); err != nil {
		logger.Error().Err(err).Msgf("failed to decode file %s to yaml", util.UsedConfigFile)
		return nil, err
	}
	return config, nil
}
