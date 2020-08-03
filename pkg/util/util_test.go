/*
 * Copyright © 2020 Mateusz Kyc
 */

package util

import (
	"os"
	"path"
	"testing"
)

func TestEnsureDirectory(t *testing.T) {
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
	mainDirectory := os.TempDir()
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
