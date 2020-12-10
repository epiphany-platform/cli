package cmd

import (
	"fmt"
	"github.com/epiphany-platform/cli/pkg/az"
	"github.com/epiphany-platform/cli/pkg/configuration"
	"github.com/epiphany-platform/cli/pkg/environment"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	tenantID       string
	subscriptionID string
	name           string
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create Service Principal",
	Long:  `Create Service Principal that can be used for authentication with Azure.`,
	PreRun: func(cmd *cobra.Command, args []string) {
		logger.Debug().Msg("create pre run called")

		err := viper.BindPFlags(cmd.Flags())
		if err != nil {
			logger.Fatal().Err(err)
		}

		tenantID = viper.GetString("tenantID")
		subscriptionID = viper.GetString("subscriptionID")
		name = viper.GetString("name")

		if tenantID == "" {
			//TODO get default tenant
			logger.Fatal().Msg("no tenantID defined")
		}
		if subscriptionID == "" {
			//TODO get default subscription
			logger.Fatal().Msg("no subscriptionID defined")
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		if !isEnvPresentAndSelected() {
			logger.Fatal().Msg("no environment selected")
		}
		pass, err := az.GeneratePassword(32, 10)
		if err != nil {
			logger.Fatal().Err(err).Msg("failed to generate password")
		}
		app, _, err := az.CreateServicePrincipal(pass, subscriptionID, tenantID, name)
		if err != nil {
			logger.Fatal().Err(err).Msg("creation of service principal on Azure failed")
		}
		credentials := az.Credentials{
			AppID:          *app.AppID,
			Password:       pass,
			Tenant:         tenantID,
			SubscriptionID: subscriptionID,
		}
		logger.Debug().Msgf("prepared credentials for further consumption: %#v", credentials)
	},
}

func init() {
	spCmd.AddCommand(createCmd)

	createCmd.Flags().String("tenantID", "", fmt.Sprintf("TenantID of AAD where Service Principal should be created"))
	createCmd.Flags().String("subscriptionID", "", fmt.Sprintf("SubsciptionID of Subscription where Service Principal should have access"))
	createCmd.Flags().String("name", "epiphany-cli", fmt.Sprintf("Display Name of Service Principal"))
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
