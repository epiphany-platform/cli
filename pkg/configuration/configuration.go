/*
 * Copyright © 2020 Mateusz Kyc
 */

package configuration

import (
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

type Config struct {
	Version            string    `yaml:"version"`
	Kind               string    `yaml:"kind"`
	CurrentEnvironment uuid.UUID `yaml:"current-environment"`
}

func (c *Config) CreateNewEnvironment(name string) error {
	env, err := environment.CreateEnvironment(name)
	if err != nil {
		panic("cannot create new environment")
	}
	c.CurrentEnvironment = env.Uuid
	return c.Save()
}

func (c *Config) SetUsedEnvironment(u uuid.UUID) error {
	c.CurrentEnvironment = u
	return c.Save()
}

func (c *Config) GetConfigFilePath() string {
	if usedConfigFile == "" {
		panic("usedConfigFile not initialized") //TODO err
	}
	return usedConfigFile
}

func (c *Config) Save() error {
	data, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
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
	return makeOrGetConfig()
}

func SetConfig(configFile string) (*Config, error) {
	usedConfigurationDirectory = path.Join(util.GetHomeDirectory(), util.DefaultConfigurationDirectory)
	util.EnsureDirectory(usedConfigurationDirectory)
	usedConfigFile = configFile
	return makeOrGetConfig()
}

func makeOrGetConfig() (*Config, error) {
	if _, err := os.Stat(usedConfigFile); os.IsNotExist(err) {
		config := &Config{
			Version: "v1",
			Kind:    "Config",
		}
		err = config.Save()
		if err != nil {
			panic(fmt.Sprintf("save file failed: %v\n", err)) //TODO err
		}
		return config, nil
	}
	return loadConfigFromUsedConfigFile()
}

func loadConfigFromUsedConfigFile() (*Config, error) {
	config := &Config{}
	file, err := os.Open(usedConfigFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	d := yaml.NewDecoder(file)
	if err := d.Decode(&config); err != nil {
		return nil, err
	}
	return config, nil
}
