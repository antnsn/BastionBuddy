package azure

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/antnsn/BastionBuddy/internal/config"
	"github.com/antnsn/BastionBuddy/internal/utils"
)

// GetBastionDetails retrieves the Bastion host details either through
// manual input or by selecting from available hosts.
func GetBastionDetails(ctx context.Context, cred *azidentity.DefaultAzureCredential, subscriptionID string) (*config.BastionHost, error) {
	selectionMethod, err := utils.SelectWithMenu([]string{"auto-select-host", "manual-input(resource-id)"}, "How would you like to specify the Bastion host?")
	if err != nil {
		if strings.Contains(err.Error(), "cancelled by user") {
			fmt.Println("\nOperation cancelled by user")
			os.Exit(0)
		}
		return nil, fmt.Errorf("failed to get selection method: %v", err)
	}

	switch selectionMethod {
	case "auto-select-host":
		return GetBastionSelection(ctx, cred, subscriptionID)
	case "manual-input(resource-id)":
		return GetBastionManualInput()
	default:
		return nil, fmt.Errorf("invalid selection method: %s", selectionMethod)
	}
}

// GetBastionManualInput prompts the user to manually input Bastion host details.
func GetBastionManualInput() (*config.BastionHost, error) {
	name, err := utils.ReadInput("Enter Bastion host name")
	if err != nil {
		if strings.Contains(err.Error(), "cancelled by user") {
			fmt.Println("\nOperation cancelled by user")
			os.Exit(0)
		}
		return nil, fmt.Errorf("failed to read Bastion host name: %v", err)
	}
	if name == "" {
		return nil, fmt.Errorf("bastion host name cannot be empty")
	}

	resourceGroup, err := utils.ReadInput("Enter Bastion resource group")
	if err != nil {
		if strings.Contains(err.Error(), "cancelled by user") {
			fmt.Println("\nOperation cancelled by user")
			os.Exit(0)
		}
		return nil, fmt.Errorf("failed to read resource group: %v", err)
	}
	if resourceGroup == "" {
		return nil, fmt.Errorf("resource group cannot be empty")
	}

	return &config.BastionHost{
		Name:          name,
		ResourceGroup: resourceGroup,
	}, nil
}

// GetBastionSelection retrieves available Bastion hosts and lets the user select one.
func GetBastionSelection(ctx context.Context, cred *azidentity.DefaultAzureCredential, subscriptionID string) (*config.BastionHost, error) {
	debugPrintf("Fetching Bastion hosts...")

	// Create resources client
	resourceClient, err := armresources.NewClient(subscriptionID, cred, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create resources client: %v", err)
	}

	// List all Bastion hosts using resource client
	filter := "resourceType eq 'Microsoft.Network/bastionHosts'"
	debugPrintf("Using filter: %s\n", filter)
	pager := resourceClient.NewListPager(&armresources.ClientListOptions{
		Filter: &filter,
	})

	var resources []*armresources.GenericResourceExpanded
	pageCount := 0
	for pager.More() {
		pageCount++
		debugPrintf("Fetching page %d...\n", pageCount)
		page, err := pager.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to list Bastion hosts: %v", err)
		}

		debugPrintf("Found %d resources on page %d\n", len(page.Value), pageCount)
		resources = append(resources, page.Value...)
	}

	if len(resources) == 0 {
		return nil, fmt.Errorf("no Bastion hosts found in subscription")
	}

	var items []string
	bastionMap := make(map[string]*config.BastionHost)

	for _, resource := range resources {
		if resource.Name == nil {
			continue
		}

		// Extract resource group from ID
		parts := strings.Split(*resource.ID, "/")
		var resourceGroup string
		for i, part := range parts {
			if strings.EqualFold(part, "resourceGroups") && i+1 < len(parts) {
				resourceGroup = parts[i+1]
				break
			}
		}
		if resourceGroup == "" {
			debugPrintf("Could not extract resource group from ID: %s\n", *resource.ID)
			continue
		}

		item := fmt.Sprintf("%s (%s)", *resource.Name, resourceGroup)
		items = append(items, item)
		bastionMap[item] = &config.BastionHost{
			Name:           *resource.Name,
			ResourceGroup:  resourceGroup,
			SubscriptionID: subscriptionID,
		}
	}

	selected, err := utils.SelectWithMenu(items, "Select Bastion host")
	if err != nil {
		if strings.Contains(err.Error(), "cancelled by user") {
			fmt.Println("\nOperation cancelled by user")
			os.Exit(0)
		}
		return nil, fmt.Errorf("failed to select Bastion host: %v", err)
	}

	return bastionMap[selected], nil
}
