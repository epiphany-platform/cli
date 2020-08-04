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

func (cc *InstalledComponentCommand) RunDocker(image string, workDirectory string, mountPath string, mounts []string) error {
	for _, m := range mounts {
		util.EnsureDirectory(path.Join(mountPath, m))
	}
	dockerJob := &docker.Job{
		Image:                image,
		Command:              cc.Command,
		Args:                 cc.Args,
		WorkDirectory:        workDirectory,
		Mounts:               mounts,
		MountPath:            mountPath,
		EnvironmentVariables: cc.Envs,
	}
	debug("will try to run docker job %+v", dockerJob)
	return dockerJob.Run()
}

func (cc *InstalledComponentCommand) String() string {
	return fmt.Sprintf("    Command:\n     Name %s\n     Description %s\n", cc.Name, cc.Description)
}

type InstalledComponentVersion struct {
	EnvironmentRef uuid.UUID                   `yaml:"environment_ref"` //TODO try to remove it
	Name           string                      `yaml:"name"`
	Type           string                      `yaml:"type"`
	Version        string                      `yaml:"version"`
	Image          string                      `yaml:"image"`
	WorkDirectory  string                      `yaml:"workdir"`
	Mounts         []string                    `yaml:"mounts"`
	Commands       []InstalledComponentCommand `yaml:"commands"`
}

func (cv *InstalledComponentVersion) Run(command string) error {
	if cv.Type == "docker" {
		mountPath := path.Join(
			usedEnvironmentDirectory,
			cv.EnvironmentRef.String(),
			cv.Name,
			cv.Version,
			util.DefaultComponentMountsSubdirectory,
		)
		for _, cc := range cv.Commands {
			if cc.Name == command {
				return cc.RunDocker(cv.Image, cv.WorkDirectory, mountPath, cv.Mounts)
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

func (cv *InstalledComponentVersion) PersistLogs(logs string) { //TODO change to zerolog
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
		errFailedToWriteFile(err)
	}
}

type Environment struct {
	Name      string                      `yaml:"name"`
	Uuid      uuid.UUID                   `yaml:"uuid"`
	Installed []InstalledComponentVersion `yaml:"installed"`
}

func (e *Environment) Save() error {
	debug("will try to marshal environment %+v", e)
	data, err := yaml.Marshal(e)
	if err != nil {
		return err
	}
	ep := path.Join(usedEnvironmentDirectory, e.Uuid.String(), util.DefaultEnvironmentConfigFileName)
	debug("will try to write marshaled data to file %s", ep)
	err = ioutil.WriteFile(ep, data, 0644)
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
			return errors.New("this version of component is already installed in environment")
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

//TODO move whole global variables initialization to one place
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
		errSaveEnvironment(err, environment.Uuid.String())
	}
	return environment, nil
}

func GetAll() ([]*Environment, error) {
	debug("will try to get all subdirectories of %s directory", usedEnvironmentDirectory)
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
			} else {
				warnNotEnvironmentDirectory(err)
			}
		}
	}
	return environments, nil
}

func Get(uuid uuid.UUID) (*Environment, error) {
	expectedFile := path.Join(usedEnvironmentDirectory, uuid.String(), util.DefaultEnvironmentConfigFileName)
	debug("will try to get environment config from file %s", expectedFile)
	if _, err := os.Stat(expectedFile); os.IsNotExist(err) {
		warnEnvironmentConfigFileNotFound(err, expectedFile)
		return nil, err
	} else {
		e, err := loadEnvironmentFromConfigFile(expectedFile)
		if err != nil {
			return nil, err
		}
		debug("got environment config %+v", e)
		return e, nil
	}
}

func loadEnvironmentFromConfigFile(configPath string) (*Environment, error) {
	e := &Environment{}
	debug("trying to open %s file", configPath)
	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	d := yaml.NewDecoder(file)
	debug("will try to decode file %s to yaml", configPath)
	if err := d.Decode(&e); err != nil {
		return nil, err
	}
	return e, nil
}
