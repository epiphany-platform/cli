package az

import (
	"context"
	"fmt"
	"github.com/rs/zerolog"
	"os"
	"time"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/authorization/mgmt/authorization"
	"github.com/Azure/azure-sdk-for-go/services/graphrbac/1.6/graphrbac"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/Azure/go-autorest/autorest/date"
	"github.com/Azure/go-autorest/autorest/to"
	"github.com/google/uuid"
	"github.com/sethvargo/go-password/password"
)

const (
	cloudName = "AzurePublicCloud"
	roleName  = "Contributor"
)

var (
	logger zerolog.Logger
)

func init() {
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	logger = zerolog.New(output).With().Str("package", "az").Caller().Timestamp().Logger()
}

// Credentials structure is used to format display information for Service Principal
type Credentials struct {
	AppID          string
	Password       string
	Tenant         string
	SubscriptionID string
}

// TODO CreateServicePrincipal has to be CLI independent for test reasons
// CreateServicePrincipal function is used to create Service Principal, returns Service Principal and related App
func CreateServicePrincipal(pass, subscriptionID, tenantID, name string) (app graphrbac.Application, sp graphrbac.ServicePrincipal, err error) {
	logger.Debug().Msg("begin CreateServicePrincipal(...)")

	authorizer, err := auth.NewAuthorizerFromCLI()
	if err != nil {
		return
	}
	logger.Debug().Msg("authorizer created from az cli command")

	env, err := azure.EnvironmentFromName(cloudName)
	if err != nil {
		return
	}
	logger.Debug().Msg("azure cloud endpoints information obtained")

	graphAuthorizer, err := auth.NewAuthorizerFromCLIWithResource(env.GraphEndpoint)
	if err != nil {
		return
	}
	logger.Debug().Msg("graph authorizer created from az cli command")

	app, err = createApplication(tenantID, name, pass, graphAuthorizer)
	if err != nil {
		return
	}
	appBytes, err := app.MarshalJSON()
	if err != nil {
		logger.Warn().Err(err).Msg("wasn't able to marshall application structure")
	}
	logger.Debug().Msgf("created application: \n %s", string(appBytes))

	sp, err = createServicePrincipal(tenantID, app, graphAuthorizer)
	if err != nil {
		return
	}
	spBytes, err := sp.MarshalJSON()
	if err != nil {
		logger.Warn().Err(err).Msg("wasn't able to marshall service principal structure")
	}
	logger.Debug().Msgf("created service principal: \n %s", string(spBytes))

	roleID, err := getRoleID(subscriptionID, roleName, authorizer)
	if err != nil {
		return
	}
	logger.Debug().Msgf("obtained id %s for role %s", roleID, roleName)

	ra, err := assignRoleToServicePrincipalWithRetries(subscriptionID, roleID, sp, authorizer)
	if err != nil {
		return
	}
	raBytes, err := ra.MarshalJSON()
	if err != nil {
		logger.Warn().Err(err).Msg("wasn't able to marshall role assignment structure")
	}
	logger.Debug().Msgf("created role assignment: \n %s", string(raBytes))

	logger.Debug().Msg("end CreateServicePrincipal(...)")
	return
}

// GeneratePassword generates Service Principal password
func GeneratePassword(length, numDigits int) (string, error) {
	logger.Debug().Msgf("will generate password of length %d with %d digits", length, numDigits)
	if numDigits > length {
		return "", fmt.Errorf("parameter 'numDigits' cannot be greater than parameter 'length'")
	}
	pass, err := password.Generate(length, numDigits, 0, false, false)
	if err != nil {
		return "", err
	}
	logger.Debug().Msgf("generated password was: %s", pass)
	return pass, nil
}

// createApplication creates an application that is used with Service Principal based on tenantID, name, pass and graphAuthorizer (autorest.Authorizer)
func createApplication(tenantID, name, password string, authorizer autorest.Authorizer) (graphrbac.Application, error) {
	logger.Debug().Msg("will create application")
	client := graphrbac.NewApplicationsClient(tenantID)
	client.Authorizer = authorizer

	keyID := uuid.New()
	t := &date.Time{
		Time: time.Now(),
	}
	t2 := &date.Time{
		Time: t.AddDate(2, 0, 0),
	}
	return client.Create(context.TODO(), graphrbac.ApplicationCreateParameters{
		DisplayName:             to.StringPtr(name),
		IdentifierUris:          &[]string{"https://" + name},
		AvailableToOtherTenants: to.BoolPtr(false),
		Homepage:                to.StringPtr("https://" + name),
		PasswordCredentials: &[]graphrbac.PasswordCredential{{
			StartDate:           t,
			EndDate:             t2,
			KeyID:               to.StringPtr(keyID.String()),
			Value:               to.StringPtr(password),
			CustomKeyIdentifier: to.ByteSlicePtr([]byte(name)),
		}},
	})
}

// createServicePrincipal creates Service Principal based on tenantID, application (graphrbac.Application) using graphAuthorizer (autorest.Authorizer)
func createServicePrincipal(tenantID string, app graphrbac.Application, authorizer autorest.Authorizer) (graphrbac.ServicePrincipal, error) {
	logger.Debug().Msg("will create service principal")
	client := graphrbac.NewServicePrincipalsClient(tenantID)
	client.Authorizer = authorizer

	return client.Create(context.TODO(), graphrbac.ServicePrincipalCreateParameters{
		AppID:          app.AppID,
		AccountEnabled: to.BoolPtr(true),
	})
}

// assignRoleToServicePrincipalWithRetries assigns role from RBAC to Service Principal
// based on subscriptionID string, sp graphrbac.ServicePrincipal, resourceManagerAuthorizer autorest.Authorizer
func assignRoleToServicePrincipalWithRetries(subscriptionID, roleID string, sp graphrbac.ServicePrincipal, authorizer autorest.Authorizer) (ra authorization.RoleAssignment, err error) {
	logger.Debug().Msg("will assign role to service principal")
	client := authorization.NewRoleAssignmentsClient(subscriptionID)
	client.Authorizer = authorizer

	for i := 0; i < 30; i++ {
		ra, err = client.Create(context.TODO(), "/subscriptions/"+subscriptionID, uuid.New().String(), authorization.RoleAssignmentCreateParameters{
			Properties: &authorization.RoleAssignmentProperties{
				RoleDefinitionID: to.StringPtr(roleID),
				PrincipalID:      sp.ObjectID,
			},
		})
		if err != nil {
			logger.Warn().Err(err).Msgf("(%d) failed to assign role to service principal.", i)
			time.Sleep(1 * time.Second)
			continue
		} else {
			return
		}
	}
	err = fmt.Errorf("maximum retries count achieved")
	return
}

// getRoleID finds roleID that is equal to roleName from given subscription
func getRoleID(subscriptionID, roleName string, authorizer autorest.Authorizer) (roleID string, err error) {
	logger.Debug().Msg("will search for role")
	client := authorization.NewRoleDefinitionsClient(subscriptionID)
	client.Authorizer = authorizer

	roleDefinitionIterator, err := client.ListComplete(context.TODO(), "/subscriptions/"+subscriptionID, "")
	if err != nil {
		return
	}

	for roleDefinitionIterator.NotDone() {
		rd := roleDefinitionIterator.Value()
		if *rd.RoleName == roleName {
			roleID = *rd.ID
			return
		}
		err = roleDefinitionIterator.NextWithContext(context.TODO())
		if err != nil {
			return
		}
	}
	err = fmt.Errorf("required role not found")
	return
}