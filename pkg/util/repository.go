/*
 * Copyright © 2020 Mateusz Kyc
 */

package util //TODO move to another package

import (
	"errors"
	"fmt"
	"github.com/mkyc/epiphany-wrapper-poc/pkg/docker"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"time"
)

const (
	github                  = "https://raw.githubusercontent.com"
	defaultRepository       = "mkyc/epiphany-wrapper-poc-repo"
	defaultRepositoryBranch = "master"
	defaultV1FileLocation   = "v1.yaml"
)

var UsedRepositoryFilePath string

type ComponentCommand struct {
	Name                 string            `yaml:"name"`
	Description          string            `yaml:"description"`
	Command              string            `yaml:"command"`
	EnvironmentVariables map[string]string `yaml:"envs"`
	CommandArguments     []string          `yaml:"args"`
}

func (cc ComponentCommand) RunDocker(image string, workDirectory string) error {
	dockerJob := &docker.Job{
		Image:                image,
		Command:              cc.Command,
		Args:                 cc.CommandArguments,
		WorkDirectory:        workDirectory,
		EnvironmentVariables: cc.EnvironmentVariables,
	}
	return dockerJob.Run()
}

type ComponentVersion struct {
	Version       string             `yaml:"version"`
	IsLatest      bool               `yaml:"latest"`
	Image         string             `yaml:"image"`
	WorkDirectory string             `yaml:"workdir"`
	Commands      []ComponentCommand `yaml:"commands"`
}

func (cv ComponentVersion) RunDocker(command string) error {
	for _, cc := range cv.Commands {
		if cc.Name == command {
			return cc.RunDocker(cv.Image, cv.WorkDirectory)
		}
	}
	return errors.New("nothing to run for this version")
}

type Component struct {
	Name     string             `yaml:"name"`
	Type     string             `yaml:"type"`
	Versions []ComponentVersion `yaml:"versions"`
}

func (c Component) Run(command string) error {
	if c.Type == "docker" {
		for _, cv := range c.Versions {
			if cv.IsLatest { //TODO make it possible to run another version
				return cv.RunDocker(command)
			}
		}
	}
	return errors.New("nothing to run in component")
}

type RepositoryV1 struct {
	Version    string      `yaml:"version"`
	Kind       string      `yaml:"kind"`
	Components []Component `yaml:"components"`
}

func ListComponents() ([]Component, error) {
	repo, err := loadOrDownloadRepository()
	if err != nil {
		return nil, err
	}
	return repo.Components, nil
}

func GetComponent(componentName string) (*Component, error) {
	repo, err := loadOrDownloadRepository()
	if err != nil {
		return nil, err
	}
	for _, c := range repo.Components {
		if c.Name == componentName {
			return &c, nil
		}
	}
	return nil, errors.New("unknown component")
}

func GetComponentWithLatestVersion(componentName string) (*Component, error) {
	c, err := GetComponent(componentName)
	if err != nil {
		return nil, err
	}
	return buildComponentWithLatestVersion(c)
}

func init() {
	repositoryFilePath, err := initRepositoryPath()
	if err != nil {
		fmt.Println("repos error") //TODO error
		os.Exit(1)
	}
	UsedRepositoryFilePath = repositoryFilePath
}

func loadOrDownloadRepository() (*RepositoryV1, error) {
	repo, err := loadRepository()
	if err != nil {
		err = getDefaultRepository()
		if err != nil {
			return nil, err
		}
		repo, err = loadRepository()
		if err != nil {
			return nil, err
		}
	}
	return repo, nil
}

func buildComponentWithLatestVersion(component *Component) (*Component, error) {
	if len(component.Versions) < 1 {
		return nil, errors.New("no versions in component")
	}
	if len(component.Versions) == 1 {
		if component.Versions[0].IsLatest {
			return component, nil
		} else {
			return nil, errors.New("component only version is not marked latest")
		}
	}
	result := &Component{
		Name: component.Name,
		Type: component.Type,
	}
	for _, v := range component.Versions {
		if v.IsLatest {
			result.Versions = append(result.Versions, v)
		}
	}
	if len(result.Versions) != 1 {
		return nil, errors.New("incorrect number of latest versions")
	}
	return result, nil
}

func getDefaultRepository() error {
	u, err := url.Parse(github)

	if err != nil {
		return fmt.Errorf("invalid url")
	}
	u.Path = path.Join(defaultRepository, defaultRepositoryBranch, defaultV1FileLocation)

	repo, err := downloadRepositoryV1Metadata(u.String())
	if err != nil {
		return err
	}
	return writeRepository(UsedRepositoryFilePath, repo)
}

func downloadRepositoryV1Metadata(repositoryUrl string) (*RepositoryV1, error) {
	client := http.Client{
		Timeout: time.Second * 5,
	}
	req, err := http.NewRequest(http.MethodGet, repositoryUrl, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "epiphany-wrapper-poc-cli")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if res.Body != nil {
		defer res.Body.Close()
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	repository := &RepositoryV1{}
	err = yaml.Unmarshal(body, repository)
	if err != nil {
		return nil, err
	}
	return repository, nil
}

func writeRepository(repositoryPath string, repository *RepositoryV1) error {
	data, err := yaml.Marshal(repository)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(repositoryPath, data, 0644)
	if err != nil {
		return err
	}
	return nil
}

func loadRepository() (*RepositoryV1, error) {
	repo := &RepositoryV1{}
	file, err := os.Open(UsedRepositoryFilePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	d := yaml.NewDecoder(file)
	if err := d.Decode(&repo); err != nil {
		return nil, err
	}
	return repo, nil
}

func initRepositoryPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	configDirPath := path.Join(home, DefaultCfgDirectory)
	if _, err = os.Stat(configDirPath); os.IsNotExist(err) {
		_ = os.Mkdir(configDirPath, 0755)
	}
	repositoryFilePath := path.Join(configDirPath, defaultV1FileLocation)
	return repositoryFilePath, nil
}
