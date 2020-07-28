/*
 * Copyright © 2020 Mateusz Kyc
 */

package environment

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/mkyc/epiphany-wrapper-poc/pkg/docker"
	"github.com/mkyc/epiphany-wrapper-poc/pkg/util"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path"
	"time"
)

var (
	usedEnvironmentDirectory string
)

type InstalledComponentCommand struct {
	Name        string            `yaml:"name"`
	Description string            `yaml:"description"`
	Command     string            `yaml:"command"`
	Envs        map[string]string `yaml:"envs"`
	Args        []string          `yaml:"args"`
}

func (cc *InstalledComponentCommand) RunDocker(image string, workDirectory string) error {
	dockerJob := &docker.Job{
		Image:                image,
		Command:              cc.Command,
		Args:                 cc.Args,
		WorkDirectory:        workDirectory,
		EnvironmentVariables: cc.Envs,
	}
	return dockerJob.Run()
}

func (cc *InstalledComponentCommand) String() string {
	return fmt.Sprintf("    Command:\n     Name %s\n     Description %s\n", cc.Name, cc.Description)
}

type InstalledComponentVersion struct {
	EnvironmentRef uuid.UUID                   `yaml:"environment_ref"`
	Name           string                      `yaml:"name"`
	Type           string                      `yaml:"type"`
	Version        string                      `yaml:"version"`
	Image          string                      `yaml:"image"`
	WorkDirectory  string                      `yaml:"workdir"`
	Commands       []InstalledComponentCommand `yaml:"commands"`
}

func (cv *InstalledComponentVersion) Run(command string) error {
	if cv.Type == "docker" {
		for _, cc := range cv.Commands {
			if cc.Name == command {
				return cc.RunDocker(cv.Image, cv.WorkDirectory)
			}
		}
	}
	return errors.New("nothing to run for this version")
}

func (cv *InstalledComponentVersion) String() string {
	var b bytes.Buffer
	b.WriteString(fmt.Sprintf("  Installed Component:\n   Name: %s\n   Type: %s\n   Version: %s\n   Image: %s\n", cv.Name, cv.Type, cv.Version, cv.Image))
	for _, cc := range cv.Commands {
		b.WriteString(cc.String())
	}
	return b.String()
}

func (cv *InstalledComponentVersion) Download() error {
	if cv.Type == "docker" {
		dockerImage := &docker.Image{Name: cv.Image}
		logs, err := dockerImage.Pull()
		cv.PersistLogs(logs)
		if err != nil {
			return err
		}
		return nil
	}
	return nil
}

func (cv *InstalledComponentVersion) PersistLogs(logs string) {
	logsPath := path.Join(
		usedEnvironmentDirectory,
		cv.EnvironmentRef.String(),
		cv.Name,
		cv.Version,
		util.DefaultComponentRunsSubdirectory,
		fmt.Sprintf("%s.log", time.Now().Format("20060102-150405.000MST")),
	)
	err := ioutil.WriteFile(logsPath, []byte(logs), 0644)
	if err != nil {
		panic("failed when writing logs")
	}
}

type Environment struct {
	Name      string                      `yaml:"name"`
	Uuid      uuid.UUID                   `yaml:"uuid"`
	Installed []InstalledComponentVersion `yaml:"installed"`
}

func (e *Environment) Save() error {
	data, err := yaml.Marshal(e)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(path.Join(usedEnvironmentDirectory, e.Uuid.String(), util.DefaultEnvironmentConfigFileName), data, 0644)
	if err != nil {
		return err
	}
	return nil
}

func (e *Environment) String() string {
	var b bytes.Buffer
	b.WriteString("Environment info:\n")
	b.WriteString(fmt.Sprintf(" Name: %s\n UUID: %s\n", e.Name, e.Uuid.String()))
	for _, ic := range e.Installed {
		b.WriteString(ic.String())
	}
	return b.String()
}

func (e *Environment) Install(newComponent InstalledComponentVersion) error {
	for _, ic := range e.Installed {
		if ic.Name == newComponent.Name && ic.Version == newComponent.Version {
			return errors.New("there is this version of component already installed in environment")
		}
	}
	e.Installed = append(e.Installed, newComponent)
	newComponentRunsDirectory := path.Join(usedEnvironmentDirectory, e.Uuid.String(), newComponent.Name, newComponent.Version, util.DefaultComponentRunsSubdirectory)
	newComponentMountsDirectory := path.Join(usedEnvironmentDirectory, e.Uuid.String(), newComponent.Name, newComponent.Version, util.DefaultComponentMountsSubdirectory)
	util.EnsureDirectory(newComponentRunsDirectory)
	util.EnsureDirectory(newComponentMountsDirectory)
	err := newComponent.Download()
	if err != nil {
		return err
	}
	return e.Save()
}

func (e *Environment) GetComponentByName(name string) (*InstalledComponentVersion, error) {
	for _, ic := range e.Installed {
		if ic.Name == name {
			return &ic, nil
		}
	}
	return nil, errors.New("no such component installed")
}

func init() {
	usedEnvironmentDirectory = path.Join(util.GetHomeDirectory(), util.DefaultConfigurationDirectory, util.DefaultEnvironmentsSubdirectory)
	util.EnsureDirectory(usedEnvironmentDirectory)
}

func Create(name string) (*Environment, error) {
	environment := &Environment{
		Name: name,
		Uuid: uuid.New(),
	}

	newEnvironmentDirectory := path.Join(usedEnvironmentDirectory, environment.Uuid.String())
	util.EnsureDirectory(newEnvironmentDirectory)
	err := environment.Save()
	if err != nil {
		panic("I wansnt able to save: " + environment.Uuid.String())
	}
	return environment, nil
}

func GetAll() ([]*Environment, error) {
	items, err := ioutil.ReadDir(usedEnvironmentDirectory)
	if err != nil {
		return nil, err
	}
	var environments []*Environment
	for _, i := range items {
		if i.IsDir() {
			e, err := Get(uuid.MustParse(i.Name()))
			if err == nil {
				environments = append(environments, e)
			}
		}
	}
	return environments, nil
}

func Get(uuid uuid.UUID) (*Environment, error) {
	expectedFile := path.Join(usedEnvironmentDirectory, uuid.String(), util.DefaultEnvironmentConfigFileName)
	if _, err := os.Stat(expectedFile); os.IsNotExist(err) {
		fmt.Println("file " + expectedFile + " does not exist!") //TODO err?
		return nil, err
	} else {
		e, err := loadEnvironmentFromConfigFile(expectedFile)
		if err != nil {
			fmt.Println("incorrect file?") //TODO warn?
		}
		return e, nil
	}
}

func loadEnvironmentFromConfigFile(configPath string) (*Environment, error) {
	e := &Environment{}
	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	d := yaml.NewDecoder(file)
	if err := d.Decode(&e); err != nil {
		return nil, err
	}
	return e, nil
}
