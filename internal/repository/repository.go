package repository

import (
	old "github.com/epiphany-platform/cli/pkg/repository"
	"github.com/epiphany-platform/cli/pkg/util"
	"gopkg.in/yaml.v2"
	"os"
	"path"
	"path/filepath"
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

func load() error {
	reposPath := path.Join(util.UsedConfigurationDirectory, repoDirectoryName)
	return filepath.Walk(reposPath, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if filepath.Ext(path) == ".yaml" || filepath.Ext(path) == ".yml" {
			v1, err2 := loadV1Repository(path)
			if err2 != nil {
				return err2
			}
			loaded.v1s = append(loaded.v1s, *v1)

		}
		return nil
	})
}

//The loadRepository method loads V1 from provided file path
func loadV1Repository(filePath string) (*old.V1, error) {
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
