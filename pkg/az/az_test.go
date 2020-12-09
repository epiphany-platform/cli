// +build integration

package az

import (
	"context"
	"os"
	"testing"
	"unicode"

	"github.com/Azure/azure-sdk-for-go/services/graphrbac/1.6/graphrbac"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/azure/auth"
)

const (
	cloudTestName        = "AzurePublicCloud"
	defaultTestPublisher = "Microsoft Services"
	roleTestName         = "Contributor"
)

var (
	tenantID       string
	subscriptionID string
	spName         string
)

// setupTestAll setups parameters needed to run all tests
func setupTestAll(t *testing.T) {
	t.Log("Initialize test")
	tenantID = os.Getenv("TENANT_ID")
	if len(tenantID) == 0 {
		t.Fatalf("expected non-empty TENANT_ID environment variable")
	}

	subscriptionID = os.Getenv("SUBSCRIPTION_ID")
	if len(subscriptionID) == 0 {
		t.Fatalf("expected non-empty SUBSCRIPTION_ID environment variable")
	}

	spName = os.Getenv("SP_NAME")
	if len(spName) == 0 {
		t.Fatalf("expected non-empty SP_NAME environment variable")
	}
}

// cleanupTestServicePrincipal cleans up Service Principal and related resources based on app and sp object ID
func cleanupTestServicePrincipal(spObjectID, appObjectID string, t *testing.T) {
	t.Log("Start deleting Service Prnicipal.")
	spClient := getTestServicePrincipalClient(tenantID, cloudTestName, t)
	appClient := getTestAppClient(tenantID, cloudTestName, t)

	_, err := spClient.Delete(context.TODO(), spObjectID)
	catch(err, t)

	_, err = appClient.Delete(context.TODO(), appObjectID)
	catch(err, t)
}

// func TestShouldSuccessfullyGeneratePassword(t *testing.T) {
// 	// when
// 	pass := GeneratePassword()

// 	// then
// 	if len(pass) < 32 {
// 		t.Error("Generated password too short.")
// 	}
// }

// func TestShouldSuccessfullyCreateServicePrincipal(t *testing.T) {

// 	// when
// 	pass := GeneratePassword()

// 	sp, app := CreateServicePrincipal(pass, subscriptionID, tenantID, spName)

// 	spClient := getTestServicePrincipalClient(tenantID, cloudTestName, t)

// 	spTest, err := spClient.Get(context.TODO(), *sp.ObjectID)
// 	catch(err, t)

// 	// then
// 	spJSON, err := sp.MarshalJSON()
// 	catch(err, t)

// 	t.Log(fmt.Sprint("App: ", string(spJSON)))

// 	if *spTest.ObjectID != *spTest.ObjectID {
// 		t.Error("Different object ID of Service Principal.")
// 	}

// 	if *sp.AppID != *spTest.AppID {
// 		t.Error("Different object ID of Service Principal.")
// 	}

// 	appClient := getTestAppClient(tenantID, cloudTestName, t)

// 	appTest, err := appClient.Get(context.TODO(), *app.ObjectID)
// 	catch(err, t)

// 	if *app.AppID != *appTest.AppID {
// 		t.Error("Different AppID of application.")
// 	}

// 	if *app.DisplayName != *appTest.DisplayName {
// 		t.Error("Different DisplayName of application.")
// 	}

// 	cleanupTestServicePrincipal(*sp.ObjectID, *app.ObjectID, t)
// }

// func TestShouldSuccessfullyCreateServicePrincipalCredentialsStruct(t *testing.T) {

// 	// when
// 	pass := GeneratePassword()
// 	sp, app := CreateServicePrincipal(pass, subscriptionID, tenantID, spName)

// 	creds := GenerateServicePrincipalCredentialsStruct(pass, tenantID, subscriptionID, *app.AppID)

// 	// then
// 	if creds.appID != *sp.AppID {
// 		t.Error("Different AppID in creds.")
// 	}

// 	if creds.password != pass {
// 		t.Error("Different password in creds.")
// 	}

// 	if creds.subscriptionID != subscriptionID {
// 		t.Error("Different subscriptionID in creds.")
// 	}

// 	if creds.tenant != tenantID {
// 		t.Error("Different tenantID in creds.")
// 	}

// 	cleanupTestServicePrincipal(*sp.ObjectID, *app.ObjectID, t)
// }

func TestShouldSuccessfullyCreateServicePrincipalAuthJSONIntegration(t *testing.T) {

	setupTestAll(t)

	if testing.Short() {
		t.Skip("skipping integration test")
	}

	// when
	pass, err := GeneratePassword(32, 10)
	if err != nil {
		t.Fatal(err)
	}
	//sp, app := CreateServicePrincipal(pass, subscriptionID, tenantID, spName)

	appID := "111111111111"

	//creds := GenerateServicePrincipalCredentialsStruct(pass, tenantID, subscriptionID, *app.AppID)
	creds := GenerateServicePrincipalCredentialsStruct(pass, tenantID, subscriptionID, appID)

	// then
	spJSON := GenerateServicePrincipalAuthJSONFromCredentialsStruct(*creds)
	t.Log(string(spJSON))

	//cleanupTestServicePrincipal(*sp.ObjectID, *app.ObjectID, t)
}

func TestShouldSuccessfullyWriteAuthJSONToFileIntegration(t *testing.T) {

	setupTestAll(t)

	if testing.Short() {
		t.Skip("skipping integration test")
	}

	// when
	pass, err := GeneratePassword(32, 10)
	if err != nil {
		t.Fatal(err)
	}
	//sp, app := CreateServicePrincipal(pass, subscriptionID, tenantID, spName)

	//creds := GenerateServicePrincipalCredentialsStruct(pass, tenantID, subscriptionID, *app.AppID)
	creds := &Credentials{
		AppID:          "appID",
		Password:       pass,
		Tenant:         "TenantID",
		SubscriptionID: "SubscriptionID",
	}

	// then
	spJSON := GenerateServicePrincipalAuthJSONFromCredentialsStruct(*creds)
	t.Log(string(spJSON))

	//cleanupTestServicePrincipal(*sp.ObjectID, *app.ObjectID, t)
}

// getEnvironment returns Azure Environment based on cloudName
func getTestEnvironment(cloudName string, t *testing.T) azure.Environment {
	env, err := azure.EnvironmentFromName(cloudName)
	catch(err, t)
	return env
}

// getGraphAuthorizer return graph authorizer based on Azure Environment
func getTestGraphAuthorizer(env azure.Environment, t *testing.T) autorest.Authorizer {
	graphAuthorizer, err := auth.NewAuthorizerFromCLIWithResource(env.GraphEndpoint)
	catch(err, t)
	return graphAuthorizer
}

// getTestAppClient gets application client for test purposes
func getTestAppClient(tenantID, cloudTestName string, t *testing.T) graphrbac.ApplicationsClient {
	env := getTestEnvironment(cloudTestName, t)
	graphAuthorizer := getTestGraphAuthorizer(env, t)

	appClient := graphrbac.NewApplicationsClient(tenantID)
	appClient.Authorizer = graphAuthorizer

	return appClient
}

// getTestServicePrincipalClient gets service principal client for test purposes
func getTestServicePrincipalClient(tenantID, cloudTestName string, t *testing.T) graphrbac.ServicePrincipalsClient {
	env := getTestEnvironment(cloudTestName, t)
	graphAuthorizer := getTestGraphAuthorizer(env, t)

	spClient := graphrbac.NewServicePrincipalsClient(tenantID)
	spClient.Authorizer = graphAuthorizer

	return spClient
}

func catch(err error, t *testing.T) {
	if err != nil {
		t.Fatal(err)
	}
}

func TestGeneratePassword(t *testing.T) {
	type args struct {
		length    int
		numDigits int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "happy path",
			args: args{
				length:    32,
				numDigits: 10,
			},
			wantErr: false,
		},
		{
			name: "too long pass",
			args: args{
				length:    100000,
				numDigits: 10000,
			},
			wantErr: true,
		},
		{
			name: "too many digits",
			args: args{
				length:    10,
				numDigits: 11,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GeneratePassword(tt.args.length, tt.args.numDigits)
			if (err != nil) != tt.wantErr {
				t.Errorf("GeneratePassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				letters, digits, others := passwordCharactersCounter(got)
				if letters != (tt.args.length-tt.args.numDigits) || digits != tt.args.numDigits || others != 0 {
					t.Errorf(`GeneratePassword() generated = %v.
It has %d letters, %d digits and %d other characters, but expected was %d letters, %d digits and 0 others`, got, letters, digits, others, tt.args.length-tt.args.numDigits, tt.args.numDigits)
				}
			}
		})
	}
}

func passwordCharactersCounter(password string) (int, int, int) {
	letters := 0
	digits := 0
	others := 0
	for _, r := range password {
		if unicode.IsLetter(r) {
			letters++
		} else if unicode.IsDigit(r) {
			digits++
		} else {
			others++
		}
	}
	return letters, digits, others
}
