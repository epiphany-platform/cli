package environment

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"time"

	"github.com/epiphany-platform/cli/pkg/auth"
	"github.com/epiphany-platform/cli/pkg/docker"
	"github.com/epiphany-platform/cli/pkg/util"
	"github.com/google/uuid"
	"github.com/mholt/archiver/v3"
	"gopkg.in/yaml.v2"
)

//InstalledComponentCommand holds information about specific command of installed component
type InstalledComponentCommand struct {
	Name        string            `yaml:"name"`
	Description string            `yaml:"description"`
	Command     string            `yaml:"command"`
	Envs        map[string]string `yaml:"envs"`
	Args        []string          `yaml:"args"`
}

//TODO add tests
func (cc *InstalledComponentCommand) RunDocker(image string, workDirectory string, mounts map[string]string, processor func(string) string) error {
	for _, v := range mounts {
		util.EnsureDirectory(v)
	}
	envs := make(map[string]string)
	for k, v := range cc.Envs {
		envs[k] = processor(v)
	}
	args := make([]string, len(cc.Args))
	for _, a := range cc.Args {
		args = append(args, processor(a))
	}
	dockerJob := &docker.Job{
		Image:                image,
		Command:              cc.Command,
		Args:                 args,
		WorkDirectory:        workDirectory,
		Mounts:               mounts,
		EnvironmentVariables: envs,
	}
	debug("will try to run docker job %+v", dockerJob)
	return dockerJob.Run()
}

//The String method is used to pretty-print InstalledComponentCommand struct
func (cc *InstalledComponentCommand) String() string {
	return fmt.Sprintf("    Command:\n     Name %s\n     Description %s\n", cc.Name, cc.Description)
}

//InstalledComponentVersion struct holds information about installed components with its details.
type InstalledComponentVersion struct {
	EnvironmentRef uuid.UUID                   `yaml:"environment_ref"` //TODO try to remove it
	Name           string                      `yaml:"name"`
	Type           string                      `yaml:"type"`
	Version        string                      `yaml:"version"`
	Image          string                      `yaml:"image"`
	WorkDirectory  string                      `yaml:"workdir"`
	Mounts         []string                    `yaml:"mounts"`
	Shared         string                      `yaml:"shared"`
	Commands       []InstalledComponentCommand `yaml:"commands"`
}

//TODO add tests
func (cv *InstalledComponentVersion) Run(command string, processor func(string) string) error {
	if cv.Type == "docker" {
		mounts := make(map[string]string)
		moduleMountPath := path.Join(
			util.UsedEnvironmentDirectory,
			cv.EnvironmentRef.String(),
			cv.Name,
			cv.Version,
			util.DefaultComponentMountsSubdirectory,
		)
		for _, m := range cv.Mounts {
			mounts[m] = path.Join(moduleMountPath, m)
		}
		if cv.Shared != "" {
			mounts[cv.Shared] = path.Join(
				util.UsedEnvironmentDirectory,
				cv.EnvironmentRef.String(),
				"/shared", //TODO to consts
			)
		}
		for _, cc := range cv.Commands {
			if cc.Name == command {
				return cc.RunDocker(cv.Image, cv.WorkDirectory, mounts, processor)
			}
		}
	}
	return errors.New("nothing to run for this version")
}

//The String method is used to pretty-print InstalledComponentVersion struct
func (cv *InstalledComponentVersion) String() string {
	var b bytes.Buffer
	b.WriteString(fmt.Sprintf("  Installed Component:\n   Name: %s\n   Type: %s\n   Version: %s\n   Image: %s\n", cv.Name, cv.Type, cv.Version, cv.Image))
	for _, cc := range cv.Commands {
		b.WriteString(cc.String())
	}
	return b.String()
}

//TODO add tests
func (cv *InstalledComponentVersion) Download() error {
	if cv.Type == "docker" {
		dockerImage := &docker.Image{Name: cv.Image}
		found, err := dockerImage.IsPulled()
		if err != nil {
			return err
		}
		if found {
			logger.Debug().Msg("image is already present, no need to download") //TODO consider --force-download switch
			return nil
		}
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
		util.UsedEnvironmentDirectory,
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

type SshConfig struct {
	RsaKeyPair auth.RsaKeyPair `yaml:"rsa-keypair"`
}

//Environment struct holds all information about managed environment with list of InstalledComponentVersion
type Environment struct {
	Name      string                      `yaml:"name"`
	Uuid      uuid.UUID                   `yaml:"uuid"`
	Installed []InstalledComponentVersion `yaml:"installed"`
	SshConfig SshConfig                   `yaml:"ssh-config,omitempty"`
}

//Save updated Environment to file
func (e *Environment) Save() error {
	if e.Uuid == uuid.Nil {
		return errors.New(fmt.Sprintf("unexpected UUID on Save: %s", e.Uuid))
	}
	debug("will try to marshal environment %+v", e)
	data, err := yaml.Marshal(e)
	if err != nil {
		return err
	}
	ep := path.Join(util.UsedEnvironmentDirectory, e.Uuid.String(), util.DefaultEnvironmentConfigFileName)
	debug("will try to write marshaled data to file %s", ep)
	err = ioutil.WriteFile(ep, data, 0644)
	if err != nil {
		return err
	}
	return nil
}

//The String method is used to pretty-print Environment struct
func (e *Environment) String() string {
	var b bytes.Buffer
	b.WriteString("Environment info:\n")
	b.WriteString(fmt.Sprintf(" Name: %s\n UUID: %s\n", e.Name, e.Uuid.String()))
	for _, ic := range e.Installed {
		b.WriteString(ic.String())
	}
	return b.String()
}

//TODO add tests
func (e *Environment) Install(newComponent InstalledComponentVersion) error {
	for _, ic := range e.Installed {
		if ic.Name == newComponent.Name && ic.Version == newComponent.Version {
			return errors.New("this version of component is already installed in environment")
		}
	}
	e.Installed = append(e.Installed, newComponent)
	newComponentRunsDirectory := path.Join(util.UsedEnvironmentDirectory, e.Uuid.String(), newComponent.Name, newComponent.Version, util.DefaultComponentRunsSubdirectory)
	newComponentMountsDirectory := path.Join(util.UsedEnvironmentDirectory, e.Uuid.String(), newComponent.Name, newComponent.Version, util.DefaultComponentMountsSubdirectory)
	util.EnsureDirectory(newComponentRunsDirectory)
	util.EnsureDirectory(newComponentMountsDirectory)
	err := newComponent.Download()
	if err != nil {
		return err
	}
	return e.Save()
}

//GetComponentByName returns first InstalledComponentVersion found by name
func (e *Environment) GetComponentByName(name string) (*InstalledComponentVersion, error) {
	for _, ic := range e.Installed {
		if ic.Name == name {
			return &ic, nil
		}
	}
	return nil, errors.New("no such component installed")
}

func (e *Environment) AddRsaKeyPair(rsaKeyPair auth.RsaKeyPair) {
	e.SshConfig.RsaKeyPair = rsaKeyPair
}

//Create new environment with given name
func Create(name string) (*Environment, error) {
	return create(name, uuid.New())
}

//create new environment with given name and uuid
func create(name string, uuid uuid.UUID) (*Environment, error) {
	debug("will try to create environment with uuid %s and name %s", uuid.String(), name)
	environment := &Environment{
		Name: name,
		Uuid: uuid,
	}
	newEnvironmentDirectory := path.Join(util.UsedEnvironmentDirectory, environment.Uuid.String())
	util.EnsureDirectory(newEnvironmentDirectory)
	err := environment.Save()
	if err != nil {
		errSaveEnvironment(err, environment.Uuid.String())
	}
	return environment, nil
}

//GetAll existing Environment
func GetAll() ([]*Environment, error) {
	debug("will try to get all subdirectories of %s directory", util.UsedEnvironmentDirectory)
	items, err := ioutil.ReadDir(util.UsedEnvironmentDirectory)
	if err != nil {
		return nil, err
	}
	var environments []*Environment
	for _, i := range items {
		debug("entered directory %s", i.Name())
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

//Get Environment bu uuid
func Get(uuid uuid.UUID) (*Environment, error) {
	expectedFile := path.Join(util.UsedEnvironmentDirectory, uuid.String(), util.DefaultEnvironmentConfigFileName)
	debug("will try to get environment config from file %s", expectedFile)
	if _, err := os.Stat(expectedFile); os.IsNotExist(err) {
		warnEnvironmentConfigFileNotFound(err, expectedFile)
		return nil, err
	} else {
		e := &Environment{}
		debug("trying to open %s file", expectedFile)
		file, err := os.Open(expectedFile)
		if err != nil {
			return nil, err
		}
		defer file.Close()
		d := yaml.NewDecoder(file)
		debug("will try to decode file %s to yaml", expectedFile)
		if err := d.Decode(&e); err != nil {
			return nil, err
		}
		debug("got environment config %+v", e)
		return e, nil
	}
}

//TODO add tests
func IsValid(envID string) (bool, error) {
	environments, err := GetAll()
	if err != nil {
		return false, err
	}
	isEnvValid := false
	for _, e := range environments {
		if e.Uuid.String() == envID {
			isEnvValid = true
			logger.Debug().Msgf("Found environment to export: %s", e.Uuid.String())
			break
		}
	}
	return isEnvValid, nil
}

//TODO add tests
func Export(envID, dstDir string) error {
	envPath := path.Join(util.UsedEnvironmentDirectory, envID)

	// Final archive name is envID + .zip extension
	err := archiver.Archive([]string{envPath}, path.Join(dstDir, envID+".zip"))
	if err != nil {
		return err
	}
	return nil
}

//TODO add tests
func Import(srcFile string) (*uuid.UUID, error) {
	// Check if environment config exists in zip archive
	// before export and verify its content
	var envConfig *Environment
	err := archiver.Walk(srcFile, func(f archiver.File) error {
		if f.Name() == util.DefaultConfigFileName {
			configContent, err := ioutil.ReadAll(f)
			if err != nil {
				return errors.New("Unable to read environment config")
			}
			envConfig = &Environment{}
			err = yaml.Unmarshal(configContent, envConfig)
			if err != nil {
				return errors.New("Cannot unmarshal config")
			}
			if envConfig.Uuid == uuid.Nil {
				return errors.New("Environment id is missing in the config")
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	envConfigDir := path.Join(util.GetHomeDirectory(), util.DefaultConfigurationDirectory, util.DefaultEnvironmentsSubdirectory)

	// Check if environment with such id is already in place
	if _, err := os.Stat(path.Join(envConfigDir, envConfig.Uuid.String())); err == nil {
		return nil, fmt.Errorf("Environment with id %s already exists", envConfig.Uuid.String())
	}

	// Unarchive specified file
	err = archiver.Unarchive(srcFile, envConfigDir)
	if err != nil {
		return nil, err
	}

	// Download all Docker images for installed components
	for _, cmp := range envConfig.Installed {
		err = cmp.Download()
		if err != nil {
			return nil, err
		}
	}

	return &envConfig.Uuid, nil
}
