/*
 * Copyright © 2020 Mateusz Kyc
 */

package cmd

import (
	"fmt"
	"github.com/epiphany-platform/cli/pkg/util"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"reflect"
	"strings"
	"testing"
)

func setup(t *testing.T, suffix string) (string, string, string, string) {
	parentDir := os.TempDir()
	configDirectory, err := ioutil.TempDir(parentDir, fmt.Sprintf("*-e-repository-%s", suffix))
	if err != nil {
		t.Fatal(err)
	}
	envsDirectory, err := ioutil.TempDir(configDirectory, "environments-*")
	if err != nil {
		t.Fatal(err)
	}

	configFile, err := ioutil.TempFile(configDirectory, "config-*.yaml")
	if err != nil {
		t.Fatal(err)
	}

	repoFile, err := ioutil.TempFile(configDirectory, "v1-*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	return configFile.Name(), configDirectory, envsDirectory, repoFile.Name()
}

func TestMain(m *testing.M) {
	err := os.Chdir("..")
	if err != nil {
		fmt.Printf("could not change dir: %v", err)
		os.Exit(1)
	}
	cmd := exec.Command("make", "build")
	err = cmd.Run()
	if err != nil {
		fmt.Printf("could not make binary for e: %v", err)
		os.Exit(1)
	}

	os.Exit(m.Run())
}

func TestCmd(t *testing.T) {
	util.UsedConfigFile, util.UsedConfigurationDirectory, util.UsedEnvironmentDirectory, util.UsedRepositoryFile = setup(t, "cmd")
	defer os.RemoveAll(util.UsedConfigurationDirectory)

	tests := []struct {
		name     string
		args     []string
		mockRepo bool
		mockEnv  bool
		envId    string
	}{
		{
			name:     "e help",
			args:     []string{"help"},
			mockRepo: false,
			mockEnv:  false,
		},
		{
			name:     "e components",
			args:     []string{"components"},
			mockRepo: false,
			mockEnv:  false,
		},
		{
			name:     "e components list",
			args:     []string{"--configDir", util.UsedConfigurationDirectory, "components", "list"},
			mockRepo: true,
			mockEnv:  false,
		},
		{
			name:     "e components info",
			args:     []string{"--configDir", util.UsedConfigurationDirectory, "components", "info", "c1"},
			mockRepo: true,
			mockEnv:  false,
		},
		{
			name:     "e components install",
			args:     []string{"--configDir", util.UsedConfigurationDirectory, "components", "install", "c1"},
			mockRepo: true,
			mockEnv:  true,
			envId:    "39c95814-1d01-4303-af15-ff079d609874",
		},
		{
			name:     "e environments",
			args:     []string{"environments"},
			mockRepo: false,
			mockEnv:  false,
		},
		{
			name:     "e environments info",
			args:     []string{"--configDir", util.UsedConfigurationDirectory, "environments", "info"},
			mockRepo: false,
			mockEnv:  true,
			envId:    "cd7b59f8-6610-468a-8d56-3d1ea2566428",
		},
		{
			name:     "e environments new",
			args:     []string{"environments", "new", "e1"},
			mockRepo: false,
			mockEnv:  false,
		},
		{
			name:     "e environments use",
			args:     []string{"--configDir", util.UsedConfigurationDirectory, "environments", "use", "2398d4b7-bd5e-4a2c-9efb-0bceaee6f89b"},
			mockRepo: false,
			mockEnv:  true,
			envId:    "2398d4b7-bd5e-4a2c-9efb-0bceaee6f89b",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir, err := os.Getwd()
			if err != nil {
				t.Fatal(err)
			}

			if tt.mockRepo {
				mock, err := loadFixture("cmd", tt.name+"_repo")
				if err != nil {
					t.Fatal(err)
				}
				err = ioutil.WriteFile(path.Join(util.UsedConfigurationDirectory, util.DefaultV1RepositoryFileName), []byte(mock), 0644)
				if err != nil {
					t.Fatal(err)
				}
			}

			if tt.mockEnv {
				configMock, err := loadFixture("cmd", tt.name+"_config")
				if err != nil {
					t.Fatal(err)
				}
				err = ioutil.WriteFile(path.Join(util.UsedConfigurationDirectory, util.DefaultConfigFileName), []byte(configMock), 0644)
				if err != nil {
					t.Fatal(err)
				}
				envConfigMock, err := loadFixture("cmd", tt.name+"_env_config")
				if err != nil {
					t.Fatal(err)
				}
				envPath := path.Join(util.UsedConfigurationDirectory, util.DefaultEnvironmentsSubdirectory, tt.envId)
				err = os.MkdirAll(envPath, 0775)
				if err != nil {
					t.Fatal(err)
				}
				err = ioutil.WriteFile(path.Join(envPath, util.DefaultEnvironmentConfigFileName), []byte(envConfigMock), 0644)
				if err != nil {
					t.Fatal(err)
				}
			}

			cmd := exec.Command(path.Join(dir, "output", "e"), tt.args...)
			got, err := cmd.CombinedOutput()
			if err != nil {
				t.Fatal(err)
			}

			want, err := loadFixture("cmd", tt.name+"_want")
			if err != nil {
				t.Fatal(err)
			}
			if !isTheSame(t, string(got), want) {
				return
			}

		})
	}
}

func isTheSame(t *testing.T, got string, want string) bool {
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got \n======\n\n%v\n\n=====\n, want \n======\n\n%v\n\n=====\n", got, want)
		return false
	}
	return true
}

func loadFixture(packageName, name string) (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	buf, err := ioutil.ReadFile(path.Join(dir, "fixtures", packageName, strings.ReplaceAll(name, " ", "_")))
	if err != nil {
		return "", err
	}
	return string(buf), nil
}
