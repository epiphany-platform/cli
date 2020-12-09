package cmd

import (
	"fmt"
	"github.com/epiphany-platform/cli/pkg/az"
	"github.com/epiphany-platform/cli/pkg/configuration"
	"github.com/epiphany-platform/cli/pkg/environment"
	"github.com/spf13/cobra"
)

var (
	tenantID       string
	subscriptionID string
	spName         string
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create Service Principal",
	Long:  `Create Service Principal that can be used for authentication with Azure.`,
	PreRun: func(cmd *cobra.Command, args []string) {
		debug("create pre run called")
	},
	Run: func(cmd *cobra.Command, args []string) {
		if !isEnvPresentAndSelected() {
			logger.Fatal().Msg("no environment selected")
		}
		pass, err := az.GeneratePassword(32, 10)
		if err != nil {
			logger.Panic().Err(err).Msg("failed to generate password")
		}
		sp, app := az.CreateServicePrincipal(pass, subscriptionID, tenantID, spName)
		debug("Create Service Principal with ObjectID: %s, AppID: %s", *sp.ObjectID, *sp.AppID)
		az.GenerateServicePrincipalCredentialsStruct(pass, tenantID, subscriptionID, *app.AppID)
	},
}

func init() {
	spCmd.AddCommand(createCmd)

	createCmd.PersistentFlags().StringVar(&tenantID, "tenantID", "", fmt.Sprintf("TenantID of AAD where Service Principal should be created"))
	createCmd.PersistentFlags().StringVar(&subscriptionID, "subscriptionID", "", fmt.Sprintf("SubsciptionID of Subscription where Service Principal should have access"))
	createCmd.PersistentFlags().StringVar(&spName, "name", "", fmt.Sprintf("Display Name of Service Principal"))
}

func isEnvPresentAndSelected() bool {
	debug("will check if isEnvPresentAndSelected()")
	config, err := configuration.GetConfig()
	if err != nil {
		logger.Fatal().Err(err).Msg("get config failed")
	}
	environments, err := environment.GetAll()
	if err != nil {
		logger.Fatal().Err(err).Msg("environments get all failed")
	}
	for _, e := range environments {
		if e.Uuid.String() == config.CurrentEnvironment.String() {
			debug("found currently selected environment %s", e.Uuid.String())
			return true
		}
	}
	debug("currently selected environment not found")
	return false
}
