/*
 * Copyright © 2020 Mateusz Kyc
 */

package configuration

import (
	"errors"
	"fmt"
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

type Kind string

const (
	KindConfig Kind = "Config"
)

type Config struct {
	Version            string    `yaml:"version"`
	Kind               Kind      `yaml:"kind"`
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

//SetUsedEnvironment to another value (NOTE: there is no additional error check)
func (c *Config) SetUsedEnvironment(u uuid.UUID) error {
	debug("changing used environment to %s", u.String())
	c.CurrentEnvironment = u
	debug("will try to save updated config %+v", c)
	return c.Save()
}

//GetConfigFilePath from usedConfigFile variable or fails if not set
func (c *Config) GetConfigFilePath() string {
	if usedConfigFile == "" {
		errIncorrectInitialization(errors.New("variable usedConfigFile not initialized"))
	}
	return usedConfigFile
}

//Save Config to usedConfigFile
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

//GetConfig sets usedConfigFile and usedConfigurationDirectory to default values and returns (existing or just initialized) Config
func GetConfig() (*Config, error) {
	debug("will try to get config file")
	if usedConfigurationDirectory == "" {
		usedConfigurationDirectory = path.Join(util.GetHomeDirectory(), util.DefaultConfigurationDirectory)
	}
	util.EnsureDirectory(usedConfigurationDirectory)
	if usedConfigFile == "" {
		usedConfigFile = path.Join(usedConfigurationDirectory, util.DefaultConfigFileName)
	}
	debug("will try to make or get configuration")
	return makeOrGetConfig()
}

//SetConfig sets variable usedConfigFile and returns (existing or just initialized) Config
func SetConfig(configFile string) (*Config, error) {
	debug("will try to set config file at %s", configFile)
	if usedConfigurationDirectory == "" {
		usedConfigurationDirectory = path.Join(util.GetHomeDirectory(), util.DefaultConfigurationDirectory)
	}
	util.EnsureDirectory(usedConfigurationDirectory)
	if usedConfigFile != "" {
		return nil, errors.New(fmt.Sprintf("usedConfigFile is %s but should be empty on set", usedConfigFile))
	}
	usedConfigFile = configFile
	debug("will try to make or get configuration")
	return makeOrGetConfig()
}

//makeOrGetConfig initializes new config file or reads existing one and returns Config
func makeOrGetConfig() (*Config, error) {
	if _, err := os.Stat(usedConfigFile); os.IsNotExist(err) {
		debug("there is no config file, will try to initialize one")
		config := &Config{
			Version: "v1",
			Kind:    KindConfig,
		}
		err = config.Save()
		if err != nil {
			errSave(err)
		}
		return config, nil
	}
	debug("will try to load existing config file from %s", usedConfigFile)
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
