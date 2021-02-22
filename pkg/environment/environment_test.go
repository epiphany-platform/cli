package environment

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"regexp"
	"testing"

	"github.com/epiphany-platform/cli/pkg/util"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

const (
	envIDValid         = "8f327570-5401-402f-ac57-8f6cd56315e7"
	configFileTemplate = `version: v1
kind: Config
current-environment: %s
`
	envConfigFileTemplate = `name: %s
uuid: %s
`
)

func setup(t *testing.T, suffix string, createEnv bool) (string, string, string) {
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	parentDir := os.TempDir()

	// Create directory to store all configuration
	mainDirectory, err := ioutil.TempDir(parentDir, fmt.Sprintf("*-e-environment-%s", suffix))
	if err != nil {
		t.Fatal(err)
	}

	// Create environments supdirectory
	envsDirectory := path.Join(mainDirectory, util.DefaultEnvironmentsSubdirectory)
	os.Mkdir(envsDirectory, 0755)

	// Create cli config file
	configFile := path.Join(mainDirectory, util.DefaultConfigFileName)
	err = ioutil.WriteFile(configFile, []byte(fmt.Sprintf(configFileTemplate, envIDValid)), 0644)
	if err != nil {
		t.Fatal(err)
	}

	if createEnv {
		// Create config directory for a valid environment
		validEnvDir := path.Join(envsDirectory, envIDValid)
		os.Mkdir(validEnvDir, 0755)

		// Create config file for a valid environment
		envConfigFile := path.Join(validEnvDir, util.DefaultEnvironmentConfigFileName)
		err = ioutil.WriteFile(envConfigFile, []byte(fmt.Sprintf(envConfigFileTemplate, "valid", envIDValid)), 0644)
		if err != nil {
			t.Fatal(err)
		}
	}

	return configFile, mainDirectory, envsDirectory
}

func TestGet(t *testing.T) {
	util.UsedConfigFile, util.UsedConfigurationDirectory, util.UsedEnvironmentDirectory = setup(t, "get", false)
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
		util.UsedConfigFile, util.UsedConfigurationDirectory, util.UsedEnvironmentDirectory = setup(t, "get-all", false)
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
	util.UsedConfigFile, util.UsedConfigurationDirectory, util.UsedEnvironmentDirectory = setup(t, "create", false)
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
	util.UsedConfigFile, util.UsedConfigurationDirectory, util.UsedEnvironmentDirectory = setup(t, "env-save", false)
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
	util.UsedConfigFile, util.UsedConfigurationDirectory, util.UsedEnvironmentDirectory = setup(t, "env-get-by-name", false)
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

func TestIsValid(t *testing.T) {
	util.UsedConfigFile, util.UsedConfigurationDirectory, util.UsedEnvironmentDirectory = setup(t, "validation", true)
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
			uuid: uuid.MustParse(envIDValid),
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
			isValid, err := IsValid(tt.uuid)
			a := assert.New(t)
			if a.NoError(err) {
				a.Equal(tt.want, isValid)
			}
		})
	}
}

func isWrongResult(t *testing.T, err error, wantErr error) bool {
	if err != nil && wantErr != nil {
		re := regexp.MustCompile(wantErr.Error())
		if !re.MatchString(err.Error()) {
			t.Errorf("got \n%v\n, want \n%v\n", err, wantErr)
			return true
		}
	} else if err == nil && wantErr != nil {
		t.Errorf("didn't got error but want: %v", wantErr)
		return true
	} else if err != nil && wantErr == nil {
		t.Errorf("didnt want error but got: %v", err)
		return true
	}
	return false
}

func assertPanic(t *testing.T, f func(), wantErr error) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		} else {
			if rc, ok := r.(string); ok {
				if isWrongResult(t, errors.New(rc), wantErr) {
					return
				}
			} else {
				t.Errorf("cannot cast recover resutl to string. r = %#v", r)
				return
			}
		}
	}()
	f()
}
