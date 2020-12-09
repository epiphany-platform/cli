package az

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
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
	AppID          string
	Password       string
	Tenant         string
	SubscriptionID string
}

// CreateServicePrincipal function is used to create Service Principal, returns Service Principal and related App
func CreateServicePrincipal(pass, subscriptionID, tenantID, spName string) (graphrbac.ServicePrincipal, graphrbac.Application) {
	info("Start creating of Azure Service Principal...")
	resourceManagerAuthorizer := getAuthorizerFromCli()

	env := getEnvironment(cloudName)

	graphAuthorizer := getGraphAuthorizer(env)

	app := createApplication(tenantID, spName, pass, graphAuthorizer)

	sp := createServicePrincipal(tenantID, app, graphAuthorizer)

	assignRoleToServicePrincipal(subscriptionID, roleName, sp, resourceManagerAuthorizer)

	return sp, app
}

// GenerateServicePrincipalCredentialsStruct generate and returns Credentials structure
func GenerateServicePrincipalCredentialsStruct(pass, tenantID, subscriptionID, appID string) *Credentials {
	debug("Generate struct.")
	creds := &Credentials{
		AppID:          appID,
		Password:       pass,
		Tenant:         tenantID,
		SubscriptionID: subscriptionID,
	}
	debug("Creds %v: ", *creds)
	return creds
}

// GenerateServicePrincipalAuthJSONFromCredentialsStruct generate JSON that can be used for
func GenerateServicePrincipalAuthJSONFromCredentialsStruct(creds Credentials) []byte {
	debug("Marshaling to JSON.")
	credsJSON, err := json.Marshal(creds)
	if err != nil {
		errFailedToMarshalJSON(err)
	}
	debug(string(credsJSON))
	return credsJSON
}

// TODO fix that to add SP to environment and not to separate file
// WriteServicePrincipalAuthJSON to JSON authorization file
func WriteServicePrincipalAuthJSON(credsJSON []byte) {
	err := ioutil.WriteFile("/tmp/dat1", credsJSON, 0644)
	if err != nil {
		errFailedToWriteJSONAuthFile(err)
	}
}

// GenerateServicePrincipalPassword generates Service Principal password
func GenerateServicePrincipalPassword() string {
	pass, err := password.Generate(32, 10, 0, false, false)
	if err != nil {
		errFailedToGeneratePassword(err)
	}
	return pass
}

// getAuthorizerFromCli returns authorizer based on local az login session
func getAuthorizerFromCli() autorest.Authorizer {
	cliAuthorizer, err := auth.NewAuthorizerFromCLI()
	if err != nil {
		errFailedToGetAuthorizerFromCli(err)
	} else {
		info("Got Azure CLI authorizer successfully .")
	}
	return cliAuthorizer
}

// TODO consider removal
// getEnvironment returns Azure Environment based on cloudName
func getEnvironment(cloudName string) azure.Environment {
	env, err := azure.EnvironmentFromName(cloudName)
	if err != nil {
		errFailedToGetEnvironment(err)
	}
	return env
}

// getGraphAuthorizer return graph authorizer based on Azure Environment
func getGraphAuthorizer(env azure.Environment) autorest.Authorizer {
	graphAuthorizer, err := auth.NewAuthorizerFromCLIWithResource(env.GraphEndpoint)
	if err != nil {
		errFailedToGetGraphAuthorizer(err)
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
		errFailedToCreateApplication(err)
	}
	appJSON, err := app.MarshalJSON()
	if err != nil {
		errFailedToMarshalJSON(err)
	}
	debug(fmt.Sprint("App: ", string(appJSON)))
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
		errFailedToCreateApplication(err)
	}

	fmt.Println("objectID: ", *sp.ObjectID)
	spJSON, err := sp.MarshalJSON()
	if err != nil {
		errFailedToMarshalJSON(err)
	}
	debug(fmt.Sprint("SP: ", string(spJSON)))
	fmt.Println(fmt.Sprint("SP: ", string(spJSON)))
	return sp
}

// assignRoleToServicePrincipal assigns role from RBAC to Service Principal
// based on subscriptionID string, sp graphrbac.ServicePrincipal, resourceManagerAuthorizer autorest.Authorizer
func assignRoleToServicePrincipal(subscriptionID, roleName string, sp graphrbac.ServicePrincipal, resourceManagerAuthorizer autorest.Authorizer) {
	info("Assigning a role to Service Principal")
	roleAssignmentClient := authorization.NewRoleAssignmentsClient(subscriptionID)
	roleAssignmentClient.Authorizer = resourceManagerAuthorizer

	roleID := getRoleID(subscriptionID, roleName, resourceManagerAuthorizer)

	roleAssignmentName := uuid.NewV4()
	for i := 0; i < 30; i++ {
		ra, err := roleAssignmentClient.Create(context.TODO(), "/subscriptions/"+subscriptionID, roleAssignmentName.String(), authorization.RoleAssignmentCreateParameters{
			Properties: &authorization.RoleAssignmentProperties{
				RoleDefinitionID: to.StringPtr(roleID),
				PrincipalID:      sp.ObjectID,
			},
		})
		if err != nil {
			warnAssignRoleToServicePrincipal(err)
			time.Sleep(1 * time.Second)
			continue
		} else {
			raJSON, err := ra.MarshalJSON()
			if err != nil {
				errFailedToMarshalJSON(err)
			}
			debug(fmt.Sprintf("\n===========\nROLE ASSIGNMENT\n%v\n===========\n", string(raJSON)))
			break
		}
	}
}

// getRoleID finds roleID that is equal to roleName from given subscription
func getRoleID(subscriptionID, roleName string, resourceManagerAuthorizer autorest.Authorizer) string {

	roleDefinitionClient := authorization.NewRoleDefinitionsClient(subscriptionID)
	roleDefinitionClient.Authorizer = resourceManagerAuthorizer

	var roleID string

	roleDefinitionIterator, err := roleDefinitionClient.ListComplete(context.TODO(), "/subscriptions/"+subscriptionID, "")
	if err != nil {
		errFailedToGetRoleDefinitionIterator(err)
	}

	for roleDefinitionIterator.NotDone() {
		rd := roleDefinitionIterator.Value()
		if *rd.RoleName == roleName {
			roleID = *rd.ID
			debug(fmt.Sprintf("RoleDefinition: %s\n", *rd.RoleName))
		}
		err = roleDefinitionIterator.NextWithContext(context.TODO())
		if err != nil {
			errFailedToIterateOverRoleDefinitions(err)
		}
	}
	return roleID
}
