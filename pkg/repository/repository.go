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
	"os"
)

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

func (v V1) ComponentsString() string {
	var b bytes.Buffer
	for _, c := range v.Components {
		for _, v := range c.Versions {
			b.WriteString(fmt.Sprintf("Component: %s:%s\n", c.Name, v.Version))
		}
	}
	return b.String()
}

func GetRepository() *V1 {
	debug("will try to get repo")
	repo, err := loadRepository()
	if err != nil {
		debug("error while loading local repo: %#v", err)
		debug("will try to download repo")
		repo, err = downloadAndPersistRepositoryV1()
		if err != nil {
			errGetRepository(err)
		}
	}
	debug("will return repo")
	return repo
}

func downloadAndPersistRepositoryV1() (*V1, error) {
	res, err := http.Get(fmt.Sprintf("%s/%s/%s/%s", util.GithubUrl, util.DefaultRepository, util.DefaultRepositoryBranch, util.DefaultV1RepositoryFileName))
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
	err = ioutil.WriteFile(util.UsedRepositoryFile, body, 0644)
	if err != nil {
		return nil, err
	}
	return repository, nil
}

func loadRepository() (*V1, error) {
	repo := &V1{}
	file, err := os.Open(util.UsedRepositoryFile)
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
