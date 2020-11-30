package az

import (
	"context"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/authorization/mgmt/authorization"
	"github.com/Azure/azure-sdk-for-go/profiles/latest/resources/mgmt/subscriptions"
	"github.com/Azure/azure-sdk-for-go/services/graphrbac/1.6/graphrbac"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/Azure/go-autorest/autorest/date"
	"github.com/Azure/go-autorest/autorest/to"
	uuid "github.com/satori/go.uuid"
	"github.com/sethvargo/go-password/password"
)

const (
	cloudName        = "AzurePublicCloud"
	defaultPublisher = "Microsoft Services"
	roleName         = "Contributor"
	appName          = ""
)

type Credentials struct {
	appID        string
	password     string
	tenant       string
	subscription string
}

// CreateSP function is used to create Service Principal
func CreateSP(tenantID, subscriptionID, spName string) {
	info("Start creating of Azure Service Principal...")
	resourceManagerAuthorizer := getAuthrorizerFromCli()

	subscriptionsClient := subscriptions.NewClient()
	subscriptionsClient.Authorizer = resourceManagerAuthorizer

	env, err := azure.EnvironmentFromName(cloudName)
	if err != nil {
		log.Fatal(err)
	}

	tenantsClient := subscriptions.NewTenantsClient()
	tenantsClient.Authorizer = resourceManagerAuthorizer
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

	// application
	appClient := graphrbac.NewApplicationsClient(tenantID)
	appClient.Authorizer = graphAuthorizer

	keyID := uuid.NewV4()
	pass, err := password.Generate(32, 10, 0, false, false)
	if err != nil {
		log.Fatal(err)
	}
	t := &date.Time{
		Time: time.Now(),
	}
	t2 := &date.Time{
		Time: t.AddDate(2, 0, 0),
	}
	app, err := appClient.Create(context.TODO(), graphrbac.ApplicationCreateParameters{
		DisplayName:             to.StringPtr(spName),
		IdentifierUris:          &[]string{"https://" + spName},
		AvailableToOtherTenants: to.BoolPtr(false),
		Homepage:                to.StringPtr("https://" + spName),
		PasswordCredentials: &[]graphrbac.PasswordCredential{{
			StartDate:           t,
			EndDate:             t2,
			KeyID:               to.StringPtr(keyID.String()),
			Value:               to.StringPtr(pass),
			CustomKeyIdentifier: to.ByteSlicePtr([]byte(spName)),
		}},
	})
	if err != nil {
		log.Fatal(err)
	}
	appJSON, err := app.MarshalJSON()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("App: ", string(appJSON))

	// sp creation
	sp, err := spClient.Create(context.TODO(), graphrbac.ServicePrincipalCreateParameters{
		AppID:          app.AppID,
		AccountEnabled: to.BoolPtr(true),
	})
	if err != nil {
		log.Fatal(err)
	}
	spJSON, err := sp.MarshalJSON()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("App: ", string(spJSON))

	roleAssignmentClient := authorization.NewRoleAssignmentsClient(subscriptionID)
	roleAssignmentClient.Authorizer = resourceManagerAuthorizer

	var roleID string
	roleAssignmentName := uuid.NewV4()
	for i := 0; i < 30; i++ {
		ra, err := roleAssignmentClient.Create(context.TODO(), "/subscriptions/"+subscriptionID, roleAssignmentName.String(), authorization.RoleAssignmentCreateParameters{
			Properties: &authorization.RoleAssignmentProperties{
				RoleDefinitionID: to.StringPtr(roleID),
				PrincipalID:      sp.ObjectID,
			},
		})
		if err != nil {
			log.Println(err)
			time.Sleep(1 * time.Second)
			continue
		} else {
			raJSON, err := ra.MarshalJSON()
			if err != nil {
				log.Fatal(err)
			}
			log.Printf("\n===========\nROLE ASSIGNMENT\n%v\n===========\n", string(raJSON))
			break
		}
	}

	c := &Credentials{
		appID:        *sp.AppID,
		password:     pass,
		tenant:       tenantID,
		subscription: subscriptionID,
	}
	log.Printf("\n===========\nCREDENCIALS\n%+v\n===========\n", c)

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
