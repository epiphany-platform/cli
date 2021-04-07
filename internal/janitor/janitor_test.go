package janitor

import (
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/epiphany-platform/cli/internal/util"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func setup(t *testing.T) string {
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	parentDir := os.TempDir()
	mainDirectory, err := ioutil.TempDir(parentDir, "*-janitor")
	if err != nil {
		t.Fatal(err)
	}
	return mainDirectory
}

func Test_setUsedConfigPaths(t *testing.T) {
	confDir := setup(t)
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
