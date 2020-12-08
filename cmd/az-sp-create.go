package cmd

import (
	"fmt"
	"github.com/epiphany-platform/cli/pkg/az"

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
		fmt.Println("create pre run called")
	},
	Run: func(cmd *cobra.Command, args []string) {
		//TODO verify if environment is created and selected
		pass := az.GenerateServicePrincipalPassword()
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
