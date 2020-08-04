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

type Kind string

const (
	KindConfig Kind = "Config"
)

type Config struct {
	Version            string    `yaml:"version"`
	Kind               Kind      `yaml:"kind"`
	CurrentEnvironment uuid.UUID `yaml:"current-environment"`
}

//TODO return newly created environment uuid
//CreateNewEnvironment in Config
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
	if util.UsedConfigFile == "" {
		errIncorrectInitialization(errors.New("variable usedConfigFile not initialized"))
	}
	return util.UsedConfigFile
}

//Save Config to usedConfigFile
func (c *Config) Save() error {
	debug("will try to marshal config %+v", c)
	data, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	debug("will try to write marshaled data to file %s", util.UsedConfigFile)
	err = ioutil.WriteFile(util.UsedConfigFile, data, 0644)
	if err != nil {
		return err
	}
	return nil
}

//GetConfig sets usedConfigFile and usedConfigurationDirectory to default values and returns (existing or just initialized) Config
func GetConfig() (*Config, error) {
	debug("will try to get config file")
	if util.UsedConfigurationDirectory == "" {
		util.UsedConfigurationDirectory = path.Join(util.GetHomeDirectory(), util.DefaultConfigurationDirectory)
	}
	util.EnsureDirectory(util.UsedConfigurationDirectory)
	if util.UsedConfigFile == "" {
		util.UsedConfigFile = path.Join(util.UsedConfigurationDirectory, util.DefaultConfigFileName)
	}
	debug("will try to make or get configuration")
	return makeOrGetConfig()
}

//SetConfigDirectory sets variable usedConfigurationDirectory and returns (existing or just initialized) Config
func SetConfigDirectory(configDir string) (*Config, error) {
	return setUsedConfigPaths(configDir, path.Join(configDir, util.DefaultConfigFileName))
}

//setUsedConfigPaths to provided values
func setUsedConfigPaths(configDir string, configFile string) (*Config, error) {
	debug("will try to set config directory to %s", configDir)
	if util.UsedConfigurationDirectory == "" {
		util.UsedConfigurationDirectory = configDir
	}
	util.EnsureDirectory(util.UsedConfigurationDirectory)
	debug("will try to set used config file")
	if util.UsedConfigFile != "" {
		return nil, errors.New(fmt.Sprintf("usedConfigFile is %s but should be empty on set", util.UsedConfigFile))
	}
	util.UsedConfigFile = configFile
	debug("will try to make or get configuration")
	return makeOrGetConfig()
}

//makeOrGetConfig initializes new config file or reads existing one and returns Config
func makeOrGetConfig() (*Config, error) {
	if _, err := os.Stat(util.UsedConfigFile); os.IsNotExist(err) {
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
	debug("will try to load existing config file from %s", util.UsedConfigFile)
	config := &Config{}
	debug("trying to open %s file", util.UsedConfigFile)
	file, err := os.Open(util.UsedConfigFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	d := yaml.NewDecoder(file)
	debug("will try to decode file %s to yaml", util.UsedConfigFile)
	if err := d.Decode(&config); err != nil {
		return nil, err
	}
	return config, nil
}
