/*
 * Copyright © 2020 Mateusz Kyc
 */

package repository

import (
	"errors"
	"fmt"
	"github.com/mkyc/epiphany-wrapper-poc/pkg/util"
	"github.com/rs/zerolog"
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"regexp"
	"testing"
)

func setup(t *testing.T, suffix string) (string, string, string, string) {
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	parentDir := os.TempDir()
	configDirectory, err := ioutil.TempDir(parentDir, fmt.Sprintf("*-e-repository-%s", suffix))
	if err != nil {
		t.Fatal(err)
	}
	envsDirectory, err := ioutil.TempDir(configDirectory, "environments-*")
	if err != nil {
		t.Fatal(err)
	}

	configFile, err := ioutil.TempFile(configDirectory, "config-*.yaml")
	if err != nil {
		t.Fatal(err)
	}

	repoFile, err := ioutil.TempFile(configDirectory, "v1-*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	return configFile.Name(), configDirectory, envsDirectory, repoFile.Name()
}

func TestComponent_JustLatestVersion(t *testing.T) {
	tests := []struct {
		name    string
		mock    *Component
		want    *Component
		wantErr error
	}{
		{
			name: "single correct",
			mock: &Component{
				Name: "c",
				Type: "t",
				Versions: []ComponentVersion{
					{
						Version:  "v1",
						IsLatest: true,
					},
				},
			},
			want: &Component{
				Name: "c",
				Type: "t",
				Versions: []ComponentVersion{
					{
						Version:  "v1",
						IsLatest: true,
					},
				},
			},
			wantErr: nil,
		},
		{
			name: "no versions",
			mock: &Component{
				Name:     "c",
				Type:     "t",
				Versions: []ComponentVersion{},
			},
			wantErr: errors.New("no versions in component"),
		},
		{
			name: "one not latest",
			mock: &Component{
				Name: "c",
				Type: "t",
				Versions: []ComponentVersion{
					{
						Version:  "v1",
						IsLatest: false,
					},
				},
			},
			wantErr: errors.New("component only version is not marked latest"),
		},
		{
			name: "multiple latest",
			mock: &Component{
				Name: "c",
				Type: "t",
				Versions: []ComponentVersion{
					{
						Version:  "v1",
						IsLatest: true,
					},
					{
						Version:  "v2",
						IsLatest: true,
					},
				},
			},
			wantErr: errors.New("incorrect number of latest versions"),
		},
		{
			name: "multiple correct",
			mock: &Component{
				Name: "c",
				Type: "t",
				Versions: []ComponentVersion{
					{
						Version:  "v1",
						IsLatest: false,
					},
					{
						Version:  "v2",
						IsLatest: false,
					},
					{
						Version:  "v3",
						IsLatest: true,
					},
				},
			},
			want: &Component{
				Name: "c",
				Type: "t",
				Versions: []ComponentVersion{
					{
						Version:  "v3",
						IsLatest: true,
					},
				},
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Component{
				Name:     tt.mock.Name,
				Type:     tt.mock.Type,
				Versions: tt.mock.Versions,
			}
			got, err := c.JustLatestVersion()
			if isWrongResult(t, err, tt.wantErr) {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestV1_GetComponentByName(t *testing.T) {
	tests := []struct {
		name          string
		mock          *V1
		componentName string
		want          *Component
		wantErr       error
	}{
		{
			name: "single correct",
			mock: &V1{
				Version: "v1",
				Kind:    "k1",
				Components: []Component{
					{
						Name: "c1",
						Type: "t1",
						Versions: []ComponentVersion{
							{
								Version:  "v1",
								IsLatest: true,
								Image:    "i1",
							},
						},
					},
				},
			},
			componentName: "c1",
			want: &Component{
				Name: "c1",
				Type: "t1",
				Versions: []ComponentVersion{
					{
						Version:  "v1",
						IsLatest: true,
						Image:    "i1",
					},
				},
			},
			wantErr: nil,
		},
		{
			name: "multiple correct",
			mock: &V1{
				Version: "v1",
				Kind:    "k1",
				Components: []Component{
					{
						Name: "c1",
						Type: "t1",
						Versions: []ComponentVersion{
							{
								Version:  "v1",
								IsLatest: true,
								Image:    "i1",
							},
						},
					},
					{
						Name: "c1",
						Type: "t2",
						Versions: []ComponentVersion{
							{
								Version:  "v2",
								IsLatest: true,
								Image:    "i2",
							},
						},
					},
				},
			},
			componentName: "c1",
			want: &Component{
				Name: "c1",
				Type: "t1",
				Versions: []ComponentVersion{
					{
						Version:  "v1",
						IsLatest: true,
						Image:    "i1",
					},
				},
			},
			wantErr: nil,
		},
		{
			name: "missing",
			mock: &V1{
				Version: "v1",
				Kind:    "k1",
				Components: []Component{
					{
						Name: "c1",
						Type: "t1",
						Versions: []ComponentVersion{
							{
								Version:  "v1",
								IsLatest: true,
								Image:    "i1",
							},
						},
					},
				},
			},
			componentName: "c2",
			wantErr:       errors.New("unknown component"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := V1{
				Version:    tt.mock.Version,
				Kind:       tt.mock.Kind,
				Components: tt.mock.Components,
			}
			got, err := v.GetComponentByName(tt.componentName)
			if isWrongResult(t, err, tt.wantErr) {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got = %#v, want %#v", got, tt.want)
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

func Test_loadRepository(t *testing.T) {
	var repoFile string
	util.UsedConfigFile, util.UsedConfigurationDirectory, util.UsedEnvironmentDirectory, repoFile = setup(t, "load-repository")
	defer os.RemoveAll(util.UsedConfigurationDirectory)

	tests := []struct {
		name         string
		repoFilePath string
		mocked       []byte
		want         *V1
		wantErr      error
	}{
		{
			name:         "correct",
			repoFilePath: repoFile,
			mocked: []byte(`version: v1
kind: k1
components:
- name: c1
  type: t1
  versions:
  - version: v1
    latest: true
`),
			want: &V1{
				Version: "v1",
				Kind:    "k1",
				Components: []Component{
					{
						Name: "c1",
						Type: "t1",
						Versions: []ComponentVersion{
							{
								Version:  "v1",
								IsLatest: true,
							},
						},
					},
				},
			},
			wantErr: nil,
		},
		{
			name:         "empty file",
			repoFilePath: repoFile,
			mocked:       []byte(""),
			wantErr:      errors.New("EOF"),
		},
		{
			name:         "incorrect file content",
			repoFilePath: repoFile,
			mocked:       []byte(`incorrect content`),
			wantErr:      errors.New("yaml: unmarshal errors:\n  line 1: cannot unmarshal !!str `incorre...` into repository.V1"),
		},
		{
			name:         "net existing file",
			repoFilePath: path.Join(util.UsedConfigurationDirectory, "not-existing-file.yaml"),
			wantErr:      errors.New("open .*-e-repository-load-repository/not-existing-file.yaml: no such file or directory"),
		},
	}
	for _, tt := range tests {
		util.UsedRepositoryFile = tt.repoFilePath
		t.Run(tt.name, func(t *testing.T) {
			if len(tt.mocked) > 0 {
				_ = ioutil.WriteFile(tt.repoFilePath, tt.mocked, 0644)
				defer ioutil.WriteFile(tt.repoFilePath, []byte(""), 0664)
			}
			got, err := loadRepository(util.UsedRepositoryFile)
			if isWrongResult(t, err, tt.wantErr) {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got = %#v, want %#v", got, tt.want)
			}
		})
	}
}
