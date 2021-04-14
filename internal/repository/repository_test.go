package repository

import (
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/epiphany-platform/cli/internal/util"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func setup(a *assert.Assertions) (string, string) {
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	parentDir := os.TempDir()
	mainDirectory, err := ioutil.TempDir(parentDir, "*-repository")
	a.NoError(err)
	reposDirectory, err := ioutil.TempDir(mainDirectory, "*-"+util.DefaultRepoDirectoryName)
	a.NoError(err)
	return mainDirectory, reposDirectory
}

func Test_inferRepoName(t *testing.T) {
	tests := []struct {
		name     string
		repoName string
		want     string
	}{
		{
			name:     "happy path 1",
			repoName: "epiphany-platform/modules",
			want:     "epiphany-platform-modules",
		},
		{
			name:     "happy path 2",
			repoName: "mkyc/my-epipany-repo",
			want:     "mkyc-my-epipany-repo",
		},
		{
			name:     "full url",
			repoName: "https://github.com/mkyc/my-epipany-repo",
			want:     "mkyc-my-epipany-repo",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := assert.New(t)
			a.Equal(inferRepoName(tt.repoName), tt.want)
		})
	}
}

func Test_persistV1RepositoryFile(t *testing.T) {
	util.UsedConfigurationDirectory, util.UsedReposDirectory = setup(assert.New(t))
	util.UsedConfigFile = ""
	util.UsedEnvironmentDirectory = ""
	util.UsedRepositoryFile = ""
	util.UsedTempDirectory = ""

	type args struct {
		inferredRepoName string
		v1               *V1
		force            bool
	}
	tests := []struct {
		name    string
		args    args
		mocked  []byte
		wantErr bool
		want    []byte
	}{
		{
			name: "new repo",
			args: args{
				inferredRepoName: "test-repo-name",
				v1: &V1{
					Version: "vvv",
					Kind:    "kkk",
					Name:    "nnn",
				},
				force: false,
			},
			mocked:  nil,
			wantErr: false,
			want: []byte(`version: vvv
kind: kkk
name: nnn
components: []
`),
		},
		{
			name: "existing no force",
			args: args{
				inferredRepoName: "existing-repo-file",
				v1: &V1{
					Version: "vvv",
					Kind:    "kkk",
					Name:    "nnn",
				},
				force: false,
			},
			mocked:  []byte(`random content`),
			wantErr: true,
			want:    nil,
		},
		{
			name: "existing but force",
			args: args{
				inferredRepoName: "existing-repo-file-2",
				v1: &V1{
					Version: "vv",
					Kind:    "kk",
					Name:    "nn",
				},
				force: true,
			},
			mocked:  []byte(`random content`),
			wantErr: false,
			want: []byte(`version: vv
kind: kk
name: nn
components: []
`),
		},
		{
			name: "nil repo",
			args: args{
				inferredRepoName: "nil-repo-name",
				v1:               nil,
				force:            false,
			},
			mocked:  nil,
			wantErr: true,
			want:    nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := assert.New(t)
			if tt.mocked != nil {
				ioutil.WriteFile(path.Join(util.UsedReposDirectory, tt.args.inferredRepoName+".yaml"), tt.mocked, 0644)
			}
			err := persistV1RepositoryFile(tt.args.inferredRepoName, tt.args.v1, tt.args.force)
			defer os.Remove(path.Join(util.UsedReposDirectory, tt.args.inferredRepoName+".yaml"))
			if tt.wantErr {
				a.Error(err)
			} else {
				a.NoError(err)
				a.FileExists(path.Join(util.UsedReposDirectory, tt.args.inferredRepoName+".yaml"))
				got, err2 := ioutil.ReadFile(path.Join(util.UsedReposDirectory, tt.args.inferredRepoName+".yaml"))
				a.NoError(err2)
				a.Equal(string(tt.want), string(got))
			}
		})
	}
}
