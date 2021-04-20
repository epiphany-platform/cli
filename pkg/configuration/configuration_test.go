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

	"github.com/epiphany-platform/cli/internal/util"
	"github.com/epiphany-platform/cli/pkg/az"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func setup(t *testing.T, suffix string) (string, string, string) {
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	parentDir := os.TempDir()
	mainDirectory, err := ioutil.TempDir(parentDir, fmt.Sprintf("*-e-configuration-%s", suffix))
	if err != nil {
		t.Fatal(err)
	}
	envDirectory, err := ioutil.TempDir(mainDirectory, fmt.Sprintf("environments-%s", suffix))
	if err != nil {
		t.Fatal(err)
	}

	tempFile, err := ioutil.TempFile(mainDirectory, "config-*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	return tempFile.Name(), mainDirectory, envDirectory
}

// Create necessary files and directories as environments are validated before switching
func prepareSetUsedEnvironmentResources(t *testing.T, envDirectory, envIDCurrent, envIDSwitch string) {
	envConfigTemplate := `name: %s
uuid: %s
`
	envDirCurrent := path.Join(envDirectory, envIDCurrent)
	envDirSwitch := path.Join(envDirectory, envIDSwitch)
	err := os.MkdirAll(envDirCurrent, 0755)
	if err != nil {
		t.Fatal(err)
	}
	err = os.Mkdir(envDirSwitch, 0755)
	if err != nil {
		t.Fatal(err)
	}
	err = ioutil.WriteFile(path.Join(envDirCurrent, util.DefaultEnvironmentConfigFileName), []byte(fmt.Sprintf(envConfigTemplate, "env1", envIDCurrent)), 0644)
	if err != nil {
		t.Fatal(err)
	}
	err = ioutil.WriteFile(path.Join(envDirSwitch, util.DefaultEnvironmentConfigFileName), []byte(fmt.Sprintf(envConfigTemplate, "env2", envIDSwitch)), 0644)
	if err != nil {
		t.Fatal(err)
	}
}

func TestConfig_SetUsedEnvironment(t *testing.T) {
	util.UsedConfigFile, util.UsedConfigurationDirectory, util.UsedEnvironmentDirectory = setup(t, "used")
	envIDCurrent := "3e5b7269-1b3d-4003-9454-9f472857633a"
	envIDSwitch := "567c0831-7e83-4b56-a2a7-ec7a8327238f"
	prepareSetUsedEnvironmentResources(t, util.UsedEnvironmentDirectory, envIDCurrent, envIDSwitch)
	defer func() {
		_ = os.RemoveAll(util.UsedConfigurationDirectory)
	}()

	type fields struct {
		Version            string
		Kind               Kind
		CurrentEnvironment uuid.UUID
	}

	tests := []struct {
		name    string
		fields  fields
		uuid    uuid.UUID
		wantErr error
		want    []byte
	}{
		{
			name: "nil to some",
			fields: fields{
				Version:            "v1",
				Kind:               KindConfig,
				CurrentEnvironment: uuid.Nil,
			},
			uuid:    uuid.MustParse(envIDSwitch),
			wantErr: nil,
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
			uuid:    uuid.MustParse(envIDSwitch),
			wantErr: nil,
			want: []byte(fmt.Sprintf(`version: v1
kind: Config
current-environment: %s
`, envIDSwitch)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := assert.New(t)
			defer func() {
				_ = ioutil.WriteFile(util.UsedConfigFile, []byte(""), 0644)
			}()
			c := &Config{
				Version:            tt.fields.Version,
				Kind:               tt.fields.Kind,
				CurrentEnvironment: tt.fields.CurrentEnvironment,
			}
			err := c.SetUsedEnvironment(tt.uuid)

			if tt.wantErr == nil {
				a.NoError(err)
			} else {
				a.EqualError(err, tt.wantErr.Error())
			}

			buf, err := ioutil.ReadFile(util.UsedConfigFile)
			a.Equalf(bytes.Compare(buf, tt.want), 0, "wanted %s but got %s", tt.want, buf)
		})
	}
}

func TestConfig_CreateNewEnvironment(t *testing.T) {
	util.UsedConfigFile, util.UsedConfigurationDirectory, util.UsedEnvironmentDirectory = setup(t, "create")
	defer func() {
		_ = os.RemoveAll(util.UsedConfigurationDirectory)
	}()

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
			a := assert.New(t)
			if len(tt.mocked) > 0 {
				_ = ioutil.WriteFile(util.UsedConfigFile, tt.mocked, 0644)
			}
			defer func() {
				_ = ioutil.WriteFile(util.UsedConfigFile, []byte(""), 0644)
			}()
			c, err := GetConfig()
			a.NoErrorf(err, "error getting configuration %v", err)
			u, err := c.CreateNewEnvironment(tt.args.name)
			if tt.wantErr == nil {
				a.NoError(err)
			} else {
				a.EqualError(err, tt.wantErr.Error())
			}
			envDir := path.Join(util.UsedEnvironmentDirectory, u.String())
			a.DirExists(envDir)
			a.FileExists(path.Join(envDir, util.DefaultEnvironmentConfigFileName))
		})
	}
}

func TestConfig_Save(t *testing.T) {
	util.UsedConfigFile, util.UsedConfigurationDirectory, util.UsedEnvironmentDirectory = setup(t, "save")
	defer func() {
		_ = os.RemoveAll(util.UsedConfigurationDirectory)
	}()

	tests := []struct {
		name    string
		fields  Config
		wantErr error
	}{
		{
			name:    "Empty config",
			wantErr: nil,
		},
		{
			name: "Minimal config",
			fields: Config{
				Version:            "v1",
				Kind:               KindConfig,
				CurrentEnvironment: uuid.Nil,
			},
			wantErr: nil,
		},
		{
			name: "New uuid",
			fields: Config{
				Version:            "v1",
				Kind:               KindConfig,
				CurrentEnvironment: uuid.New(),
			},
			wantErr: nil,
		},
		{
			name: "Existing uuid",
			fields: Config{
				Version:            "v1",
				Kind:               KindConfig,
				CurrentEnvironment: uuid.MustParse("654e92b3-f06c-43c8-b152-6f2c5557f8af"),
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := assert.New(t)
			c := tt.fields
			err := c.Save()
			if tt.wantErr == nil {
				a.NoError(err)
			} else {
				a.EqualError(err, tt.wantErr.Error())
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
			a := assert.New(t)
			c := &Config{
				Version:            "v1",
				Kind:               KindConfig,
				CurrentEnvironment: tt.fields.CurrentEnvironment,
			}
			c.AddAzureCredentials(tt.args.credentials)
			a.Truef(reflect.DeepEqual(c, tt.want), "got = %#v, want = %#v", c, tt.want)
		})
	}
}

func TestGetConfig(t *testing.T) {
	var tempFile, tempDirectory string
	tempFile, tempDirectory, _ = setup(t, "get")
	defer func() {
		_ = os.RemoveAll(tempDirectory)
	}()

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
			wantErr:    errors.New("open " + path.Join(tempDirectory, "non-existing-config-directory") + ": no such file or directory"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			util.UsedConfigurationDirectory = tempDirectory
			util.UsedConfigFile = tt.configPath

			a := assert.New(t)
			util.UsedConfigFile = tt.configPath
			if len(tt.mocked) > 0 {
				_ = ioutil.WriteFile(tt.configPath, tt.mocked, 0644)
			}
			defer func() {
				_ = ioutil.WriteFile(tt.configPath, []byte(""), 0644)
			}()
			got, err := GetConfig()
			if tt.wantErr == nil {
				a.NoError(err)
			} else {
				a.EqualError(err, tt.wantErr.Error())
			}
			a.Truef(reflect.DeepEqual(got, tt.want), "got = %#v, want = %#v", got, tt.want)
		})
	}

}
