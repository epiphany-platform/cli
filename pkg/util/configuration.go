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
	Version      string        `yaml:"version"`
	Kind         string        `yaml:"kind"`
	Environments []Environment `yaml:"environments"`
}

type Environment struct {
	Name string `yaml:"name"`
	Uuid string `yaml:"uuid"`
}

func InitDefaultConfiguration() string {
	return initDefaultConfiguration()
}

func CreateNewEnvironment(configPath string, environmentName string) error {
	config, err := loadConfig(configPath)
	if err != nil {
		return err
	}
	config.Environments = append(config.Environments, Environment{
		Name: environmentName,
		Uuid: uuid.New().String(),
	})
	err = writeConfiguration(configPath, config)
	if err != nil {
		return err
	}
	return nil
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
