package util

import (
	"fmt"
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
	Version string `yaml:"version"`
	Kind    string `yaml:"kind"`
}

func InitDefaultConfiguration() string {
	return initDefaultConfiguration()
}

func initDefaultConfiguration() string {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(err) //TODO log
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
			fmt.Println(err) //TODO log
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
