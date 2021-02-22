package configuration

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"testing"

	"github.com/epiphany-platform/cli/pkg/az"
	"github.com/epiphany-platform/cli/pkg/util"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

func setup(t *testing.T, suffix string) (string, string) {
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	parentDir := os.TempDir()
	mainDirectory, err := ioutil.TempDir(parentDir, fmt.Sprintf("*-e-configuration-%s", suffix))
	if err != nil {
		t.Fatal(err)
	}

	tempFile, err := ioutil.TempFile(mainDirectory, "config-*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	return tempFile.Name(), mainDirectory
}

func TestConfig_GetConfigFilePath(t *testing.T) {
	tempFile, tempDirectory := setup(t, "get")
	defer os.RemoveAll(tempDirectory)

	tests := []struct {
		name        string
		mocked      string
		want        string
		shouldPanic bool
	}{
		{
			name:        "correct",
			mocked:      tempFile,
			want:        tempFile,
			shouldPanic: false,
		},
		{
			name:        "incorrect",
			mocked:      "",
			want:        tempFile,
			shouldPanic: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			util.UsedConfigFile = tt.mocked
			c := &Config{}
			f := func() {
				if got := c.GetConfigFilePath(); got != tt.want {
					t.Errorf("got %s, want %s", got, tt.want)
				}
			}
			if tt.shouldPanic {
				assertPanic(t, f)
			} else {
				f()
			}
		})
	}
}

func TestConfig_SetUsedEnvironment(t *testing.T) {
	tempFile, tempDirectory := setup(t, "used")
	envIDCurrent := "3e5b7269-1b3d-4003-9454-9f472857633a"
	envIDSwitch := "567c0831-7e83-4b56-a2a7-ec7a8327238f"
	envConfigTemplate := `name: %s
uuid: %s
`
	defer os.RemoveAll(tempDirectory)

	type fields struct {
		Version            string
		Kind               Kind
		CurrentEnvironment uuid.UUID
	}

	tests := []struct {
		name       string
		fields     fields
		uuid       uuid.UUID
		configPath string
		wantErr    error
		want       []byte
	}{
		{
			name: "nil to some",
			fields: fields{
				Version:            "v1",
				Kind:               KindConfig,
				CurrentEnvironment: uuid.Nil,
			},
			uuid:       uuid.MustParse(envIDSwitch),
			configPath: tempFile,
			wantErr:    nil,
			want: []byte(fmt.Sprintf(`version: v1
kind: Config
current-environment: %s
`, envIDSwitch)),
		},
		{
			name: "some to another",
			fields: fields{
				Version:            "v1",
				Kind:               KindConfig,
				CurrentEnvironment: uuid.MustParse(envIDCurrent),
			},
			uuid:       uuid.MustParse(envIDSwitch),
			configPath: tempFile,
			wantErr:    nil,
			want: []byte(fmt.Sprintf(`version: v1
kind: Config
current-environment: %s
`, envIDSwitch)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create necessary files and directories as environments are validated before switching
			util.UsedConfigurationDirectory = tempDirectory
			util.UsedEnvironmentDirectory = path.Join(util.UsedConfigurationDirectory, util.DefaultEnvironmentsSubdirectory)
			envDirCurrent := path.Join(util.UsedEnvironmentDirectory, envIDCurrent)
			envDirSwitch := path.Join(util.UsedEnvironmentDirectory, envIDSwitch)
			os.MkdirAll(envDirCurrent, 0755)
			os.Mkdir(envDirSwitch, 0755)
			ioutil.WriteFile(path.Join(envDirCurrent, util.DefaultEnvironmentConfigFileName), []byte(fmt.Sprintf(envConfigTemplate, "env1", envIDCurrent)), 0644)
			ioutil.WriteFile(path.Join(envDirSwitch, util.DefaultEnvironmentConfigFileName), []byte(fmt.Sprintf(envConfigTemplate, "env2", envIDSwitch)), 0644)

			util.UsedConfigFile = tt.configPath
			defer ioutil.WriteFile(tt.configPath, []byte(""), 0644)
			c := &Config{
				Version:            tt.fields.Version,
				Kind:               tt.fields.Kind,
				CurrentEnvironment: tt.fields.CurrentEnvironment,
			}
			err := c.SetUsedEnvironment(tt.uuid)

			if isWrongResult(t, err, tt.wantErr) {
				return
			}

			buf, err := ioutil.ReadFile(tt.configPath)
			if bytes.Compare(buf, tt.want) != 0 {
				t.Errorf("wanted %s but got %s", tt.want, buf)
			}
		})
	}
}

func TestConfig_CreateNewEnvironment(t *testing.T) {
	util.UsedConfigFile, util.UsedConfigurationDirectory = setup(t, "create")
	defer os.RemoveAll(util.UsedConfigurationDirectory)

	type args struct {
		name string
	}
	tests := []struct {
		name    string
		args    args
		mocked  []byte
		wantErr error
	}{
		{
			name: "from nil",
			args: args{
				name: "e1",
			},
			mocked: []byte(`version: v1
kind: Config
current-environment: 00000000-0000-0000-0000-000000000000`),
			wantErr: nil,
		},
		{
			name: "from not nil",
			args: args{
				name: "e1",
			},
			mocked: []byte(`version: v1
kind: Config
current-environment: b3d7be89-461e-41eb-b130-0b4db1555d85`),
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if len(tt.mocked) > 0 {
				_ = ioutil.WriteFile(util.UsedConfigFile, tt.mocked, 0644)
			}
			defer ioutil.WriteFile(util.UsedConfigFile, []byte(""), 0664)
			c, err := GetConfig()
			if err != nil {
				t.Errorf("error getting configuration %v", err)
				return
			}
			err = c.CreateNewEnvironment(tt.args.name)
			if isWrongResult(t, err, tt.wantErr) {
				return
			}
			//TODO check for directory location after fixing global variables creation
		})
	}
}

func TestConfig_Save(t *testing.T) {
	util.UsedConfigFile, util.UsedConfigurationDirectory = setup(t, "save")
	defer os.RemoveAll(util.UsedConfigurationDirectory)

	tests := []struct {
		name    string
		fields  Config
		wantErr bool
	}{
		{
			name:    "empty config",
			wantErr: false,
		},
		{
			name: "minimal config",
			fields: Config{
				Version:            "v1",
				Kind:               KindConfig,
				CurrentEnvironment: uuid.Nil,
			},
			wantErr: false,
		},
		{
			name: "new uuid",
			fields: Config{
				Version:            "v1",
				Kind:               KindConfig,
				CurrentEnvironment: uuid.New(),
			},
			wantErr: false,
		},
		{
			name: "existing uuid",
			fields: Config{
				Version:            "v1",
				Kind:               KindConfig,
				CurrentEnvironment: uuid.MustParse("654e92b3-f06c-43c8-b152-6f2c5557f8af"),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.fields
			if err := c.Save(); (err != nil) != tt.wantErr {
				t.Errorf("Save() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConfig_AddAzureCredentials(t *testing.T) {
	type fields struct {
		CurrentEnvironment uuid.UUID
	}
	type args struct {
		credentials az.Credentials
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Config
	}{
		{
			name: "happy path",
			fields: fields{
				CurrentEnvironment: uuid.MustParse("3e5b7269-1b3d-4003-9454-9f472857633a"),
			},
			args: args{
				credentials: az.Credentials{
					AppID:          "app-id-1",
					Password:       "some-strong-pass",
					Tenant:         "some-tenant-id",
					SubscriptionID: "some-subscription-id",
				},
			},
			want: &Config{
				Version:            "v1",
				Kind:               "Config",
				CurrentEnvironment: uuid.MustParse("3e5b7269-1b3d-4003-9454-9f472857633a"),
				AzureConfig: AzureConfig{
					Credentials: az.Credentials{
						AppID:          "app-id-1",
						Password:       "some-strong-pass",
						Tenant:         "some-tenant-id",
						SubscriptionID: "some-subscription-id",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				Version:            "v1",
				Kind:               KindConfig,
				CurrentEnvironment: tt.fields.CurrentEnvironment,
			}
			c.AddAzureCredentials(tt.args.credentials)
			if !reflect.DeepEqual(c, tt.want) {
				t.Errorf("got = %v, want %v", c, tt.want)
			}
		})
	}
}

func Test_setUsedConfigPaths(t *testing.T) {
	tempFile, tempDir := setup(t, "set")
	defer os.RemoveAll(tempDir)

	tests := []struct {
		name       string
		configDir  string
		configFile string
		mocked     []byte
		want       *Config
		wantErr    error
	}{
		{
			name:       "empty file",
			configDir:  tempDir,
			configFile: tempFile,
			wantErr:    errors.New("EOF"),
		},
		{
			name:       "correct",
			configDir:  tempDir,
			configFile: tempFile,
			mocked: []byte(`version: v1
kind: Config
current-environment: 3e5b7269-1b3d-4003-9454-9f472857633a`),
			want: &Config{
				Version:            "v1",
				Kind:               KindConfig,
				CurrentEnvironment: uuid.MustParse("3e5b7269-1b3d-4003-9454-9f472857633a"),
			},
			wantErr: nil,
		},
		{
			name:       "incorrect",
			configDir:  tempDir,
			configFile: tempFile,
			mocked:     []byte("incorrect file"),
			wantErr:    errors.New("yaml: unmarshal errors:\n  line 1: cannot unmarshal !!str `incorre...` into configuration.Config"),
		},
		{
			name:       "correct with null",
			configDir:  tempDir,
			configFile: tempFile,
			mocked: []byte(`version: v1
kind: Config
current-environment: 00000000-0000-0000-0000-000000000000`),
			want: &Config{
				Version:            "v1",
				Kind:               KindConfig,
				CurrentEnvironment: uuid.MustParse("00000000-0000-0000-0000-000000000000"),
			},
			wantErr: nil,
		},
		{
			name:       "not existing directory",
			configDir:  path.Join(tempDir, "non-existing-config-directory"),
			configFile: path.Join(tempDir, "non-existing-config-directory", "non-existing-file.yaml"),
			wantErr:    nil,
			want: &Config{
				Version:            "v1",
				Kind:               KindConfig,
				CurrentEnvironment: uuid.MustParse("00000000-0000-0000-0000-000000000000"),
			},
		},
		{
			name:       "not existing file",
			configDir:  tempDir,
			configFile: path.Join(tempDir, "another-not-existing-file.yaml"),
			wantErr:    nil,
			want: &Config{
				Version:            "v1",
				Kind:               KindConfig,
				CurrentEnvironment: uuid.MustParse("00000000-0000-0000-0000-000000000000"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			util.UsedConfigFile = ""
			util.UsedConfigurationDirectory = ""
			util.UsedEnvironmentDirectory = ""
			util.UsedRepositoryFile = ""
			util.UsedTempDirectory = ""
			if len(tt.mocked) > 0 {
				_ = ioutil.WriteFile(tt.configFile, tt.mocked, 0644)
			}
			defer ioutil.WriteFile(tt.configFile, []byte(""), 0664)
			got, err := setUsedConfigPaths(tt.configDir, tt.configFile)

			if isWrongResult(t, err, tt.wantErr) {
				t.Errorf("got err = %v, want err %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_makeOrGetConfig(t *testing.T) {
	tempFile, tempDirectory := setup(t, "make")
	defer os.RemoveAll(tempDirectory)

	tests := []struct {
		name       string
		configPath string
		mocked     []byte
		want       *Config
		wantErr    error
	}{
		{
			name:       "empty",
			configPath: tempFile,
			wantErr:    errors.New("EOF"),
		},
		{
			name:       "correct",
			configPath: tempFile,
			mocked: []byte(`version: v1
kind: Config
current-environment: 3e5b7269-1b3d-4003-9454-9f472857633a`),
			want: &Config{
				Version:            "v1",
				Kind:               KindConfig,
				CurrentEnvironment: uuid.MustParse("3e5b7269-1b3d-4003-9454-9f472857633a"),
			},
			wantErr: nil,
		},
		{
			name:       "incorrect",
			configPath: tempFile,
			mocked:     []byte("incorrect file"),
			wantErr:    errors.New("yaml: unmarshal errors:\n  line 1: cannot unmarshal !!str `incorre...` into configuration.Config"),
		},
		{
			name:       "correct with null",
			configPath: tempFile,
			mocked: []byte(`version: v1
kind: Config
current-environment: 00000000-0000-0000-0000-000000000000`),
			want: &Config{
				Version:            "v1",
				Kind:               KindConfig,
				CurrentEnvironment: uuid.MustParse("00000000-0000-0000-0000-000000000000"),
			},
			wantErr: nil,
		},
		{
			name:       "not existing",
			configPath: path.Join(tempDirectory, "non-existing-config-directory"),
			wantErr:    nil,
			want: &Config{
				Version:            "v1",
				Kind:               KindConfig,
				CurrentEnvironment: uuid.MustParse("00000000-0000-0000-0000-000000000000"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			util.UsedConfigFile = tt.configPath
			if len(tt.mocked) > 0 {
				_ = ioutil.WriteFile(tt.configPath, tt.mocked, 0644)
			}
			defer ioutil.WriteFile(tt.configPath, []byte(""), 0664)
			got, err := makeOrGetConfig()
			if isWrongResult(t, err, tt.wantErr) {
				t.Errorf("got err = %v, want err %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got = %v, want %v", got, tt.want)
			}
		})
	}
}

func isWrongResult(t *testing.T, err error, wantErr error) bool {
	if err != nil && wantErr != nil {
		if wantErr.Error() != err.Error() {
			t.Errorf("got error %v, want error %v", err, wantErr)
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

func assertPanic(t *testing.T, f func()) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()
	f()
}
