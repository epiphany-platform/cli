/*
 * Copyright © 2020 Mateusz Kyc
 */

package configuration

import (
	"errors"
	"github.com/google/uuid"
	"github.com/mkyc/epiphany-wrapper-poc/pkg/environment"
	"github.com/mkyc/epiphany-wrapper-poc/pkg/util"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path"
)

var (
	usedConfigFile             string
	usedConfigurationDirectory string
)

type Config struct {
	Version            string    `yaml:"version"`
	Kind               string    `yaml:"kind"`
	CurrentEnvironment uuid.UUID `yaml:"current-environment"`
}

func (c *Config) CreateNewEnvironment(name string) error {
	debug("will try to create environment %s", name)
	env, err := environment.Create(name)
	if err != nil {
		errCreateEnvironment(err)
	}
	c.CurrentEnvironment = env.Uuid
	debug("will try to save updated config %+v", c)
	return c.Save()
}

func (c *Config) SetUsedEnvironment(u uuid.UUID) error {
	debug("changing used environment to %s", u.String())
	c.CurrentEnvironment = u
	debug("will try to save updated config %+v", c)
	return c.Save()
}

func (c *Config) GetConfigFilePath() string {
	if usedConfigFile == "" {
		errIncorrectInitialization(errors.New("variable usedConfigFile not initialized"))
	}
	return usedConfigFile
}

func (c *Config) Save() error {
	debug("will try to marshal config %+v", c)
	data, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	debug("will try to write marshaled data to file %s", usedConfigFile)
	err = ioutil.WriteFile(usedConfigFile, data, 0644)
	if err != nil {
		return err
	}
	return nil
}

func GetConfig() (*Config, error) {
	usedConfigurationDirectory = path.Join(util.GetHomeDirectory(), util.DefaultConfigurationDirectory)
	util.EnsureDirectory(usedConfigurationDirectory)
	usedConfigFile = path.Join(usedConfigurationDirectory, util.DefaultConfigFileName)
	debug("will try to make or get configuration")
	return makeOrGetConfig()
}

func SetConfig(configFile string) (*Config, error) {
	debug("will try to set not default config file at %s", configFile)
	usedConfigurationDirectory = path.Join(util.GetHomeDirectory(), util.DefaultConfigurationDirectory)
	util.EnsureDirectory(usedConfigurationDirectory)
	usedConfigFile = configFile
	return makeOrGetConfig()
}

func makeOrGetConfig() (*Config, error) {
	if _, err := os.Stat(usedConfigFile); os.IsNotExist(err) {
		debug("there is no config file, will try to initialize one")
		config := &Config{
			Version: "v1",
			Kind:    "Config",
		}
		err = config.Save()
		if err != nil {
			errSave(err)
		}
		return config, nil
	}
	debug("will try to load existing config file from %s", usedConfigFile)
	return loadConfigFromUsedConfigFile()
}

func loadConfigFromUsedConfigFile() (*Config, error) {
	config := &Config{}
	debug("trying to open %s file", usedConfigFile)
	file, err := os.Open(usedConfigFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	d := yaml.NewDecoder(file)
	debug("will try to decode file %s to yaml", usedConfigFile)
	if err := d.Decode(&config); err != nil {
		return nil, err
	}
	return config, nil
}
