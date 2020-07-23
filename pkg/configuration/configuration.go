/*
 * Copyright © 2020 Mateusz Kyc
 */

package configuration

import (
	"fmt"
	"github.com/google/uuid"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path"
)

const (
	DefaultCfgDirectory             string = ".e"
	defaultCfgFile                  string = "config.yaml"
	defaultEnvironmentsSubdirectory string = "environments"
)

var (
	usedConfigFile             string
	usedConfigurationDirectory string
)

type Config struct {
	Version            string        `yaml:"version"`
	Kind               string        `yaml:"kind"`
	Environments       []Environment `yaml:"environments"`
	CurrentEnvironment uuid.UUID     `yaml:"current-environment"`
}

func (c *Config) CreateNewEnvironment(name string) error {
	newUuid := uuid.New()
	c.Environments = append(c.Environments, Environment{
		Name: name,
		Uuid: newUuid,
	})
	c.CurrentEnvironment = newUuid
	newEnvironmentDirectory := path.Join(usedConfigurationDirectory, defaultEnvironmentsSubdirectory, newUuid.String())
	ensureDirectory(newEnvironmentDirectory)
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

type Environment struct {
	Name string    `yaml:"name"`
	Uuid uuid.UUID `yaml:"uuid"`
}

func NewConfig() (*Config, error) {
	usedConfigurationDirectory = path.Join(getHomeDirectory(), DefaultCfgDirectory)
	ensureDirectory(usedConfigurationDirectory)
	usedConfigFile = path.Join(usedConfigurationDirectory, defaultCfgFile)
	return makeOrGetConfig()
}

func SetConfig(configFile string) (*Config, error) {
	usedConfigurationDirectory = path.Join(getHomeDirectory(), DefaultCfgDirectory)
	ensureDirectory(usedConfigurationDirectory)
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

func getHomeDirectory() string {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(fmt.Sprintf("finding home dir failed: %v\n", err)) //TODO err
	}
	return home
}

func ensureDirectory(directory string) {
	err := os.MkdirAll(directory, 0755)
	if err != nil {
		panic(fmt.Sprintf("directory creation failed: %v\n", err)) //TODO err
	}
}
