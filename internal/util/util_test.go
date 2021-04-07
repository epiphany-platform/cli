package util

import (
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/rs/zerolog"
)

func setup() {
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
}
func TestEnsureDirectory(t *testing.T) {
	setup()
	tests := []struct {
		name string
		want string
	}{
		{
			name: "one level of directories",
			want: "test1",
		},
		{
			name: "two levels of directories",
			want: "test2l1/test2l2",
		},
	}
	parentDir := os.TempDir()
	mainDirectory, err := ioutil.TempDir(parentDir, "*-e-util")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(mainDirectory)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expectedDirectory := path.Join(mainDirectory, tt.want)
			EnsureDirectory(expectedDirectory)
			if _, err := os.Stat(expectedDirectory); os.IsNotExist(err) {
				t.Errorf("expected directory not found: %s", tt.want)
			}
		})
	}
}
