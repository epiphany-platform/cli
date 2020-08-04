/*
 * Copyright © 2020 Mateusz Kyc
 */

package environment

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/mkyc/epiphany-wrapper-poc/pkg/util"
	"github.com/rs/zerolog"
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"testing"
)

func setup(t *testing.T, suffix string) (string, string, string) {
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	parentDir := os.TempDir()
	mainDirectory, err := ioutil.TempDir(parentDir, fmt.Sprintf("*-e-environment-%s", suffix))
	if err != nil {
		t.Fatal(err)
	}
	envsDirectory, err := ioutil.TempDir(mainDirectory, "environments-*")
	if err != nil {
		t.Fatal(err)
	}

	tempFile, err := ioutil.TempFile(mainDirectory, "config-*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	return tempFile.Name(), mainDirectory, envsDirectory
}

func Test_create(t *testing.T) {
	util.UsedConfigFile, util.UsedConfigurationDirectory, util.UsedEnvironmentDirectory = setup(t, "create")
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
			got, err := create(tt.args.name, uuid.MustParse(tt.args.uuid))
			if isWrongResult(t, err, tt.wantErr) {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got = %+v, want %+v", got, tt.want)
				return
			}
			expectedConfigFile := path.Join(util.UsedEnvironmentDirectory, got.Uuid.String(), util.DefaultEnvironmentConfigFileName)
			if _, err := os.Stat(expectedConfigFile); os.IsNotExist(err) {
				t.Errorf("expected to find file %s but didn't find", expectedConfigFile)
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
