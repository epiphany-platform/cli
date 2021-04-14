package repository

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/epiphany-platform/cli/internal/logger"
	"github.com/epiphany-platform/cli/internal/util"

	"gopkg.in/yaml.v2"
)

var (
	loaded     repositories
	httpClient = &http.Client{}
)

type repositories struct {
	v1s []V1
	// add next versions of repositories here
}

func init() {
	logger.Initialize()
}

//ComponentCommand struct contains information about specific command provided by component to be executed
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

//ComponentVersion struct contains information about version of component available to be installed
type ComponentVersion struct {
	Name          string             `yaml:"-"`
	Type          string             `yaml:"-"`
	Version       string             `yaml:"version"`
	IsLatest      bool               `yaml:"latest"`
	Image         string             `yaml:"image"`
	WorkDirectory string             `yaml:"workdir"`
	Mounts        []string           `yaml:"mounts"`
	Shared        string             `yaml:"shared"`
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

//Component struct is main element in repository identifying component and gathering all versions of it
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

//V1 struct is entrypoint repository for version 1 of used repository structure
type V1 struct {
	Version    string      `yaml:"version"`
	Kind       string      `yaml:"kind"`
	Name       string      `yaml:"name"`
	Components []Component `yaml:"components"`
}

func Init() error {
	return Install(util.DefaultRepository, false, util.DefaultRepositoryBranch)
}

func List() (string, error) {
	err := load()
	if err != nil {
		logger.Error().Err(err).Msg("unable to load repos")
		return "", err
	}

	var sb strings.Builder
	for _, v1 := range loaded.v1s {
		sb.WriteString(fmt.Sprintf("Repository: %s\n", v1.Name))
		for _, c := range v1.Components {
			for _, v := range c.Versions {
				sb.WriteString(fmt.Sprintf("\tModule: %s:%s\n", c.Name, v.Version))
			}
		}
	}
	return sb.String(), nil
}

func Install(repo string, force bool, branch string) error {
	err := load()
	if err != nil {
		logger.Error().Err(err).Msg("unable to load repos")
		return err
	}
	inferredRepoName := inferRepoName(repo)
	if !force {
		for _, v1 := range loaded.v1s {
			if v1.Name == inferredRepoName {
				logger.Debug().Msgf("looks like repo with name %s is already installed", inferredRepoName)
				return nil
			}
		}
	}

	logger.Debug().Msgf("will install %s", repo)
	b := util.DefaultRepositoryBranch
	if branch != "" {
		b = branch
	}
	u, err := url.Parse(repo)
	if err != nil {
		return err
	}
	r, err := downloadV1Repository(fmt.Sprintf("%s/%s/%s/%s", util.GithubUrl, u.Path, b, util.DefaultV1RepositoryFileName))
	if err != nil {
		return err
	}
	if r.Name == "" {
		r.Name = inferredRepoName
	}
	return persistV1RepositoryFile(inferredRepoName, r, force)
}

func Search(name string) (string, error) {
	err := load()
	if err != nil {
		logger.Error().Err(err).Msg("unable to load repos")
		return "", err
	}

	var sb strings.Builder
	for _, v1 := range loaded.v1s {
		for _, c := range v1.Components {
			if c.Name == name {
				for _, v := range c.Versions {
					sb.WriteString(fmt.Sprintf("%s/%s:%s\n", v1.Name, c.Name, v.Version))
				}
			}
		}
	}
	return sb.String(), nil
}

func GetModule(repoName, moduleName, moduleVersion string) (*ComponentVersion, error) {
	err := load()
	if err != nil {
		logger.Error().Err(err).Msg("unable to load repos")
		return nil, err
	}
	for _, v1 := range loaded.v1s {
		if repoName != "" && repoName == v1.Name {
			for _, c := range v1.Components {
				if c.Name == moduleName {
					for _, v := range c.Versions {
						if moduleVersion != "" && moduleVersion == v.Version {
							v.Name = c.Name
							v.Type = c.Type
							return &v, nil
						}
					}
				}
			}
		}
	}
	return nil, nil
}

func load() error {
	loaded = repositories{}
	reposPath := path.Join(util.UsedConfigurationDirectory, util.DefaultRepoDirectoryName)
	return filepath.Walk(reposPath, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if filepath.Ext(path) == ".yaml" || filepath.Ext(path) == ".yml" {
			v1, err2 := decodeV1Repository(path)
			if err2 != nil {
				return err2
			}
			loaded.v1s = append(loaded.v1s, *v1)
		}
		return nil
	})
}

//The decodeV1Repository method loads V1 from provided file path
func decodeV1Repository(filePath string) (*V1, error) {
	repo := &V1{}
	file, err := os.Open(filePath)
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

//The downloadV1Repository method retrieves file from provided url, unmarshalls it to V1 and returns obtained V1 struct.
func downloadV1Repository(url string) (*V1, error) {
	logger.Trace().Msgf("will try to download repo from: %s", url)
	res, err := httpClient.Get(url)
	if err != nil {
		logger.Error().Err(err).Msg("wasn't able to perform http GET on repo URL")
		return nil, err
	}
	if res.Body != nil {
		defer res.Body.Close()
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logger.Error().Err(err).Msg("wasn't able to read response body")
		return nil, err
	}
	logger.Trace().Msgf("got response body: \n%s", string(body))
	if len(body) == 0 {
		err = errors.New("empty response body")
		logger.Error().Err(err).Msgf("got nothing from: %s", url)
		return nil, err
	}
	r := &V1{}
	err = yaml.Unmarshal(body, r)
	if err != nil {
		logger.Error().Err(err).Msg("wasn't able to unmarshal body into correct yaml")
		return nil, err
	}
	return r, nil
}

func inferRepoName(repo string) string {
	u, err := url.Parse(repo)
	if err != nil {
		logger.Fatal().Err(err).Msgf("url.Parse(%s) failed", repo)
	}
	logger.Debug().Msgf("url.Parse(%s) Path: %s", repo, u.Path)
	reg, _ := regexp.Compile("[^a-zA-Z0-9]+")
	np := u.Path
	if np[0] == '/' {
		np = np[1:]
	}
	return reg.ReplaceAllString(np, "-")
}

func persistV1RepositoryFile(inferredRepoName string, v1 *V1, force bool) error {
	if v1 == nil {
		err := errors.New("nil repository")
		logger.Error().Err(err).Msg("incorrect nil parameter")
		return err
	}
	b, err := yaml.Marshal(v1)
	if err != nil {
		logger.Error().Err(err).Msg("wasn't able to marshal repo object into yaml")
		return err
	}
	filePath := path.Join(util.UsedReposDirectory, inferredRepoName+".yaml")
	if _, err = os.Stat(filePath); err == nil {
		logger.Debug().Msg("file " + filePath + " already exists")
		if !force {
			return errors.New("repo file already exists (use '--force' if you know what you do)")
		}
	}
	logger.Debug().Msgf("will write yaml to file: %s", filePath)
	return ioutil.WriteFile(filePath, b, 0644)
}
