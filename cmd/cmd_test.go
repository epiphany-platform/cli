package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"reflect"
	"strings"
	"testing"

	"github.com/epiphany-platform/cli/internal/util"

	"github.com/stretchr/testify/assert"
)

func setup(t *testing.T, suffix string) (string, string, string, string, string) {
	parentDir := os.TempDir()
	configDirectory, err := ioutil.TempDir(parentDir, fmt.Sprintf("*-e-repository-%s", suffix))
	if err != nil {
		t.Fatal(err)
	}
	envsDirectory := path.Join(configDirectory, util.DefaultEnvironmentsSubdirectory)
	err = os.Mkdir(envsDirectory, 0755)
	if err != nil {
		t.Fatal(err)
	}
	tempDirectory := path.Join(configDirectory, util.DefaultEnvironmentsTempSubdirectory)
	err = os.Mkdir(tempDirectory, 0755)
	if err != nil {
		t.Fatal(err)
	}

	configFile := path.Join(configDirectory, util.DefaultConfigFileName)

	repoFile := path.Join(configDirectory, util.DefaultV1RepositoryFileName)

	return configFile, configDirectory, envsDirectory, repoFile, tempDirectory
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
	util.UsedConfigFile, util.UsedConfigurationDirectory, util.UsedEnvironmentDirectory, util.UsedRepositoryFile, util.UsedTempDirectory = setup(t, "cmd")
	defer os.RemoveAll(util.UsedConfigurationDirectory)

	tests := []struct {
		name        string
		args        []string
		mockRepo    bool
		mockEnv     bool
		envId       string
		shouldFail  bool
		checkOutput bool
	}{
		{
			name:        "e help",
			args:        []string{"help"},
			mockRepo:    false,
			mockEnv:     false,
			shouldFail:  false,
			checkOutput: true,
		},
		{
			name:        "e components",
			args:        []string{"components"},
			mockRepo:    false,
			mockEnv:     false,
			shouldFail:  false,
			checkOutput: true,
		},
		{
			name:        "e components list",
			args:        []string{"--configDir", util.UsedConfigurationDirectory, "components", "list"},
			mockRepo:    true,
			mockEnv:     false,
			shouldFail:  false,
			checkOutput: true,
		},
		{
			name:        "e components install",
			args:        []string{"--configDir", util.UsedConfigurationDirectory, "components", "install", "c1"},
			mockRepo:    true,
			mockEnv:     true,
			envId:       "39c95814-1d01-4303-af15-ff079d609874",
			shouldFail:  false,
			checkOutput: true,
		},
		{
			name:        "e environments",
			args:        []string{"environments"},
			mockRepo:    false,
			mockEnv:     false,
			shouldFail:  false,
			checkOutput: true,
		},
		{
			name:        "e environments info",
			args:        []string{"--configDir", util.UsedConfigurationDirectory, "environments", "info"},
			mockRepo:    false,
			mockEnv:     true,
			envId:       "cd7b59f8-6610-468a-8d56-3d1ea2566428",
			shouldFail:  false,
			checkOutput: true,
		},
		{
			name:        "e environments new",
			args:        []string{"--configDir", util.UsedConfigurationDirectory, "environments", "new", "e1"},
			mockRepo:    false,
			mockEnv:     false,
			shouldFail:  false,
			checkOutput: false,
		},
		{
			name:        "e environments use",
			args:        []string{"--configDir", util.UsedConfigurationDirectory, "environments", "use", "2398d4b7-bd5e-4a2c-9efb-0bceaee6f89b"},
			mockRepo:    false,
			mockEnv:     true,
			envId:       "2398d4b7-bd5e-4a2c-9efb-0bceaee6f89b",
			shouldFail:  false,
			checkOutput: false,
		},
		{
			name:        "e environments export",
			args:        []string{"--configDir", util.UsedConfigurationDirectory, "environments", "export", "--id", "2e309e00-aea0-41bc-8344-9813970ec2a6", "--destination", util.UsedConfigurationDirectory},
			mockRepo:    false,
			mockEnv:     true,
			envId:       "2e309e00-aea0-41bc-8344-9813970ec2a6",
			shouldFail:  false,
			checkOutput: false,
		},
		{
			name:        "e environments export wrong env id",
			args:        []string{"--configDir", util.UsedConfigurationDirectory, "environments", "export", "--id", "fcfd81e4-27a8-4ee6-8bb3-f71b8218ba6d"},
			mockRepo:    false,
			mockEnv:     false,
			shouldFail:  true,
			checkOutput: false,
		},
		{
			name:        "e environments export wrong destination",
			args:        []string{"--configDir", util.UsedConfigurationDirectory, "environments", "export", "--id", "ba03a2ba-8fa0-4c15-ac07-894af3dbb365", "--destination", "/fake/path"},
			mockRepo:    false,
			mockEnv:     true,
			envId:       "ba03a2ba-8fa0-4c15-ac07-894af3dbb365",
			shouldFail:  true,
			checkOutput: false,
		},
		{
			name:        "e environments import",
			args:        []string{"--configDir", util.UsedConfigurationDirectory, "environments", "import", "--from", path.Join(util.UsedConfigurationDirectory, "ba03a2ba-8fa0-4c15-ac07-894af3dbb365.zip")},
			mockRepo:    false,
			mockEnv:     false,
			envId:       "ba03a2ba-8fa0-4c15-ac07-894af3dbb365",
			shouldFail:  true,
			checkOutput: false,
		},
		{
			name:        "e environments import existing",
			args:        []string{"--configDir", util.UsedConfigurationDirectory, "environments", "import", "--from", path.Join(util.UsedConfigurationDirectory, "ba03a2ba-8fa0-4c15-ac07-894af3dbb365.zip")},
			mockRepo:    false,
			mockEnv:     false,
			shouldFail:  true,
			checkOutput: false,
		},
		{
			name:        "e environments import not existing",
			args:        []string{"--configDir", util.UsedConfigurationDirectory, "environments", "import", "--from", path.Join(util.UsedConfigurationDirectory, "2e309e00-aea0-41bc-8344-9813970ec2a6.zip")},
			mockRepo:    false,
			mockEnv:     false,
			envId:       "2e309e00-aea0-41bc-8344-9813970ec2a6",
			shouldFail:  false,
			checkOutput: false,
		},
		{
			name:        "e environments import wrong source file",
			args:        []string{"--configDir", util.UsedConfigurationDirectory, "environments", "import", "--from", "/fake/path"},
			mockRepo:    false,
			mockEnv:     false,
			shouldFail:  true,
			checkOutput: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := assert.New(t)
			dir, err := os.Getwd()
			if err != nil {
				t.Fatal(err)
			}

			if tt.mockRepo {
				mock, err := loadFixture("cmd", tt.name+"_repo")
				if err != nil {
					t.Fatal(err)
				}
				err = ioutil.WriteFile(util.UsedRepositoryFile, []byte(mock), 0644)
				if err != nil {
					t.Fatal(err)
				}
			}

			if tt.mockEnv {
				configMock, err := loadFixture("cmd", tt.name+"_config")
				if err != nil {
					t.Fatal(err)
				}
				err = ioutil.WriteFile(util.UsedConfigFile, []byte(configMock), 0644)
				if err != nil {
					t.Fatal(err)
				}
				envConfigMock := fmt.Sprintf(`name: e1
uuid: %s
`, tt.envId)
				envPath := path.Join(util.UsedEnvironmentDirectory, tt.envId)
				err = os.MkdirAll(envPath, 0755)
				if err != nil {
					t.Fatal(err)
				}
				err = ioutil.WriteFile(path.Join(envPath, util.DefaultEnvironmentConfigFileName), []byte(envConfigMock), 0644)
				if err != nil {
					t.Fatal(err)
				}
			}

			if tt.name == "e environments import not existing" {
				err = os.RemoveAll(path.Join(util.UsedEnvironmentDirectory, tt.envId))
				if err != nil {
					t.Fatal(err)
				}
			}

			cmd := exec.Command(path.Join(dir, "output", "e"), tt.args...)
			got, err := cmd.CombinedOutput()
			if tt.shouldFail {
				a.Error(err)
			} else {
				a.NoError(err)
			}
			if tt.checkOutput {
				want, err := loadFixture("cmd", tt.name+"_want")
				if err != nil {
					t.Fatal(err)
				}
				a.Truef(reflect.DeepEqual(string(got), want), "got \n======\n\n%v\n\n=====\n, want \n======\n\n%v\n\n=====\n", string(got), want)
			}
		})
	}
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
