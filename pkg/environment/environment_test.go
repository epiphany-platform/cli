package environment

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/epiphany-platform/cli/pkg/util"
	"github.com/google/uuid"
	"github.com/mholt/archiver/v3"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

const configFileTemplate = `version: v1
kind: Config
current-environment: %s
`

func setup(t *testing.T, suffix string) (string, string, string, string) {
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	parentDir := os.TempDir()

	// Create directory to store all configuration
	mainDirectory, err := ioutil.TempDir(parentDir, fmt.Sprintf("*-e-environment-%s", suffix))
	if err != nil {
		t.Fatal(err)
	}

	// Create environments subdirectory
	envsDirectory := path.Join(mainDirectory, util.DefaultEnvironmentsSubdirectory)
	err = os.Mkdir(envsDirectory, 0755)
	if err != nil {
		t.Fatal(err)
	}

	// Create temp directory
	tempDirectory := path.Join(mainDirectory, util.DefaultEnvironmentsTempSubdirectory)
	err = os.Mkdir(util.UsedTempDirectory, 0755)
	if err != nil {
		t.Fatal(err)
	}

	// Create cli config file
	configFile := path.Join(mainDirectory, util.DefaultConfigFileName)
	err = ioutil.WriteFile(configFile, []byte(fmt.Sprintf(configFileTemplate, uuid.Nil)), 0644)
	if err != nil {
		t.Fatal(err)
	}

	return configFile, mainDirectory, envsDirectory, tempDirectory
}

func prepareImportTestPrereqs(t *testing.T, usedConfigurationDirectory, usedEnvironmentDirectory string) (string, string, string, string) {
	// Create environments
	importEnvValidFirst, err := Create("export-valid1")
	if err != nil {
		t.Fatal(err)
	}
	importEnvValidSecond, err := Create("export-valid2")
	if err != nil {
		t.Fatal(err)
	}
	importEnvInvalidFirst, err := Create("export-invalid1")
	if err != nil {
		t.Fatal(err)
	}
	// Remove config file for one of environments to make it invalid
	os.Remove(path.Join(usedEnvironmentDirectory, importEnvInvalidFirst.Uuid.String(), util.DefaultEnvironmentConfigFileName))

	importEnvInvalidSecond, err := Create("export-invalid2")
	if err != nil {
		t.Fatal(err)
	}

	// Create archives
	// Create a valid archive
	importDirValidFirst := path.Join(usedEnvironmentDirectory, importEnvValidFirst.Uuid.String())
	importFileValidFirst := path.Join(usedConfigurationDirectory, importEnvValidFirst.Uuid.String()+".zip")
	err = archiver.Archive([]string{importDirValidFirst}, importFileValidFirst)
	if err != nil {
		t.Fatal(err)
	}

	// Create a valid archive with a changed name
	importDirValidSecond := path.Join(usedEnvironmentDirectory, importEnvValidSecond.Uuid.String())
	importFileValidSecond := path.Join(usedConfigurationDirectory, "changed.zip")
	err = archiver.Archive([]string{importDirValidSecond}, importFileValidSecond)
	if err != nil {
		t.Fatal(err)
	}

	// Create an invalid archive with missing env config file
	importDirInvalidFirst := path.Join(usedEnvironmentDirectory, importEnvInvalidFirst.Uuid.String())
	importFileInvalidFirst := path.Join(usedConfigurationDirectory, importEnvInvalidFirst.Uuid.String()+".zip")
	err = archiver.Archive([]string{importDirInvalidFirst}, importFileInvalidFirst)
	if err != nil {
		t.Fatal(err)
	}

	// Create a valid archive for existing environment
	importDirInvalidSecond := path.Join(usedEnvironmentDirectory, importEnvInvalidSecond.Uuid.String())
	importFileInvalidSecond := path.Join(usedConfigurationDirectory, importEnvInvalidSecond.Uuid.String()+".zip")
	err = archiver.Archive([]string{importDirInvalidSecond}, importFileInvalidSecond)
	if err != nil {
		t.Fatal(err)
	}

	// Remove environments to be able to export except one left by purpose
	os.RemoveAll(path.Join(usedEnvironmentDirectory, importEnvValidFirst.Uuid.String()))
	os.RemoveAll(path.Join(usedEnvironmentDirectory, importEnvValidSecond.Uuid.String()))
	os.RemoveAll(path.Join(usedEnvironmentDirectory, importEnvInvalidFirst.Uuid.String()))

	return importFileValidFirst, importFileValidSecond, importFileInvalidFirst, importFileInvalidSecond
}

func TestGet(t *testing.T) {
	util.UsedConfigFile, util.UsedConfigurationDirectory, util.UsedEnvironmentDirectory, _ = setup(t, "get")
	defer os.RemoveAll(util.UsedConfigurationDirectory)

	type args struct {
		uuid uuid.UUID
	}
	tests := []struct {
		name      string
		args      args
		mocked    []byte
		want      *Environment
		wantErr   error
		isPattern bool
	}{
		{
			name: "correct",
			args: args{
				uuid: uuid.MustParse("fccf6810-32c4-4500-9414-2de45d2c4097"),
			},
			mocked: []byte(`name: e1
uuid: fccf6810-32c4-4500-9414-2de45d2c4097
installed: []`),
			want: &Environment{
				Name:      "e1",
				Uuid:      uuid.MustParse("fccf6810-32c4-4500-9414-2de45d2c4097"),
				Installed: []InstalledComponentVersion{},
			},
			wantErr:   nil,
			isPattern: false,
		},
		{
			name: "missing file",
			args: args{
				uuid: uuid.MustParse("816789aa-7839-4f2b-ac74-b66344e4fbe8"),
			},
			wantErr:   errors.New("no such file or directory"),
			isPattern: true,
		},
		{
			name: "incorrect file",
			args: args{
				uuid: uuid.MustParse("66d4cd70-4375-4737-b6ce-7e13f3cc93f9"),
			},
			mocked:    []byte(`incorrect file`),
			wantErr:   errors.New("yaml: unmarshal errors:\n  line 1: cannot unmarshal !!str `incorre...` into environment.Environment"),
			isPattern: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := assert.New(t)
			if len(tt.mocked) > 0 {
				envDir := path.Join(util.UsedEnvironmentDirectory, tt.args.uuid.String())
				err := os.MkdirAll(envDir, 0755)
				if err != nil {
					t.Fatal(err)
				}
				envConfigFile := path.Join(envDir, util.DefaultEnvironmentConfigFileName)
				err = ioutil.WriteFile(envConfigFile, tt.mocked, 0644)
				if err != nil {
					t.Fatal(err)
				}
			}
			got, err := Get(tt.args.uuid)
			if tt.wantErr == nil {
				a.NoError(err)
			} else if tt.isPattern {
				a.Contains(err.Error(), tt.wantErr.Error())
			} else {
				a.EqualError(err, tt.wantErr.Error())
			}
			a.Truef(reflect.DeepEqual(got, tt.want), "got = %#v, want = %#v", got, tt.want)
		})
	}
}

func TestGetAll(t *testing.T) {
	type mocked struct {
		subdirectory  string
		configName    string
		configContent []byte
	}
	tests := []struct {
		name        string
		mocked      []mocked
		want        []*Environment
		wantErr     error
		shouldPanic bool
	}{
		{
			name: "correct",
			mocked: []mocked{
				{
					subdirectory: "45764648-162a-4526-bdd0-71a438fd6ceb",
					configName:   "config.yaml",
					configContent: []byte(`name: e2
uuid: 45764648-162a-4526-bdd0-71a438fd6ceb
installed: []`),
				},
				{
					subdirectory: "4af1705c-c48f-4ca2-be08-53c673da835c",
					configName:   "config.yaml",
					configContent: []byte(`name: e1
uuid: 4af1705c-c48f-4ca2-be08-53c673da835c
installed: []`),
				},
			},
			want: []*Environment{
				{
					Name:      "e2",
					Uuid:      uuid.MustParse("45764648-162a-4526-bdd0-71a438fd6ceb"),
					Installed: []InstalledComponentVersion{},
				},
				{
					Name:      "e1",
					Uuid:      uuid.MustParse("4af1705c-c48f-4ca2-be08-53c673da835c"),
					Installed: []InstalledComponentVersion{},
				},
			},
			wantErr:     nil,
			shouldPanic: false,
		},
		{
			name: "subdirectory name not uuid",
			mocked: []mocked{
				{
					subdirectory: "45764648-162a-4526-bdd0-71a438fd6ceb",
					configName:   "config.yaml",
					configContent: []byte(`name: e2
uuid: 45764648-162a-4526-bdd0-71a438fd6ceb
installed: []`),
				},
				{
					subdirectory: "incorrect-directory-name",
					configName:   "config.yaml",
					configContent: []byte(`name: e1
uuid: 4af1705c-c48f-4ca2-be08-53c673da835c
installed: []`),
				},
			},
			wantErr:     errors.New("uuid: Parse(incorrect-directory-name): invalid UUID length: 24"),
			shouldPanic: true,
		},
		{
			name: "incorrect config file name",
			mocked: []mocked{
				{
					subdirectory: "45764648-162a-4526-bdd0-71a438fd6ceb",
					configName:   "config.yaml",
					configContent: []byte(`name: e2
uuid: 45764648-162a-4526-bdd0-71a438fd6ceb
installed: []`),
				},
				{
					subdirectory: "4af1705c-c48f-4ca2-be08-53c673da835c",
					configName:   "incorrect-config.yaml",
					configContent: []byte(`name: e1
uuid: 4af1705c-c48f-4ca2-be08-53c673da835c
installed: []`),
				},
			},
			want: []*Environment{
				{
					Name:      "e2",
					Uuid:      uuid.MustParse("45764648-162a-4526-bdd0-71a438fd6ceb"),
					Installed: []InstalledComponentVersion{},
				},
			},
			wantErr:     nil,
			shouldPanic: false,
		},
		{
			name: "incorrect config file content",
			mocked: []mocked{
				{
					subdirectory: "45764648-162a-4526-bdd0-71a438fd6ceb",
					configName:   "config.yaml",
					configContent: []byte(`name: e2
uuid: 45764648-162a-4526-bdd0-71a438fd6ceb
installed: []`),
				},
				{
					subdirectory:  "4af1705c-c48f-4ca2-be08-53c673da835c",
					configName:    "config.yaml",
					configContent: []byte(`incorrect content`),
				},
			},
			want: []*Environment{
				{
					Name:      "e2",
					Uuid:      uuid.MustParse("45764648-162a-4526-bdd0-71a438fd6ceb"),
					Installed: []InstalledComponentVersion{},
				},
			},
			wantErr:     nil,
			shouldPanic: false,
		},
	}
	for _, tt := range tests {
		util.UsedConfigFile, util.UsedConfigurationDirectory, util.UsedEnvironmentDirectory, _ = setup(t, "get-all")
		t.Run(tt.name, func(t *testing.T) {
			a := assert.New(t)
			if tt.mocked != nil && len(tt.mocked) > 0 {
				for _, m := range tt.mocked {
					envDir := path.Join(util.UsedEnvironmentDirectory, m.subdirectory)
					err := os.MkdirAll(envDir, 0755)
					if err != nil {
						t.Fatal(err)
					}
					envConfigFile := path.Join(envDir, m.configName)
					err = ioutil.WriteFile(envConfigFile, m.configContent, 0644)
					if err != nil {
						t.Fatal(err)
					}
				}
			}

			if tt.shouldPanic {
				a.PanicsWithValue(tt.wantErr.Error(), func() { GetAll() })
			} else {
				got, err := GetAll()
				if tt.wantErr == nil {
					a.NoError(err)
				} else {
					a.EqualError(err, tt.wantErr.Error())
				}
				a.Truef(reflect.DeepEqual(got, tt.want), "got = %#v, want = %#v", got, tt.want)
			}
		})
		os.RemoveAll(util.UsedConfigurationDirectory)
	}
}

func Test_create(t *testing.T) {
	util.UsedConfigFile, util.UsedConfigurationDirectory, util.UsedEnvironmentDirectory, _ = setup(t, "create")
	defer os.RemoveAll(util.UsedConfigurationDirectory)

	type args struct {
		name string
		uuid string
	}
	tests := []struct {
		name    string
		args    args
		want    *Environment
		wantErr error
	}{
		{
			name: "correct",
			args: args{
				name: "e1",
				uuid: "b03bb900-5d49-4421-a45e-eeeb40e0a5d5",
			},
			want: &Environment{
				Name: "e1",
				Uuid: uuid.MustParse("b03bb900-5d49-4421-a45e-eeeb40e0a5d5"),
			},
			wantErr: nil,
		},
		{
			name: "empty name",
			args: args{
				name: "",
				uuid: "66d4cd70-4375-4737-b6ce-7e13f3cc93f9",
			},
			want: &Environment{
				Uuid: uuid.MustParse("66d4cd70-4375-4737-b6ce-7e13f3cc93f9"),
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := assert.New(t)
			got, err := create(tt.args.name, uuid.MustParse(tt.args.uuid))
			if tt.wantErr == nil {
				a.NoError(err)
			} else {
				a.EqualError(err, tt.wantErr.Error())
			}
			a.Truef(reflect.DeepEqual(got, tt.want), "got = %#v, want = %#v", got, tt.want)
			expectedConfigFile := path.Join(util.UsedEnvironmentDirectory, got.Uuid.String(), util.DefaultEnvironmentConfigFileName)
			a.FileExistsf(expectedConfigFile, "expected to find file %s but didn't find", expectedConfigFile)
		})
	}
}

func TestEnvironment_Save(t *testing.T) {
	util.UsedConfigFile, util.UsedConfigurationDirectory, util.UsedEnvironmentDirectory, _ = setup(t, "env-save")
	defer os.RemoveAll(util.UsedConfigurationDirectory)

	tests := []struct {
		name        string
		environment *Environment
		wantContent []byte
		wantErr     error
	}{
		{
			name: "correct",
			environment: &Environment{
				Name: "e1",
				Uuid: uuid.MustParse("10d52c05-029e-4794-a790-79d6c2af40b6"),
			},
			wantContent: []byte(`name: e1
uuid: 10d52c05-029e-4794-a790-79d6c2af40b6
installed: []
`),
			wantErr: nil,
		},
		{
			name: "missing uuid",
			environment: &Environment{
				Name: "e1",
			},
			wantErr: errors.New("unexpected UUID on Save: 00000000-0000-0000-0000-000000000000"),
		},
		{
			name: "missing name",
			environment: &Environment{
				Uuid: uuid.MustParse("10d52c05-029e-4794-a790-79d6c2af40b6"),
			},
			wantContent: []byte(`name: ""
uuid: 10d52c05-029e-4794-a790-79d6c2af40b6
installed: []
`),
			wantErr: nil,
		},
		{
			name: "with installed components",
			environment: &Environment{
				Name: "x",
				Uuid: uuid.MustParse("3e5b7269-1b3d-4003-9454-9f472857633a"),
				Installed: []InstalledComponentVersion{
					{
						EnvironmentRef: uuid.MustParse("3e5b7269-1b3d-4003-9454-9f472857633a"),
						Name:           "x",
						Type:           "x",
						Version:        "x",
						Image:          "x",
						WorkDirectory:  "x",
						Mounts:         []string{"x"},
						Shared:         "x",
						Commands:       []InstalledComponentCommand{},
					},
				},
			},
			wantContent: []byte(`name: x
uuid: 3e5b7269-1b3d-4003-9454-9f472857633a
installed:
- environment_ref: 3e5b7269-1b3d-4003-9454-9f472857633a
  name: x
  type: x
  version: x
  image: x
  workdir: x
  mounts:
  - x
  shared: x
  commands: []
`),
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := assert.New(t)
			dir := path.Join(util.UsedEnvironmentDirectory, tt.environment.Uuid.String())
			err := os.MkdirAll(dir, 0755)
			if err != nil {
				t.Fatal(err)
			}
			e := &Environment{
				Name:      tt.environment.Name,
				Uuid:      tt.environment.Uuid,
				Installed: tt.environment.Installed,
			}
			err = e.Save()
			if tt.wantErr == nil {
				a.NoError(err)
			} else {
				a.EqualError(err, tt.wantErr.Error())
			}
			if len(tt.wantContent) > 0 {
				expectedConfigFile := path.Join(dir, util.DefaultEnvironmentConfigFileName)
				a.FileExistsf(expectedConfigFile, "expected to find file %s but didn't find", expectedConfigFile)
				savedBytes, _ := ioutil.ReadFile(expectedConfigFile)
				a.Truef(bytes.Equal(tt.wantContent, savedBytes), "saved file is \n%s\n\n but expected is \n\n%s\n", string(savedBytes), string(tt.wantContent))
			}
		})
	}
}

func TestEnvironment_GetComponentByName(t *testing.T) {
	util.UsedConfigFile, util.UsedConfigurationDirectory, util.UsedEnvironmentDirectory, _ = setup(t, "env-get-by-name")
	defer os.RemoveAll(util.UsedConfigurationDirectory)

	tests := []struct {
		name          string
		environment   *Environment
		componentName string
		want          *InstalledComponentVersion
		wantErr       error
	}{
		{
			name: "correct single",
			environment: &Environment{
				Name: "e1",
				Uuid: uuid.MustParse("10d52c05-029e-4794-a790-79d6c2af40b6"),
				Installed: []InstalledComponentVersion{
					{
						EnvironmentRef: uuid.MustParse("10d52c05-029e-4794-a790-79d6c2af40b6"),
						Name:           "c1",
						Type:           "docker",
						Version:        "v1",
						Image:          "i1",
					},
				},
			},
			componentName: "c1",
			want: &InstalledComponentVersion{
				EnvironmentRef: uuid.MustParse("10d52c05-029e-4794-a790-79d6c2af40b6"),
				Name:           "c1",
				Type:           "docker",
				Version:        "v1",
				Image:          "i1",
			},
			wantErr: nil,
		},
		{
			name: "correct multiple",
			environment: &Environment{
				Name: "e1",
				Uuid: uuid.MustParse("10d52c05-029e-4794-a790-79d6c2af40b6"),
				Installed: []InstalledComponentVersion{
					{
						EnvironmentRef: uuid.MustParse("10d52c05-029e-4794-a790-79d6c2af40b6"),
						Name:           "c1",
						Type:           "d2",
						Version:        "v1",
						Image:          "i1",
					},
					{
						EnvironmentRef: uuid.MustParse("10d52c05-029e-4794-a790-79d6c2af40b6"),
						Name:           "c2",
						Type:           "d2",
						Version:        "v2",
						Image:          "i3",
					},
				},
			},
			componentName: "c1",
			want: &InstalledComponentVersion{
				EnvironmentRef: uuid.MustParse("10d52c05-029e-4794-a790-79d6c2af40b6"),
				Name:           "c1",
				Type:           "d2",
				Version:        "v1",
				Image:          "i1",
			},
			wantErr: nil,
		},
		{
			name: "missing",
			environment: &Environment{
				Name: "e1",
				Uuid: uuid.MustParse("10d52c05-029e-4794-a790-79d6c2af40b6"),
				Installed: []InstalledComponentVersion{
					{
						EnvironmentRef: uuid.MustParse("10d52c05-029e-4794-a790-79d6c2af40b6"),
						Name:           "c2",
						Type:           "d2",
						Version:        "v2",
						Image:          "i3",
					},
				},
			},
			componentName: "c1",
			wantErr:       errors.New("no such component installed"),
		},
		{
			name: "empty",
			environment: &Environment{
				Name:      "e1",
				Uuid:      uuid.MustParse("10d52c05-029e-4794-a790-79d6c2af40b6"),
				Installed: []InstalledComponentVersion{},
			},
			componentName: "c1",
			wantErr:       errors.New("no such component installed"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := assert.New(t)
			e := &Environment{
				Name:      tt.environment.Name,
				Uuid:      tt.environment.Uuid,
				Installed: tt.environment.Installed,
			}
			got, err := e.GetComponentByName(tt.componentName)
			if tt.wantErr == nil {
				a.NoError(err)
			} else {
				a.EqualError(err, tt.wantErr.Error())
			}
			a.Truef(reflect.DeepEqual(got, tt.want), "got = %#v, want = %#v", got, tt.want)
		})
	}
}

func TestIsExisting(t *testing.T) {
	util.UsedConfigFile, util.UsedConfigurationDirectory, util.UsedEnvironmentDirectory, _ = setup(t, "validation")
	validEnv, err := Create("validation")
	if err != nil {
		t.Fatal(err)
	}
	fakeEnvID, err := uuid.NewRandom()
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(util.UsedConfigurationDirectory)

	tests := []struct {
		name string
		uuid uuid.UUID
		want bool
	}{
		{
			name: "Existing environment",
			uuid: validEnv.Uuid,
			want: true,
		},
		{
			name: "Fake environment",
			uuid: fakeEnvID,
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			IsExisting, err := IsExisting(tt.uuid)
			a := assert.New(t)
			if a.NoError(err) {
				a.Equal(tt.want, IsExisting)
			}
		})
	}
}

func TestExport(t *testing.T) {
	util.UsedConfigFile, util.UsedConfigurationDirectory, util.UsedEnvironmentDirectory, util.UsedTempDirectory = setup(t, "export")
	defer os.RemoveAll(util.UsedConfigurationDirectory)

	// Create the environment to export
	exportEnv, err := Create("export")
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name       string
		destDir    string
		fileExists bool
		wantErr    error
		isPattern  bool
	}{
		{
			name:       "Successful export",
			destDir:    util.UsedConfigurationDirectory,
			fileExists: false,
			wantErr:    nil,
			isPattern:  false,
		},
		{
			name:       "Target file already exists",
			destDir:    util.UsedConfigurationDirectory,
			fileExists: true,
			wantErr:    errors.New("file already exists"),
			isPattern:  true,
		},
		{
			name:       "Writable path that does not exist",
			destDir:    path.Join(util.UsedConfigurationDirectory, "fake"),
			fileExists: false,
			wantErr:    nil,
			isPattern:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := assert.New(t)
			if tt.fileExists {
				_, err := os.Create(path.Join(tt.destDir, exportEnv.Uuid.String()+".zip"))
				if err != nil {
					t.Fatal(err)
				}
			}
			err = exportEnv.Export(tt.destDir)
			if tt.wantErr == nil {
				a.NoError(err)
				a.FileExists(path.Join(tt.destDir, exportEnv.Uuid.String()+".zip"))
			} else if tt.isPattern {
				a.Contains(err.Error(), tt.wantErr.Error())
			} else {
				a.EqualError(err, tt.wantErr.Error())
			}
		})
	}
}

func TestImport(t *testing.T) {
	util.UsedConfigFile, util.UsedConfigurationDirectory, util.UsedEnvironmentDirectory, _ = setup(t, "import")
	defer os.RemoveAll(util.UsedConfigurationDirectory)

	importFileValidFirst, importFileValidSecond, importFileInvalidFirst, importFileInvalidSecond := prepareImportTestPrereqs(t, util.UsedConfigurationDirectory, util.UsedEnvironmentDirectory)

	tests := []struct {
		name    string
		from    string
		wantErr error
	}{
		{
			name:    "Valid source file",
			from:    importFileValidFirst,
			wantErr: nil,
		},
		{
			name:    "Valid source file with a changed name",
			from:    importFileValidSecond,
			wantErr: nil,
		},
		{
			name:    "Source file that does not contain environment's config file",
			from:    importFileInvalidFirst,
			wantErr: errors.New("Missing environment config file"),
		},
		{
			name:    "Existing environment",
			from:    importFileInvalidSecond,
			wantErr: fmt.Errorf("Environment with id %s already exists", strings.Trim(filepath.Base(importFileInvalidSecond), filepath.Ext(importFileInvalidSecond))),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := assert.New(t)
			got, err := Import(tt.from)
			if tt.wantErr == nil {
				a.NoError(err)
				a.DirExists(path.Join(util.UsedEnvironmentDirectory, got.String()))
				a.FileExists(path.Join(util.UsedEnvironmentDirectory, got.String(), util.DefaultEnvironmentConfigFileName))
			} else {
				a.EqualError(err, tt.wantErr.Error())
			}

		})
	}
}
