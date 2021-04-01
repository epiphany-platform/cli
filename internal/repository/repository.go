package repository

import (
	"errors"
	"fmt"
	"github.com/epiphany-platform/cli/internal/logger"
	old "github.com/epiphany-platform/cli/pkg/repository"
	"github.com/epiphany-platform/cli/pkg/util"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	repoDirectoryName = "repos"
)

var loaded repositories

type repositories struct {
	v1s []old.V1
	// add next versions of repositories here
}

func init() {
	logger.Initialize()
}

func List() (string, error) {
	err := load()
	if err != nil {
		logger.Panic().Err(err).Msg("unable to load repos")
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

func Install(repoName string, force bool, branch string) error {
	err := load()
	if err != nil {
		logger.Panic().Err(err).Msg("unable to load repos")
	}

	logger.Debug().Msgf("will install %s", repoName)
	b := util.DefaultRepositoryBranch
	if branch != "" {
		b = branch
	}
	r, err := downloadV1Repository(fmt.Sprintf("%s/%s/%s/%s", util.GithubUrl, repoName, b, util.DefaultV1RepositoryFileName))
	if err != nil {
		return err
	}
	inferredRepoName := inferName(repoName)
	if r.Name == "" {
		r.Name = inferredRepoName
	}
	return persistV1RepositoryFile(inferredRepoName, r, force)
}

func Search(name string) (string, error) {
	err := load()
	if err != nil {
		logger.Panic().Err(err).Msg("unable to load repos")
	}

	var sb strings.Builder
	for _, v1 := range loaded.v1s {
		// TODO implement SearchComponent() and don't use GetComponentByName()
		c, _ := v1.GetComponentByName(name)
		if c != nil {
			for _, v := range c.Versions {
				sb.WriteString(fmt.Sprintf("%s/%s:%s\n", v1.Name, c.Name, v.Version))
			}
		}
	}
	return sb.String(), nil
}

func load() error {
	loaded = repositories{}
	reposPath := path.Join(util.UsedConfigurationDirectory, repoDirectoryName)
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
func decodeV1Repository(filePath string) (*old.V1, error) {
	repo := &old.V1{}
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
func downloadV1Repository(url string) (*old.V1, error) {
	logger.Trace().Msgf("will try to download repo from: %s", url)
	res, err := http.Get(url)
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
	r := &old.V1{}
	err = yaml.Unmarshal(body, r)
	if err != nil {
		logger.Error().Err(err).Msg("wasn't able to unmarshal body into correct yaml")
		return nil, err
	}
	return r, nil
}

func inferName(repo string) string {
	reg, _ := regexp.Compile("[^a-zA-Z0-9]+")
	return reg.ReplaceAllString(repo, "-")
}

func persistV1RepositoryFile(inferredRepoName string, v1 *old.V1, force bool) error {
	b, err := yaml.Marshal(v1)
	if err != nil {
		logger.Error().Err(err).Msg("wasn't able to marshal repo object into yaml")
		return err
	}
	newFilePath := path.Join(util.UsedConfigurationDirectory, repoDirectoryName, inferredRepoName+".yaml")
	if _, err = os.Stat(newFilePath); err == nil {
		logger.Debug().Msg("file " + newFilePath + " already exists")
		if !force {
			return errors.New("repo file already exists (use '--force' if you know what you do)")
		}
	}
	logger.Debug().Msg("will write yaml to file")
	return ioutil.WriteFile(newFilePath, b, 0644)
}
