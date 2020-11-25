package repository

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/epiphany-platform/cli/pkg/util"
	"gopkg.in/yaml.v2"
)

//ComponentCommand struct contains information about specific command provided by component to be executed
type ComponentCommand struct {
	Name        string            `yaml:"name"`
	Description string            `yaml:"description"`
	Command     string            `yaml:"command"`
	Envs        map[string]string `yaml:"envs"`
	Args        []string          `yaml:"args"`
}

//The String method is used to pretty-print ComponentCommand struct
func (cc *ComponentCommand) String() string {
	return fmt.Sprintf("    Command:\n     Name %s\n     Description %s\n", cc.Name, cc.Description)
}

//ComponentVersion struct contains information about version of component available to be installed
type ComponentVersion struct {
	Version       string             `yaml:"version"`
	IsLatest      bool               `yaml:"latest"`
	Image         string             `yaml:"image"`
	WorkDirectory string             `yaml:"workdir"`
	Mounts        []string           `yaml:"mounts"`
	Commands      []ComponentCommand `yaml:"commands"`
}

//The String method is used to pretty-print ComponentVersion struct
func (cv *ComponentVersion) String() string {
	var b bytes.Buffer
	b.WriteString(fmt.Sprintf("  Component Version:\n   Version: %s\n   Image: %s\n", cv.Version, cv.Image))
	for _, cc := range cv.Commands {
		b.WriteString(cc.String())
	}
	return b.String()
}

//Component struct is main element in repository identifying component and gathering all versions of it
type Component struct {
	Name     string             `yaml:"name"`
	Type     string             `yaml:"type"`
	Versions []ComponentVersion `yaml:"versions"`
}

//The String method is used to pretty-print Component struct
func (c *Component) String() string {
	var b bytes.Buffer
	b.WriteString(fmt.Sprintf("Component:\n Name: %s\n Type: %s\n", c.Name, c.Type))
	for _, cv := range c.Versions {
		b.WriteString(cv.String())
	}
	return b.String()
}

//The JustLatestVersion method returns Component with just one latest ComponentVersion marked as IsLatest
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

//V1 struct is entrypoint repository for version 1 of used repository structure
type V1 struct {
	Version    string      `yaml:"version"`
	Kind       string      `yaml:"kind"`
	Components []Component `yaml:"components"`
}

//The GetComponentByName method gets first Component matching name parameter from V1 repository
func (v V1) GetComponentByName(name string) (*Component, error) {
	for _, c := range v.Components {
		if c.Name == name {
			return &c, nil
		}
	}
	return nil, errors.New("unknown component")
}

//The ComponentsString method is used to pretty-print V1 repository Component list
func (v V1) ComponentsString() string {
	var b bytes.Buffer
	for _, c := range v.Components {
		for _, v := range c.Versions {
			b.WriteString(fmt.Sprintf("Component: %s:%s\n", c.Name, v.Version))
		}
	}
	return b.String()
}

//The GetRepository method checks if there is already cached repository file and returns V1 struct. If there is no
//cache file it will try to download it from default location, persist it to cache file and return V1 as well.
func GetRepository() *V1 {
	debug("will try to get repo")
	repo, err := loadRepository(util.UsedRepositoryFile)
	if err != nil {
		debug("error while loading local repo: %#v", err)
		debug("will try to download repo")
		repo, err = downloadAndPersistRepositoryV1(fmt.Sprintf("%s/%s/%s/%s", util.GithubUrl, util.DefaultRepository, util.DefaultRepositoryBranch, util.DefaultV1RepositoryFileName))
		if err != nil {
			errGetRepository(err)
		}
	}
	debug("will return repo")
	return repo
}

//The downloadAndPersistRepositoryV1 method retrieves file from provided url, unmarshalls it to V1 and writes file to
//util.UsedRepositoryFile. Eventually it also returns obtained V1 struct.
func downloadAndPersistRepositoryV1(url string) (*V1, error) {
	res, err := http.Get(url)
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

//The loadRepository method loads V1 from provided file path
func loadRepository(repoFilePath string) (*V1, error) {
	repo := &V1{}
	file, err := os.Open(repoFilePath)
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
