package repository

import "testing"

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
			if got := inferRepoName(tt.repoName); got != tt.want {
				t.Errorf("inferRepoName() = %v, want %v", got, tt.want)
			}
		})
	}
}
