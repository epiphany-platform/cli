package az

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

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

func TestMain(m *testing.M) {
	setup()
	log.Println("Run tests")
	exitVal := m.Run()
	log.Println("Finish test")
	os.Exit(exitVal)
}

// initializes test with creation of key pair and checks if variables need to run tests are setup
func setup() {
	log.Println("Initialize test")
	tenantID = os.Getenv("TENANT_ID")
	if len(tenantID) == 0 {
		log.Fatalf("expected non-empty TENANT_ID environment variable")
	}

	subscriptionID = os.Getenv("SUBSCRIPTION_ID")
	if len(subscriptionID) == 0 {
		log.Fatalf("expected non-empty SUBSCRIPTION_ID environment variable")
	}

	spName = os.Getenv("SP_NAME")
	if len(spName) == 0 {
		log.Fatalf("expected non-empty SP_NAME environment variable")
	}

}

func TestShouldSuccessfullyGeneratePassword(t *testing.T) {
	pass := GenerateServicePrincipalPassword()
	if len(pass) < 32 {
		t.Error("Generated password too short.")
	}
}

func TestShouldSuccessfullyCreateServicePrincipal(t *testing.T) {

	// given

	// when
	pass := GenerateServicePrincipalPassword()

	sp, app := CreateServicePrincipal(pass, subscriptionID, tenantID, spName)

	env := getEnvironment(cloudTestName)

	graphAuthorizer := getGraphAuthorizer(env)

	spClient := graphrbac.NewServicePrincipalsClient(tenantID)
	spClient.Authorizer = graphAuthorizer

	spTest, err := spClient.Get(context.TODO(), *sp.ObjectID)
	if err != nil {
		t.Error(err)
	}

	spJSON, err := sp.MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(fmt.Sprint("App: ", string(spJSON)))

	if *spTest.ObjectID != *spTest.ObjectID {
		t.Error("Different object ID of Service Principal.")
	}

	if *sp.AppID != *spTest.AppID {
		t.Error("Different object ID of Service Principal.")
	}

	appClient := graphrbac.NewApplicationsClient(tenantID)
	appClient.Authorizer = graphAuthorizer

	appTest, err := appClient.Get(context.TODO(), *app.ObjectID)
	if err != nil {
		t.Fatal(err)
	}

	if *app.AppID != *appTest.AppID {
		t.Error("Different AppID of application.")
	}

	if *app.DisplayName != *appTest.DisplayName {
		t.Error("Different DisplayName of application.")
	}

	_, err = spClient.Delete(context.TODO(), *sp.ObjectID)
	if err != nil {
		t.Fatal(err)
	}

	_, err = appClient.Delete(context.TODO(), *app.ObjectID)
	if err != nil {
		t.Fatal(err)
	}
}

func TestShouldSuccessfullyCreateServicePrincipalCredentialsStruct(t *testing.T) {

	env := getEnvironment(cloudTestName)

	graphAuthorizer := getGraphAuthorizer(env)

	spClient := graphrbac.NewServicePrincipalsClient(tenantID)
	spClient.Authorizer = graphAuthorizer

	appClient := graphrbac.NewApplicationsClient(tenantID)
	appClient.Authorizer = graphAuthorizer

	pass := GenerateServicePrincipalPassword()
	sp, app := CreateServicePrincipal(pass, subscriptionID, tenantID, spName)

	creds := GenerateServicePrincipalCredentialsStruct(pass, tenantID, subscriptionID, sp, app)

	if creds.appID != *sp.AppID {
		t.Error("Different AppID in creds.")
	}

	if creds.password != pass {
		t.Error("Different password in creds.")
	}

	if creds.subscriptionID != subscriptionID {
		t.Error("Different subscriptionID in creds.")
	}

	if creds.tenant != tenantID {
		t.Error("Different tenantID in creds.")
	}

	_, err := spClient.Delete(context.TODO(), *sp.ObjectID)
	if err != nil {
		t.Fatal(err)
	}

	_, err = appClient.Delete(context.TODO(), *app.ObjectID)
	if err != nil {
		t.Fatal(err)
	}
}

func TestShouldSuccessfullyCreateServicePrincipalAuthJSON(t *testing.T) {

	env := getEnvironment(cloudTestName)

	graphAuthorizer := getGraphAuthorizer(env)

	spClient := graphrbac.NewServicePrincipalsClient(tenantID)
	spClient.Authorizer = graphAuthorizer

	appClient := graphrbac.NewApplicationsClient(tenantID)
	appClient.Authorizer = graphAuthorizer

	pass := GenerateServicePrincipalPassword()
	sp, app := CreateServicePrincipal(pass, subscriptionID, tenantID, spName)

	creds := GenerateServicePrincipalCredentialsStruct(pass, tenantID, subscriptionID, sp, app)

	GenerateServicePrincipalAuthJSONFromCredentialsStruct(creds)

	_, err := spClient.Delete(context.TODO(), *sp.ObjectID)
	if err != nil {
		t.Fatal(err)
	}

	_, err = appClient.Delete(context.TODO(), *app.ObjectID)
	if err != nil {
		t.Fatal(err)
	}
}

// getEnvironment returns Azure Environment based on cloudName
func getTestEnvironment(cloudName string) azure.Environment {
	env, err := azure.EnvironmentFromName(cloudName)
	if err != nil {
		errFailedToGetEnvironment(err)
	}
	return env
}

// getGraphAuthorizer return graph authorizer based on Azure Environment
func getTestGraphAuthorizer(env azure.Environment) autorest.Authorizer {
	graphAuthorizer, err := auth.NewAuthorizerFromCLIWithResource(env.GraphEndpoint)
	if err != nil {
		errFailedToGetGraphAuthrorizer(err)
	}
	return graphAuthorizer
}
