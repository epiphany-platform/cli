package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/epiphany-platform/cli/internal/util"

	"github.com/stretchr/testify/assert"
)

func setup(t *testing.T, suffix string) (string, string, string, string, string) {
	parentDir := os.TempDir()
	configDirectory, err := ioutil.TempDir(parentDir, fmt.Sprintf("*-e-repository-%s", suffix))
	if err != nil {
		t.Fatal(err)
	}
	envsDirectory := path.Join(configDirectory, util.DefaultEnvironmentsSubdirectory)
	err = os.Mkdir(envsDirectory, 0755)
	if err != nil {
		t.Fatal(err)
	}
	tempDirectory := path.Join(configDirectory, util.DefaultEnvironmentsTempSubdirectory)
	err = os.Mkdir(tempDirectory, 0755)
	if err != nil {
		t.Fatal(err)
	}

	reposDirectory := path.Join(configDirectory, util.DefaultRepoDirectoryName)
	err = os.Mkdir(reposDirectory, 0755)
	if err != nil {
		t.Fatal(err)
	}

	configFile := path.Join(configDirectory, util.DefaultConfigFileName)

	return configFile, configDirectory, envsDirectory, reposDirectory, tempDirectory
}

func TestMain(m *testing.M) {
	err := os.Chdir("..")
	if err != nil {
		fmt.Printf("could not change dir: %v", err)
		os.Exit(1)
	}
	cmd := exec.Command("make", "build")
	err = cmd.Run()
	if err != nil {
		fmt.Printf("could not make binary for e: %v", err)
		os.Exit(1)
	}

	os.Exit(m.Run())
}

func TestAz(t *testing.T) {
	util.UsedConfigFile, util.UsedConfigurationDirectory, util.UsedEnvironmentDirectory, util.UsedReposDirectory, util.UsedTempDirectory = setup(t, "az")
	defer os.RemoveAll(util.UsedConfigurationDirectory)

	tests := []struct {
		name    string
		args    []string
		want    []string
		wantErr bool
	}{
		{
			name:    "e az",
			args:    []string{"--configDir", util.UsedConfigurationDirectory, "az"},
			want:    []string{"Available Commands:\n  sp"},
			wantErr: false,
		},
		{
			name:    "e az sp",
			args:    []string{"--configDir", util.UsedConfigurationDirectory, "az", "sp"},
			want:    []string{"Available Commands:\n  create"},
			wantErr: false,
		},
		{
			name:    "e az sp create",
			args:    []string{"--configDir", util.UsedConfigurationDirectory, "az", "sp", "create"},
			want:    []string{"no tenantID defined"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := assert.New(t)
			dir, err := os.Getwd()
			a.NoError(err)

			cmd := exec.Command(path.Join(dir, "output", "e"), tt.args...)
			got, err := cmd.CombinedOutput()
			if tt.wantErr {
				a.Error(err)
			} else {
				a.NoError(err)
			}

			for _, w := range tt.want {
				a.Contains(string(got), w)
			}
		})
	}
}

func TestEnvironments(t *testing.T) {
	util.UsedConfigFile, util.UsedConfigurationDirectory, util.UsedEnvironmentDirectory, util.UsedReposDirectory, util.UsedTempDirectory = setup(t, "environments")
	defer os.RemoveAll(util.UsedConfigurationDirectory)

	tests := []struct {
		name    string
		args    []string
		want    []string
		wantErr bool
	}{
		{
			name:    "e environments info",
			args:    []string{"--configDir", util.UsedConfigurationDirectory, "environments", "info"},
			want:    []string{"Environment info:\n Name: " + time.Now().Format("060102")},
			wantErr: false,
		},
		{
			name:    "e environments list",
			args:    []string{"--configDir", util.UsedConfigurationDirectory, "environments", "list"},
			want:    []string{fmt.Sprintf(") | %s-", time.Now().Format("060102"))},
			wantErr: false,
		},
		{
			name:    "e environments new",
			args:    []string{"--configDir", util.UsedConfigurationDirectory, "environments", "new", "t1", "--logLevel", "debug"},
			want:    []string{"Created an environment with id "},
			wantErr: false,
		},
		{
			name:    "e environments use",
			args:    []string{"--configDir", util.UsedConfigurationDirectory, "environments", "use", "something-incorrect"},
			want:    []string{"Error: invalid UUID length: 19"},
			wantErr: false,
		},
		{
			name:    "e environments run no args",
			args:    []string{"--configDir", util.UsedConfigurationDirectory, "environments", "run"},
			want:    []string{"Usage:\n  e environments run"},
			wantErr: false,
		},
		{
			name:    "e environments run one arg",
			args:    []string{"--configDir", util.UsedConfigurationDirectory, "environments", "run", "something"},
			want:    []string{"Usage:\n  e environments run"},
			wantErr: false,
		},
		{
			name:    "e environments run two incorrect args",
			args:    []string{"--configDir", util.UsedConfigurationDirectory, "environments", "run", "something", "incorrect"},
			want:    []string{"no such component installed"},
			wantErr: true,
		},
		{
			name:    "e environments export",
			args:    []string{"--configDir", util.UsedConfigurationDirectory, "environments", "export", "--destination", util.UsedConfigurationDirectory, "--logLevel", "debug"},
			want:    []string{"Export operation finished correctly (environment id:"},
			wantErr: false,
		},
		{
			name:    "e environments export wrong env id",
			args:    []string{"--configDir", util.UsedConfigurationDirectory, "environments", "export", "--id", "fcfd81e4-27a8-4ee6-8bb3-f71b8218ba6d"},
			want:    []string{"Environment not found (environment id:"},
			wantErr: true,
		},
		{
			name:    "e environments export wrong destination",
			args:    []string{"--configDir", util.UsedConfigurationDirectory, "environments", "export", "--destination", "/fake/path"},
			want:    []string{"Unable to export environment (environment id:"},
			wantErr: true,
		},
		{
			name:    "e environments import not existing",
			args:    []string{"--configDir", util.UsedConfigurationDirectory, "environments", "import", "--from", path.Join(util.UsedConfigurationDirectory, "ba03a2ba-8fa0-4c15-ac07-894af3dbb365.zip")},
			want:    []string{"Incorrect file path specified"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := assert.New(t)
			dir, err := os.Getwd()
			a.NoError(err)

			cmd := exec.Command(path.Join(dir, "output", "e"), tt.args...)
			got, err := cmd.CombinedOutput()
			if tt.wantErr {
				a.Error(err)
			} else {
				a.NoError(err)
			}

			for _, w := range tt.want {
				a.Contains(string(got), w)
			}
		})
	}
}

func TestHelp(t *testing.T) {
	util.UsedConfigFile, util.UsedConfigurationDirectory, util.UsedEnvironmentDirectory, util.UsedReposDirectory, util.UsedTempDirectory = setup(t, "help")
	defer os.RemoveAll(util.UsedConfigurationDirectory)

	tests := []struct {
		name            string
		args            []string
		wantSubcommands []string
		wantFlags       []string
	}{
		{
			name:            "e --help",
			args:            []string{"--help"},
			wantSubcommands: []string{"az", "environments", "help", "module", "repos", "ssh"},
			wantFlags:       []string{"configDir", "help", "logLevel"},
		},
		{
			name:            "e az --help",
			args:            []string{"az", "--help"},
			wantSubcommands: []string{"sp"},
			wantFlags:       []string{"configDir", "help", "logLevel"},
		},
		{
			name:            "e az sp --help",
			args:            []string{"az", "sp", "--help"},
			wantSubcommands: []string{"create"},
			wantFlags:       []string{"configDir", "help", "logLevel"},
		},
		{
			name:            "e az sp create --help",
			args:            []string{"az", "sp", "create", "--help"},
			wantSubcommands: []string{},
			wantFlags:       []string{"configDir", "help", "logLevel", "name", "subscriptionID", "tenantID"},
		},
		{
			name:            "e environments --help",
			args:            []string{"environments", "--help"},
			wantSubcommands: []string{"export", "import", "info", "list", "new", "run", "use"},
			wantFlags:       []string{"configDir", "help", "logLevel"},
		},
		{
			name:            "e environments export --help",
			args:            []string{"environments", "export", "--help"},
			wantSubcommands: []string{},
			wantFlags:       []string{"configDir", "help", "logLevel", "destination", "id"},
		},
		{
			name:            "e environments import --help",
			args:            []string{"environments", "import", "--help"},
			wantSubcommands: []string{},
			wantFlags:       []string{"configDir", "help", "logLevel", "from"},
		},
		{
			name:            "e environments info --help",
			args:            []string{"environments", "info", "--help"},
			wantSubcommands: []string{},
			wantFlags:       []string{"configDir", "help", "logLevel"},
		},
		{
			name:            "e environments list --help",
			args:            []string{"environments", "list", "--help"},
			wantSubcommands: []string{},
			wantFlags:       []string{"configDir", "help", "logLevel"},
		},
		{
			name:            "e environments new --help",
			args:            []string{"environments", "new", "--help"},
			wantSubcommands: []string{},
			wantFlags:       []string{"configDir", "help", "logLevel", "name"},
		},
		{
			name:            "e environments run --help",
			args:            []string{"environments", "run", "--help"},
			wantSubcommands: []string{},
			wantFlags:       []string{"configDir", "help", "logLevel"},
		},
		{
			name:            "e environments use --help",
			args:            []string{"environments", "use", "--help"},
			wantSubcommands: []string{},
			wantFlags:       []string{"configDir", "help", "logLevel"},
		},
		{
			name:            "e module --help",
			args:            []string{"module", "--help"},
			wantSubcommands: []string{"info", "install", "search"},
			wantFlags:       []string{"configDir", "help", "logLevel"},
		},
		{
			name:            "e module info --help",
			args:            []string{"module", "info", "--help"},
			wantSubcommands: []string{},
			wantFlags:       []string{"configDir", "help", "logLevel"},
		},
		{
			name:            "e module install --help",
			args:            []string{"module", "install", "--help"},
			wantSubcommands: []string{},
			wantFlags:       []string{"configDir", "help", "logLevel"},
		},
		{
			name:            "e module search --help",
			args:            []string{"module", "search", "--help"},
			wantSubcommands: []string{},
			wantFlags:       []string{"configDir", "help", "logLevel"},
		},
		{
			name:            "e repos --help",
			args:            []string{"repos", "--help"},
			wantSubcommands: []string{"install", "list"},
			wantFlags:       []string{"configDir", "help", "logLevel"},
		},
		{
			name:            "e repos install --help",
			args:            []string{"repos", "install", "--help"},
			wantSubcommands: []string{},
			wantFlags:       []string{"configDir", "help", "logLevel", "branch", "force"},
		},
		{
			name:            "e repos list --help",
			args:            []string{"repos", "list", "--help"},
			wantSubcommands: []string{},
			wantFlags:       []string{"configDir", "help", "logLevel"},
		},
		{
			name:            "e ssh --help",
			args:            []string{"ssh", "--help"},
			wantSubcommands: []string{"keygen"},
			wantFlags:       []string{"configDir", "help", "logLevel"},
		},
		{
			name:            "e ssh keygen --help",
			args:            []string{"ssh", "keygen", "--help"},
			wantSubcommands: []string{"create"},
			wantFlags:       []string{"configDir", "help", "logLevel"},
		},
		{
			name:            "e ssh keygen create --help",
			args:            []string{"ssh", "keygen", "create", "--help"},
			wantSubcommands: []string{},
			wantFlags:       []string{"configDir", "help", "logLevel"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := assert.New(t)
			dir, err := os.Getwd()
			a.NoError(err)

			cmd := exec.Command(path.Join(dir, "output", "e"), tt.args...)
			got, err := cmd.CombinedOutput()
			a.NoError(err)

			subCommands := extractSubcommandsNames(string(got))
			a.ElementsMatch(tt.wantSubcommands, subCommands)
			flags := extractFlagsNames(string(got))
			a.ElementsMatch(tt.wantFlags, flags)
		})
	}
}

func TestModule(t *testing.T) {
	util.UsedConfigFile, util.UsedConfigurationDirectory, util.UsedEnvironmentDirectory, util.UsedReposDirectory, util.UsedTempDirectory = setup(t, "module")
	defer os.RemoveAll(util.UsedConfigurationDirectory)

	tests := []struct {
		name     string
		args     []string
		mockRepo map[string][]byte
		want     []string
	}{
		{
			name: "e module info",
			args: []string{"--configDir", util.UsedConfigurationDirectory, "module", "info", "example-repo/c1:0.1.0"},
			mockRepo: map[string][]byte{
				"example-repo.yaml": []byte(`version: v1
kind: Repository
name: example-repo
components:
  - name: c1
    type: docker
    versions:
      - version: 0.1.0
        latest: true
        image: "docker.io/hashicorp/terraform:0.12.28"
        workdir: "/terraform"
        mounts:
          - "/terraform"
        commands:
          - name: init
            description: "initializes terraform in local directory"
            command: init
            envs:
              TF_LOG: WARN
`),
			},
			want: []string{"Version: 0.1.0", "Image: docker.io/hashicorp/terraform:0.12.28"},
		},
		{
			name: "e module install",
			args: []string{"--configDir", util.UsedConfigurationDirectory, "module", "install", "example-repo/c1:0.1.0"},
			mockRepo: map[string][]byte{
				"example-repo.yaml": []byte(`version: v1
kind: Repository
name: example-repo
components:
  - name: c1
    type: docker
    versions:
      - version: 0.1.0
        latest: true
        image: "docker.io/hashicorp/terraform:0.12.28"
        workdir: "/terraform"
        mounts:
          - "/terraform"
        commands:
          - name: init
            description: "initializes terraform in local directory"
            command: init
            envs:
              TF_LOG: WARN
`),
			},
			want: []string{"Installed module c1:0.1.0 to environment"},
		},
		{
			name: "e module search",
			args: []string{"--configDir", util.UsedConfigurationDirectory, "module", "search", "c1"},
			mockRepo: map[string][]byte{
				"example-repo.yaml": []byte(`version: v1
kind: Repository
name: example-repo
components:
  - name: c1
    type: docker
    versions:
      - version: 0.1.0
        latest: true
        image: "docker.io/hashicorp/terraform:0.12.28"
        workdir: "/terraform"
        mounts:
          - "/terraform"
        commands:
          - name: init
            description: "initializes terraform in local directory"
            command: init
            envs:
              TF_LOG: WARN
`),
			},
			want: []string{"example-repo/c1:0.1.0"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := assert.New(t)
			dir, err := os.Getwd()
			a.NoError(err)

			if tt.mockRepo != nil {
				for k, v := range tt.mockRepo {
					err := ioutil.WriteFile(path.Join(util.UsedReposDirectory, k), v, 0644)
					a.NoError(err)
				}
			}

			cmd := exec.Command(path.Join(dir, "output", "e"), tt.args...)
			got, err := cmd.CombinedOutput()
			a.NoError(err)

			for _, w := range tt.want {
				a.Contains(string(got), w)
			}
		})
	}
}

func TestRepos(t *testing.T) {
	util.UsedConfigFile, util.UsedConfigurationDirectory, util.UsedEnvironmentDirectory, util.UsedReposDirectory, util.UsedTempDirectory = setup(t, "repos")
	defer os.RemoveAll(util.UsedConfigurationDirectory)

	tests := []struct {
		name string
		args []string
		want []string
	}{
		{
			name: "e repos list",
			args: []string{"--configDir", util.UsedConfigurationDirectory, "repos", "list"},
			want: []string{"Module: terraform:0.1.0"},
		},
		{
			name: "e repos install",
			args: []string{"--configDir", util.UsedConfigurationDirectory, "repos", "install", "mkyc/my-epipany-repo", "--logLevel", "debug"},
			want: []string{"will install mkyc/my-epipany-repo"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := assert.New(t)
			dir, err := os.Getwd()
			a.NoError(err)

			cmd := exec.Command(path.Join(dir, "output", "e"), tt.args...)
			got, err := cmd.CombinedOutput()
			a.NoError(err)

			for _, w := range tt.want {
				a.Contains(string(got), w)
			}
		})
	}
}

func TestSsh(t *testing.T) {
	util.UsedConfigFile, util.UsedConfigurationDirectory, util.UsedEnvironmentDirectory, util.UsedReposDirectory, util.UsedTempDirectory = setup(t, "ssh")
	defer os.RemoveAll(util.UsedConfigurationDirectory)

	tests := []struct {
		name string
		args []string
		want []string
	}{
		{
			name: "e ssh",
			args: []string{"--configDir", util.UsedConfigurationDirectory, "ssh"},
			want: []string{"Available Commands:\n  keygen"},
		},
		{
			name: "e ssh keygen",
			args: []string{"--configDir", util.UsedConfigurationDirectory, "ssh", "keygen"},
			want: []string{"Available Commands:\n  create"},
		},
		{
			name: "e ssh keygen create",
			args: []string{"--configDir", util.UsedConfigurationDirectory, "ssh", "keygen", "create", "--logLevel", "debug"},
			want: []string{"correctly saved private and public key files:"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := assert.New(t)
			dir, err := os.Getwd()
			a.NoError(err)

			cmd := exec.Command(path.Join(dir, "output", "e"), tt.args...)
			got, err := cmd.CombinedOutput()
			a.NoError(err)

			for _, w := range tt.want {
				a.Contains(string(got), w)
			}
		})
	}
}

func extractSubcommandsNames(in string) []string {
	commandsSectionExtractor := regexp.MustCompile("Available Commands:([\\S\\s]*?)Flags:")
	commandsSection := commandsSectionExtractor.FindString(in)
	commandsNamesExtractor := regexp.MustCompile("(?m)^\\s\\s[^\\s,]*[\\s]*")
	commandsNames := commandsNamesExtractor.FindAllString(commandsSection, -1)
	for i, m := range commandsNames {
		commandsNames[i] = strings.TrimSpace(m)
	}
	return commandsNames
}

func extractFlagsNames(input string) []string {
	useLineRemover := regexp.MustCompile("(?m)[\r\n]+^.*Use \"e.*$")
	inputWithoutUseLine := useLineRemover.ReplaceAllString(input, "")
	flagsSectionExtractor := regexp.MustCompile("Flags:([\\S\\s]*?)$")
	flagsSection := flagsSectionExtractor.FindString(inputWithoutUseLine)
	flagsNamesExtractor := regexp.MustCompile("--[a-zA-Z]*")
	flagsNames := flagsNamesExtractor.FindAllString(flagsSection, -1)
	for i, m := range flagsNames {
		flagsNames[i] = strings.TrimLeft(strings.TrimSpace(m), "-")
	}
	return flagsNames
}
