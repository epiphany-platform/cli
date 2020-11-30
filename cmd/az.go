/*

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
	tenantID      string
	subsciptionID string
	spName        string
)

// azCmd represents the az command
var azCmd = &cobra.Command{
	Use:   "az",
	Short: "Enable access to set of commands used to work with Azure cloud",
	Long: `Enable access to set of commands used to work with Azure cloud:
	- authentication - let you access authentication options - e.g. create Service Principal`,
	PreRun: func(cmd *cobra.Command, args []string) {
		debug("az called")
	},
	Run: func(cmd *cobra.Command, args []string) {
		az.CreateSP(tenantID, subsciptionID, spName)
	},
}

func init() {
	rootCmd.AddCommand(azCmd)
	azCmd.PersistentFlags().StringVar(&tenantID, "tenantID", "", fmt.Sprintf("TenantID of AAD where Service Principal should be created"))
	azCmd.PersistentFlags().StringVar(&subsciptionID, "subsciptionID", "", fmt.Sprintf("SubsciptionID of Subscription where Service Principal should have access"))
	azCmd.PersistentFlags().StringVar(&spName, "spName", "", fmt.Sprintf("Display Name of Service Principal"))
}
