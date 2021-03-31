package repository

import (
	"errors"
	"fmt"
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

func List() (string, error) {
	if len(loaded.v1s) < 1 {
		err := load()
		if err != nil {
			return "", err
		}
	}
	var sb strings.Builder
	for _, v1 := range loaded.v1s {
		// TODO add name here
		sb.WriteString("add name here\n")
		sb.WriteString(v1.ComponentsString())
	}
	return sb.String(), nil
}

func Install(repoName string) error {
	fmt.Printf("will install %s\n", repoName)
	r, err := downloadV1Repository(fmt.Sprintf("%s/%s/%s/%s", util.GithubUrl, repoName, util.DefaultRepositoryBranch, util.DefaultV1RepositoryFileName))
	if err != nil {
		return err
	}
	return persistV1RepositoryFile(repoName, r)
}

func load() error {
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
	r := &old.V1{}
	err = yaml.Unmarshal(body, r)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func inferName(repo string) string {
	reg, _ := regexp.Compile("[^a-zA-Z0-9]+")
	return reg.ReplaceAllString(repo, "-")
}

func persistV1RepositoryFile(repo string, v1 *old.V1) error {
	b, err := yaml.Marshal(v1)
	if err != nil {
		return err
	}
	newFilePath := path.Join(util.UsedConfigurationDirectory, repoDirectoryName, inferName(repo)+".yaml")
	if _, err = os.Stat(newFilePath); err == nil {
		return errors.New("repo file already exists")
	}

	return ioutil.WriteFile(newFilePath, b, 0644)
}
