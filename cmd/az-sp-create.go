package cmd

import (
	"errors"
	"github.com/epiphany-platform/cli/pkg/az"
	"github.com/epiphany-platform/cli/pkg/configuration"
	"github.com/epiphany-platform/cli/pkg/environment"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	tenantID                string
	subscriptionID          string
	newServicePrincipalName string
)

// azSpCreateCmd represents the create command
var azSpCreateCmd = &cobra.Command{
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
		newServicePrincipalName = viper.GetString("name")

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
		config, err := isEnvPresentAndSelected()
		if err != nil {
			logger.Fatal().Msg("no environment selected")
		}
		pass, err := az.GeneratePassword(32, 10, 5)
		if err != nil {
			logger.Fatal().Err(err).Msg("failed to generate password")
		}
		app, _, err := az.CreateServicePrincipal(pass, subscriptionID, tenantID, newServicePrincipalName)
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
		config.AddAzureCredentials(credentials)
		err = config.Save()
		if err != nil {
			logger.Fatal().Err(err).Msg("failed to save config file")
		}
	},
}

func init() {
	spCmd.AddCommand(azSpCreateCmd)

	azSpCreateCmd.Flags().String("tenantID", "", "TenantID of AAD where service principal should be created")
	azSpCreateCmd.Flags().String("subscriptionID", "", "SubscriptionID of subscription where service principal should have access")
	azSpCreateCmd.Flags().String("name", "epiphany-cli", "Display Name of service principal")
}

func isEnvPresentAndSelected() (config *configuration.Config, err error) {
	debug("will check if isEnvPresentAndSelected()")
	config, err = configuration.GetConfig()
	if err != nil {
		return
	}
	environments, err := environment.GetAll()
	if err != nil {
		return
	}
	for _, e := range environments {
		if e.Uuid.String() == config.CurrentEnvironment.String() {
			debug("found currently selected environment %s", e.Uuid.String())
			return
		}
	}
	debug("currently selected environment not found")
	err = errors.New("currently selected environment not found")
	return
}
