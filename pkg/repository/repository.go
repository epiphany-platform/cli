/*
 * Copyright © 2020 Mateusz Kyc
 */

package repository //TODO move to another package

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/mkyc/epiphany-wrapper-poc/pkg/util"
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
	Name        string            `yaml:"name"`
	Description string            `yaml:"description"`
	Command     string            `yaml:"command"`
	Envs        map[string]string `yaml:"envs"`
	Args        []string          `yaml:"args"`
}

func (cc *ComponentCommand) String() string {
	return fmt.Sprintf("    Command:\n     Name %s\n     Description %s\n", cc.Name, cc.Description)
}

type ComponentVersion struct {
	Version       string             `yaml:"version"`
	IsLatest      bool               `yaml:"latest"`
	Image         string             `yaml:"image"`
	WorkDirectory string             `yaml:"workdir"`
	Mounts        []string           `yaml:"mounts"`
	Commands      []ComponentCommand `yaml:"commands"`
}

func (cv *ComponentVersion) String() string {
	var b bytes.Buffer
	b.WriteString(fmt.Sprintf("  Component Version:\n   Version: %s\n   Image: %s\n", cv.Version, cv.Image))
	for _, cc := range cv.Commands {
		b.WriteString(cc.String())
	}
	return b.String()
}

type Component struct {
	Name     string             `yaml:"name"`
	Type     string             `yaml:"type"`
	Versions []ComponentVersion `yaml:"versions"`
}

func (c *Component) String() string {
	var b bytes.Buffer
	b.WriteString(fmt.Sprintf("Component:\n Name: %s\n Type: %s\n", c.Name, c.Type))
	for _, cv := range c.Versions {
		b.WriteString(cv.String())
	}
	return b.String()
}

func (c *Component) JustLatestVersion() (*Component, error) {
	if len(c.Versions) < 1 {
		return nil, errors.New("no versions in component")
	}
	if len(c.Versions) == 1 {
		if c.Versions[0].IsLatest {
			return c, nil
		} else {
			return nil, errors.New("component only version is not marked latest")
		}
	}
	result := &Component{
		Name: c.Name,
		Type: c.Type,
	}
	for _, v := range c.Versions {
		if v.IsLatest {
			result.Versions = append(result.Versions, v)
		}
	}
	if len(result.Versions) != 1 {
		return nil, errors.New("incorrect number of latest versions")
	}
	return result, nil
}

type V1 struct {
	Version    string      `yaml:"version"`
	Kind       string      `yaml:"kind"`
	Components []Component `yaml:"components"`
}

func (v V1) GetComponentByName(name string) (*Component, error) {
	for _, c := range v.Components {
		if c.Name == name {
			return &c, nil
		}
	}
	return nil, errors.New("unknown component")
}

func GetRepository() *V1 {
	v1, err := loadOrDownloadRepository()
	if err != nil {
		errGetRepository(err)
	}
	return v1
}

func (v V1) ComponentsString() string {
	var b bytes.Buffer
	for _, c := range v.Components {
		for _, v := range c.Versions {
			b.WriteString(fmt.Sprintf("Component: %s:%s\n", c.Name, v.Version))
		}
	}
	return b.String()
}

func init() { //TODO move it to configuration
	repositoryFilePath, err := initRepositoryPath()
	if err != nil {
		errInitRepository(err)
	}
	UsedRepositoryFilePath = repositoryFilePath
}

func loadOrDownloadRepository() (*V1, error) {
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

func downloadRepositoryV1Metadata(repositoryUrl string) (*V1, error) {
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
	repository := &V1{}
	err = yaml.Unmarshal(body, repository)
	if err != nil {
		return nil, err
	}
	return repository, nil
}

func writeRepository(repositoryPath string, repository *V1) error {
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

func loadRepository() (*V1, error) {
	repo := &V1{}
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

func initRepositoryPath() (string, error) { //TODO move this responsibility to configuration?
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	configDirPath := path.Join(home, util.DefaultConfigurationDirectory)
	if _, err = os.Stat(configDirPath); os.IsNotExist(err) {
		_ = os.Mkdir(configDirPath, 0755)
	}
	repositoryFilePath := path.Join(configDirPath, defaultV1FileLocation)
	return repositoryFilePath, nil
}
