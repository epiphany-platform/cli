package auth

import (
	"io/ioutil"
	"reflect"
	"testing"
)

func TestGenerateRsaKeyPair(t *testing.T) {
	type args struct {
		directory string
	}
	tests := []struct {
		name    string
		args    args
		wantKp  RsaKeyPair
		wantErr bool
	}{
		{
			name: "happy path",
			args: args{
				directory: "hp",
			},
			wantKp: RsaKeyPair{
				Name: rsaKeyName,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {

		dir, err := ioutil.TempDir("", tt.args.directory)
		if err != nil {
			t.Error(err)
		}

		t.Run(tt.name, func(t *testing.T) {
			gotKp, err := GenerateRsaKeyPair(dir)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateRsaKeyPair() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotKp, tt.wantKp) {
				t.Errorf("GenerateRsaKeyPair() gotKp = %v, want %v", gotKp, tt.wantKp)
			}
		})
	}
}

func Test_generateRsaKeyPair(t *testing.T) {
	type args struct {
		directory string
		name      string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "happy path",
			args: args{
				directory: "hp",
				name:      rsaKeyName,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {

		dir, err := ioutil.TempDir("", tt.args.directory)
		if err != nil {
			t.Error(err)
		}
		t.Run(tt.name, func(t *testing.T) {
			if err := generateRsaKeyPair(dir, tt.args.name); (err != nil) != tt.wantErr {
				t.Errorf("generateRsaKeyPair() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
