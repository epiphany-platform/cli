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

func TestShouldSuccessfullyCreateServicePrincipal(t *testing.T) {

	// given

	// when
	creds := CreateSP(subscriptionID, tenantID, spName)

	t.Log(fmt.Sprintf("\n===========\nCREDENCIALS\n%+v\n===========\n", creds))

	env := getEnvironment(cloudTestName)

	graphAuthorizer := getGraphAuthorizer(env)

	spClient := graphrbac.NewServicePrincipalsClient(tenantID)
	spClient.Authorizer = graphAuthorizer

	sp, err := spClient.Get(context.TODO(), creds.appID)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(sp)
	// then
	// if !matched {
	// 	t.Error("Expected to find expression matching:\n", expectedFileContentRegexp, "\nbut found:\n", fileContent)
	// }

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
	} else {
		info("Got Azure Graph authorizer successfully .")
	}
	return graphAuthorizer
}
