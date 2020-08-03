/*
 * Copyright © 2020 Mateusz Kyc
 */

package configuration

import (
	"bytes"
	"errors"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"testing"
)

func setup(t *testing.T) (string, string) {
	zerolog.SetGlobalLevel(zerolog.WarnLevel)
	parentDir := os.TempDir()
	mainDirectory, err := ioutil.TempDir(parentDir, "*-e-configuration-save")
	if err != nil {
		t.Fatal(err)
	}

	tempFile, err := ioutil.TempFile(mainDirectory, "config-*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	return tempFile.Name(), mainDirectory
}

func TestConfig_Save(t *testing.T) {
	tempFile, tempDirectory := setup(t)
	defer os.RemoveAll(tempDirectory)

	usedConfigFile = tempFile

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

func TestSetConfig(t *testing.T) {
	tempFile, tempDirectory := setup(t)
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
			configPath: path.Join(tempDirectory, "non-existing-config-file-name.yaml"),
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
			usedConfigFile = ""
			if len(tt.mocked) > 0 {
				_ = ioutil.WriteFile(tt.configPath, tt.mocked, 0644)
			}
			defer ioutil.WriteFile(tt.configPath, []byte(""), 0664)
			got, err := SetConfig(tt.configPath)

			if isWrongResult(t, err, tt.wantErr) {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got = %v, want %v", got, tt.wantErr)
			}
		})
	}
}

func Test_makeOrGetConfig(t *testing.T) {
	tempFile, tempDirectory := setup(t)
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
			configPath: path.Join(tempDirectory, "non-existing-config-file-name.yaml"),
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
			usedConfigFile = tt.configPath
			if len(tt.mocked) > 0 {
				_ = ioutil.WriteFile(tt.configPath, tt.mocked, 0644)
			}
			defer ioutil.WriteFile(tt.configPath, []byte(""), 0664)
			got, err := makeOrGetConfig()
			if isWrongResult(t, err, tt.wantErr) {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_GetConfigFilePath(t *testing.T) {
	tempFile, tempDirectory := setup(t)
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
			usedConfigFile = tt.mocked
			c := &Config{}
			defer func() { recover() }()
			f := func() {
				if got := c.GetConfigFilePath(); got != tt.want {
					t.Errorf("got %s, want %s", got, tt.want)
				}
			}
			if tt.shouldPanic {
				assertPanic(t, f)
			}
			f()
		})
	}
}

func TestConfig_SetUsedEnvironment(t *testing.T) {
	tempFile, tempDirectory := setup(t)
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
			uuid:       uuid.MustParse("3e5b7269-1b3d-4003-9454-9f472857633a"),
			configPath: tempFile,
			wantErr:    nil,
			want: []byte(`version: v1
kind: Config
current-environment: 3e5b7269-1b3d-4003-9454-9f472857633a
`),
		},
		{
			name: "some to another",
			fields: fields{
				Version:            "v1",
				Kind:               KindConfig,
				CurrentEnvironment: uuid.MustParse("3e5b7269-1b3d-4003-9454-9f472857633a"),
			},
			uuid:       uuid.MustParse("567c0831-7e83-4b56-a2a7-ec7a8327238f"),
			configPath: tempFile,
			wantErr:    nil,
			want: []byte(`version: v1
kind: Config
current-environment: 567c0831-7e83-4b56-a2a7-ec7a8327238f
`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			usedConfigFile = tt.configPath
			defer ioutil.WriteFile(tt.configPath, []byte(""), 0664)
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
