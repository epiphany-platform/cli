/*
 * Copyright © 2020 Mateusz Kyc
 */

package util

import (
	"fmt"
	"github.com/google/uuid"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path"
)

const (
	DefaultCfgDirectory string = ".e"
	DefaultCfgFile      string = "config.yaml"
)

type Config struct {
	Version            string        `yaml:"version"`
	Kind               string        `yaml:"kind"`
	Environments       []Environment `yaml:"environments"`
	CurrentEnvironment uuid.UUID     `yaml:"current-environment"`
}

type Environment struct {
	Name string    `yaml:"name"`
	Uuid uuid.UUID `yaml:"uuid"`
}

func InitDefaultConfiguration() string {
	return initDefaultConfiguration()
}

func CreateNewEnvironment(configPath string, environmentName string) error {
	return createNewEnvironment(configPath, environmentName)
}

func GetConfig(configPath string) (*Config, error) {
	return loadConfig(configPath)
}

func SetUsedEnvironment(configPath string, uuid uuid.UUID) error {
	return setUsedEnvironment(configPath, uuid)
}

func setUsedEnvironment(configPath string, uuid uuid.UUID) error {
	config, err := loadConfig(configPath)
	if err != nil {
		return err
	}
	config.CurrentEnvironment = uuid
	return writeConfiguration(configPath, config)
}

func createNewEnvironment(configPath string, environmentName string) error {
	config, err := loadConfig(configPath)
	if err != nil {
		return err
	}
	newUuid := uuid.New()
	config.Environments = append(config.Environments, Environment{
		Name: environmentName,
		Uuid: newUuid,
	})
	config.CurrentEnvironment = newUuid
	return writeConfiguration(configPath, config)
}

func initDefaultConfiguration() string {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(err) //TODO err?
		os.Exit(1)
	}

	configDirPath := path.Join(home, DefaultCfgDirectory)
	if _, err = os.Stat(configDirPath); os.IsNotExist(err) {
		_ = os.Mkdir(configDirPath, 0755)
	}
	configFilePath := path.Join(configDirPath, DefaultCfgFile)
	if _, err = os.Stat(configFilePath); os.IsNotExist(err) {
		err = createInitConfigFile(configFilePath)
		if err != nil {
			fmt.Println(err) //TODO warn?
			os.Exit(1)
		}
	}
	return configFilePath
}

func createInitConfigFile(configPath string) error {
	config := &Config{
		Version: "v1",
		Kind:    "Config",
	}
	return writeConfiguration(configPath, config)
}

func loadConfig(configPath string) (*Config, error) {
	config := &Config{}
	file, err := os.Open(configPath)
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

func writeConfiguration(configPath string, config *Config) error {
	data, err := yaml.Marshal(config)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(configPath, data, 0644)
	if err != nil {
		return err
	}
	return nil
}
