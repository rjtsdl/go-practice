package main

import (
	"log"

	"github.com/Azure/azure-sdk-for-go/arm/keyvault"
	"github.com/Azure/go-autorest/autorest/azure"
)

func main() {
	vaultsClient := keyvault.NewVaultsClient("c1089427-83d3-4286-9f35-5af546a6eb67")
	spt, err := getSpt("72f988bf-86f1-41af-91ab-2d7cd011db47",
		"https://jiren1025swarm.westus2.cloudapp.azure.com",
		"resource",
		"mTxm4Arntoy97U3FWVKmDpwO8znkG5BP")
	vaultsClient.Authorizer = spt
	vault, err := vaultsClient.Get("ACSKeyVault-Int", "ACSKeyVault-Int")
	if err != nil {
		log.Fatal(err)
	}
	log.Print(vault)
}

func getSpt(tenantID, clientID, resource, clientSec string) (*azure.ServicePrincipalToken, error) {
	oauthConfig, err := azure.PublicCloud.OAuthConfigForTenant(tenantID)
	if err != nil {
		return nil, err
	}
	spt, err := azure.NewServicePrincipalTokenWithSecret(
		*oauthConfig,
		clientID,
		resource,
		&azure.ServicePrincipalTokenSecret{
			ClientSecret: clientSec,
		},
		nil)
	return spt, err
}
