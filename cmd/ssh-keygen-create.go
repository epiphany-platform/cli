package cmd

import (
	"path"

	"github.com/epiphany-platform/cli/internal/logger"
	"github.com/epiphany-platform/cli/internal/util"
	"github.com/epiphany-platform/cli/pkg/auth"
	"github.com/spf13/cobra"
)

// sshKeygenCreateCmd represents the create command
var sshKeygenCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create new ssh keypair in current environment.",
	Long:  `This command creates new ssh keypair in current environment and stores information about it in environment.`,
	PreRun: func(cmd *cobra.Command, args []string) {
		logger.Debug().Msg("create called")
	},
	Run: func(cmd *cobra.Command, args []string) {
		kp, err := auth.GenerateRsaKeyPair(path.Join(
			util.UsedEnvironmentDirectory,
			currentEnvironment.Uuid.String(),
			"/shared", //TODO to consts
		))
		if err != nil {
			logger.Fatal().Err(err).Msg("generate rsa keypair failed")
		}
		currentEnvironment.AddRsaKeyPair(kp)
		err = currentEnvironment.Save()
		if err != nil {
			logger.Fatal().Err(err).Msg("save env failed")
		}
	},
}

func init() {
	sshKeygenCmd.AddCommand(sshKeygenCreateCmd)
}
