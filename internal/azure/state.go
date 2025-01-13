package azure

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
)

var (
	// subID is the Azure subscription ID
	subID string

	// cred is the Azure credential
	cred *azidentity.DefaultAzureCredential

	// ctx is the context for Azure operations
	ctx context.Context
)

func init() {
	ctx = context.Background()
	var err error
	cred, err = azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		panic(err)
	}
}
