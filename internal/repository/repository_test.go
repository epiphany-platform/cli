package repository

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"testing"

	"github.com/epiphany-platform/cli/internal/util"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func setup(a *assert.Assertions) (string, string) {
	zerolog.SetGlobalLevel(zerolog.TraceLevel)
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
			util.UsedConfigurationDirectory, util.UsedReposDirectory = setup(a)
			util.UsedConfigFile = ""
			util.UsedEnvironmentDirectory = ""
			util.UsedRepositoryFile = ""
			util.UsedTempDirectory = ""
			defer os.RemoveAll(util.UsedConfigurationDirectory)

			if tt.mocked != nil {
				err := ioutil.WriteFile(path.Join(util.UsedReposDirectory, tt.args.inferredRepoName+".yaml"), tt.mocked, 0644)
				a.NoError(err)
			}
			err := persistV1RepositoryFile(tt.args.inferredRepoName, tt.args.v1, tt.args.force)
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

func Test_downloadV1Repository(t *testing.T) {
	util.UsedConfigurationDirectory = ""
	util.UsedReposDirectory = ""
	util.UsedConfigFile = ""
	util.UsedEnvironmentDirectory = ""
	util.UsedRepositoryFile = ""
	util.UsedTempDirectory = ""

	type args struct {
		url string
	}
	tests := []struct {
		name    string
		args    args
		want    *V1
		wantErr bool
	}{
		{
			name: "happy path",
			args: args{
				url: fmt.Sprintf("/%s/%s/%s", "test-user/test-repo", util.DefaultRepositoryBranch, util.DefaultV1RepositoryFileName),
			},
			want: &V1{
				Version:    "v",
				Kind:       "k",
				Name:       "n",
				Components: []Component{},
			},
			wantErr: false,
		},
		{
			name: "incorrect url",
			args: args{
				url: "/test-user/test-repo/incorrect/url.yaml",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := assert.New(t)

			server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				t.Logf("got test request to: %s", req.URL.String())
				if req.URL.String() == "/test-user/test-repo/HEAD/v1.yaml" {
					_, _ = rw.Write([]byte(`version: v
kind: k
name: n
components: []
`))
				}
			}))
			defer server.Close()

			httpClient = server.Client()
			got, err := downloadV1Repository(server.URL + tt.args.url)
			if tt.wantErr {
				a.Error(err)
			} else {
				a.NoError(err)
			}
			a.Equal(tt.want, got)
		})
	}
}

func Test_decodeV1Repository(t *testing.T) {
	type args struct {
		fileName string
	}
	tests := []struct {
		name    string
		mocked  []byte
		args    args
		want    *V1
		wantErr bool
	}{
		{
			name: "happy path",
			args: args{fileName: "repo-file.yaml"},
			mocked: []byte(`version: v1
kind: Repository
name: c
components: []
`),
			want: &V1{
				Version:    "v1",
				Kind:       "Repository",
				Name:       "c",
				Components: []Component{},
			},
			wantErr: false,
		},
		{
			name:    "missing repo file",
			args:    args{fileName: "not-existing-repo-file.yaml"},
			mocked:  nil,
			want:    nil,
			wantErr: true,
		},
		{
			name:    "empty file",
			args:    args{fileName: "empty-repo-file.yaml"},
			mocked:  []byte(``),
			want:    nil,
			wantErr: true,
		},
		{
			name:    "not a yaml",
			args:    args{fileName: "not-a-yaml-repo-file.yaml"},
			mocked:  []byte(`not a yaml`),
			want:    nil,
			wantErr: true,
		},
		{
			name:    "missing version field",
			args:    args{fileName: "missing-version-field-repo-file.yaml"},
			mocked:  []byte(`kind: a`),
			want:    nil,
			wantErr: true,
		},
		{
			name:    "missing kind field",
			args:    args{fileName: "missing-kind-field-repo-file.yaml"},
			mocked:  []byte(`version: a`),
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := assert.New(t)
			util.UsedConfigurationDirectory, util.UsedReposDirectory = setup(a)
			util.UsedConfigFile = ""
			util.UsedEnvironmentDirectory = ""
			util.UsedRepositoryFile = ""
			util.UsedTempDirectory = ""
			defer os.RemoveAll(util.UsedConfigurationDirectory)

			if tt.mocked != nil {
				err := ioutil.WriteFile(path.Join(util.UsedReposDirectory, tt.args.fileName), tt.mocked, 0644)
				a.NoError(err)
			}
			got, err := decodeV1Repository(path.Join(util.UsedReposDirectory, tt.args.fileName))
			if tt.wantErr {
				a.Error(err)
			} else {
				a.NoError(err)
				a.Equal(tt.want, got)
			}
		})
	}
}

func Test_load(t *testing.T) {
	tests := []struct {
		name    string
		mocked  map[string][]byte
		wantErr bool
		wantLen int
	}{
		{
			name: "happy path",
			mocked: map[string][]byte{
				"first-repo.yaml": []byte(`version: v1
kind: Repository
name: first
components: []
`),
				"second-repo.yaml": []byte(`version: v1
kind: Repository
name: second
components: []
`),
			},
			wantErr: false,
			wantLen: 2,
		},
		{
			name: "incorrect extension",
			mocked: map[string][]byte{
				"first-repo.incorrect": []byte(`version: v1
kind: Repository
name: first
components: []
`),
				"correct-repo.yaml": []byte(`version: v1
kind: Repository
name: correct
components: []
`),
			},
			wantErr: false,
			wantLen: 1,
		},
		{
			name: "incorrect content",
			mocked: map[string][]byte{
				"incorrect-repo.yaml": []byte(`incorrect`),
				"correct-repo.yaml": []byte(`version: v1
kind: Repository
name: correct
components: []
`),
			},
			wantErr: false,
			wantLen: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := assert.New(t)
			util.UsedConfigurationDirectory, util.UsedReposDirectory = setup(a)
			util.UsedConfigFile = ""
			util.UsedEnvironmentDirectory = ""
			util.UsedRepositoryFile = ""
			util.UsedTempDirectory = ""
			defer os.RemoveAll(util.UsedConfigurationDirectory)

			if tt.mocked != nil {
				for k, v := range tt.mocked {
					err := ioutil.WriteFile(path.Join(util.UsedReposDirectory, k), v, 0644)
					a.NoError(err)
				}
			}
			err := load()
			if tt.wantErr {
				a.Error(err)
			} else {
				a.NoError(err)
				a.Len(loaded.v1s, tt.wantLen)
			}
		})
	}
}
