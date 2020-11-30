package az

import (
	"context"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/authorization/mgmt/authorization"
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
)

// Credentials structure is used to format display information for Service Principal
type Credentials struct {
	appID          string
	password       string
	tenant         string
	subscriptionID string
}

// CreateSP function is used to create Service Principal
func CreateSP(subscriptionID, tenantID, spName string) {
	info("Start creating of Azure Service Principal...")
	resourceManagerAuthorizer := getAuthrorizerFromCli()

	env := getEnvironment(cloudName)

	graphAuthorizer := getGraphAuthorizer(env)

	pass := generatePassword()

	app := createApplication(tenantID, spName, pass, graphAuthorizer)

	sp := createServicePrincipal(tenantID, app, graphAuthorizer)

	assignRoleToServicePrincipal(subscriptionID, sp, resourceManagerAuthorizer)

	creds := &Credentials{
		appID:          *sp.AppID,
		password:       pass,
		tenant:         tenantID,
		subscriptionID: subscriptionID,
	}
	log.Printf("\n===========\nCREDENCIALS\n%+v\n===========\n", creds)
	info("Azure Service Principal created.")
}

func generatePassword() string {
	pass, err := password.Generate(32, 10, 0, false, false)
	if err != nil {
		log.Fatal(err)
	}
	return pass
}

// getAuthrorizerFromCli returns authorizer based on local az login session
func getAuthrorizerFromCli() autorest.Authorizer {
	cliAuthorizer, err := auth.NewAuthorizerFromCLI()
	if err != nil {
		log.Fatal("Ups...", err)
	} else {
		info("Got Azure CLI authorizer successfully .")
	}
	return cliAuthorizer
}

// getEnvironment returns Azure Environment based on cloudName
func getEnvironment(cloudName string) azure.Environment {
	env, err := azure.EnvironmentFromName(cloudName)
	if err != nil {
		log.Fatal(err)
	}
	return env
}

// getGraphAuthorizer return graph authorizer based on Azure Environment
func getGraphAuthorizer(env azure.Environment) autorest.Authorizer {
	graphAuthorizer, err := auth.NewAuthorizerFromCLIWithResource(env.GraphEndpoint)
	if err != nil {
		log.Fatal("Ups...", err)
	} else {
		info("Got Azure Graph authorizer successfully .")
	}
	return graphAuthorizer
}

// createApplication creates an application that is used with Service Principal based on tenantID, spName, pass and graphAuthorizer (autorest.Authorizer)
func createApplication(tenantID, spName, pass string, graphAuthorizer autorest.Authorizer) graphrbac.Application {
	info("Creating an application")
	appClient := graphrbac.NewApplicationsClient(tenantID)
	appClient.Authorizer = graphAuthorizer

	keyID := uuid.NewV4()
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

	return app
}

// createServicePrincipal creates Service Principal based on tenantID, application (graphrbac.Application) using graphAuthorizer (autorest.Authorizer)
func createServicePrincipal(tenantID string, app graphrbac.Application, graphAuthorizer autorest.Authorizer) graphrbac.ServicePrincipal {
	info("Creating a Service Principal")
	spClient := graphrbac.NewServicePrincipalsClient(tenantID)
	spClient.Authorizer = graphAuthorizer

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
	return sp
}

// assignRoleToServicePrincipal assigns role from RBAC to Service Principal
// based on subscriptionID string, sp graphrbac.ServicePrincipal, resourceManagerAuthorizer autorest.Authorizer
func assignRoleToServicePrincipal(subscriptionID string, sp graphrbac.ServicePrincipal, resourceManagerAuthorizer autorest.Authorizer) {
	info("Assigning a role to Service Principal")
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
}
