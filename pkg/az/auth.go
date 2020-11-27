package az

import (
	"context"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/resources/mgmt/subscriptions"
	"github.com/Azure/azure-sdk-for-go/services/graphrbac/1.6/graphrbac"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/azure/auth"
)

const (
	cloudName        = "AzurePublicCloud"
	defaultPublisher = "Microsoft Services"
	roleName         = "Contributor"
	appName          = "20201031-auto-test-1"
)

// CreateSP function is used to create Service Principal
func CreateSP(tenantID, spName string) {
	info("Start creating of Azure Service Principal...")
	authorizer := getAuthrorizerFromCli()

	subscriptionsClient := subscriptions.NewClient()
	subscriptionsClient.Authorizer = authorizer

	env, err := azure.EnvironmentFromName(cloudName)
	if err != nil {
		log.Fatal(err)
	}

	tenantsClient := subscriptions.NewTenantsClient()
	tenantsClient.Authorizer = authorizer
	tenantsIterator, err := tenantsClient.ListComplete(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	tenantsCount := 0
	for tenantsIterator.NotDone() {
		tenantsCount++
		ten := tenantsIterator.Value()
		if *ten.TenantID == tenantID {
			log.Printf("Tenantid: %s, name: %s found.\n", *ten.TenantID, *ten.DisplayName)
			tenantID = *ten.TenantID
		}
		err = tenantsIterator.NextWithContext(context.TODO())
		if err != nil {
			log.Fatal(err)
		}
	}

	graphAuthorizer, err := auth.NewAuthorizerFromCLIWithResource(env.GraphEndpoint)

	info("Getting SP Client")
	spClient := graphrbac.NewServicePrincipalsClient(tenantID)
	spClient.Authorizer = graphAuthorizer

	spPresent := checkIfSPWithDisplayNamePresent(spClient, spName)

	if spPresent {
		info("Service Principal with name already exists.")
	} else {
		info("Service Principal with name doesn't exist.")
	}

}

func getAuthrorizerFromCli() autorest.Authorizer {
	cliAuthorizer, err := auth.NewAuthorizerFromCLI()
	if err != nil {
		log.Fatal("Ups...")
	} else {
		info("Got Azure CLI authorizer successfully .")
	}
	return cliAuthorizer
}

func checkIfSPWithDisplayNamePresent(spClient graphrbac.ServicePrincipalsClient, spName string) bool {
	info("Getting SP iterator")
	spFilter := "displayname eq '" + spName + "'"
	spIterator, err := spClient.ListComplete(context.TODO(), spFilter)
	if err != nil {
		log.Fatal(err)
	}

	start := time.Now()
	spCount := 0
	for spIterator.NotDone() {
		spCount++
		err = spIterator.NextWithContext(context.TODO())
		if err != nil {
			log.Fatal(err)
		}
	}
	elapsed := time.Since(start)
	log.Printf("SP listing took %s", elapsed)
	log.Printf("Found %d service principals.", spCount)
	if spCount > 0 {
		return true
	}
	return false

}
