package cmd

import (
	"github.com/epiphany-platform/cli/internal/logger"
	"github.com/epiphany-platform/cli/pkg/auth"
	"github.com/epiphany-platform/cli/pkg/configuration"
	"github.com/epiphany-platform/cli/pkg/environment"
	"github.com/epiphany-platform/cli/pkg/util"
	"path"

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
		config, err := configuration.GetConfig()
		if err != nil {
			logger.Fatal().Err(err).Msg("get config failed")
		}
		env, err := environment.Get(config.CurrentEnvironment)
		if err != nil {
			logger.Fatal().Err(err).Msg("get environments details failed")
		}
		kp, err := auth.GenerateRsaKeyPair(path.Join(
			util.UsedEnvironmentDirectory,
			env.Uuid.String(),
			"/shared", //TODO to consts
		))
		if err != nil {
			logger.Fatal().Err(err).Msg("generate rsa keypair failed")
		}
		env.AddRsaKeyPair(kp)
		err = env.Save()
		if err != nil {
			logger.Fatal().Err(err).Msg("save env failed")
		}
	},
}

func init() {
	sshKeygenCmd.AddCommand(sshKeygenCreateCmd)
}
