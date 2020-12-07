/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
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

// createServicePrincipalCmd represents the createServicePrincipal command
var createServicePrincipalCmd = &cobra.Command{
	Use:   "create-service-principal",
	Short: "Create Service Principal",
	Long:  `Create Service Principal that can be used for authentication with Azure.`,
	PreRun: func(cmd *cobra.Command, args []string) {
		fmt.Println("createServicePrincipal called")
	},
	Run: func(cmd *cobra.Command, args []string) {
		pass := az.GenerateServicePrincipalPassword()
		sp, app := az.CreateServicePrincipal(pass, subscriptionID, tenantID, spName)
		debug("Create Service Principal with ObjectID: %s, AppID: %s", *sp.ObjectID, *sp.AppID)
		az.GenerateServicePrincipalCredentialsStruct(pass, tenantID, subscriptionID, *app.AppID)
	},
}

func init() {
	azCmd.AddCommand(createServicePrincipalCmd)

	createServicePrincipalCmd.PersistentFlags().StringVar(&tenantID, "tenantID", "", fmt.Sprintf("TenantID of AAD where Service Principal should be created"))
	createServicePrincipalCmd.PersistentFlags().StringVar(&subscriptionID, "subscriptionID", "", fmt.Sprintf("SubsciptionID of Subscription where Service Principal should have access"))
	createServicePrincipalCmd.PersistentFlags().StringVar(&spName, "spName", "", fmt.Sprintf("Display Name of Service Principal"))
}
