package janitor

import (
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/epiphany-platform/cli/internal/util"
	"github.com/epiphany-platform/cli/pkg/configuration"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func setup(a *assert.Assertions) string {
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	parentDir := os.TempDir()
	mainDirectory, err := ioutil.TempDir(parentDir, "*-janitor")
	a.NoError(err)
	return mainDirectory
}

func Test_setUsedConfigPaths(t *testing.T) {
	confDir := setup(assert.New(t))
	defer os.RemoveAll(confDir)

	tests := []struct {
		name      string
		configDir string
		wantDirs  []string
	}{
		{
			name:      "happy path",
			configDir: path.Join(confDir, "hp"),
			wantDirs:  []string{"hp", "hp/environments", "hp/tmp", "hp/repos"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := assert.New(t)
			util.UsedConfigFile = ""
			util.UsedConfigurationDirectory = ""
			util.UsedEnvironmentDirectory = ""
			util.UsedRepositoryFile = ""
			util.UsedTempDirectory = ""
			util.UsedReposDirectory = ""

			setUsedConfigPaths(tt.configDir)

			for _, d := range tt.wantDirs {
				dir := path.Join(confDir, d)
				a.DirExists(dir)
			}
		})
	}
}

func Test_ensureConfig(t *testing.T) {
	confDir := setup(assert.New(t))
	defer os.RemoveAll(confDir)
	util.UsedConfigFile = ""
	util.UsedConfigurationDirectory = ""
	util.UsedEnvironmentDirectory = ""
	util.UsedRepositoryFile = ""
	util.UsedTempDirectory = ""
	util.UsedReposDirectory = ""
	setUsedConfigPaths(confDir)

	tests := []struct {
		name   string
		exists bool
		want   string
	}{
		{
			name:   "does not exist",
			exists: false,
			want:   path.Join(confDir, "config.yaml"),
		},
		{
			name:   "exists",
			exists: true,
			want:   path.Join(confDir, "config.yaml"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := assert.New(t)
			if tt.exists {
				_ = ioutil.WriteFile(tt.want, []byte("random content"), 0644)
			}
			defer os.Remove(tt.want)
			err := ensureConfig()
			a.NoError(err)
			a.FileExists(tt.want)
		})
	}
}

func Test_ensureEnvironment(t *testing.T) {
	confDir := setup(assert.New(t))
	defer os.RemoveAll(confDir)
	util.UsedConfigFile = ""
	util.UsedConfigurationDirectory = ""
	util.UsedEnvironmentDirectory = ""
	util.UsedRepositoryFile = ""
	util.UsedTempDirectory = ""
	util.UsedReposDirectory = ""
	setUsedConfigPaths(confDir)

	tests := []struct {
		name   string
		mocked []byte
		want   string
	}{
		{
			name: "correct",
			mocked: []byte(`version: v1
kind: Config
current-environment: 3e5b7269-1b3d-4003-9454-9f472857633a`),
			want: "3e5b7269-1b3d-4003-9454-9f472857633a",
		},
		{
			name: "correct with null",
			mocked: []byte(`version: v1
kind: Config
current-environment: 00000000-0000-0000-0000-000000000000`),
			want: "new",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := assert.New(t)
			_ = ioutil.WriteFile(util.UsedConfigFile, tt.mocked, 0644)
			defer os.Remove(util.UsedConfigFile)
			err := ensureConfig()
			a.NoError(err)
			err = ensureEnvironment()
			a.NoError(err)
			c, err := configuration.GetConfig()
			a.NoError(err)
			if tt.want != "new" {
				a.Equal(uuid.MustParse(tt.want), c.CurrentEnvironment)
			} else {
				a.NotEqual(uuid.Nil, c.CurrentEnvironment)
			}
		})
	}
}
