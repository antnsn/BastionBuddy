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

// GetTargetResource prompts the user to select a target resource.
func GetTargetResource(ctx context.Context, cred *azidentity.DefaultAzureCredential, subscriptionID string) (*config.TargetResource, error) {
	selectionMethod, err := utils.SelectWithMenu([]string{"select-resource", "manual-input"}, "How would you like to specify the target resource?")
	if err != nil {
		if strings.Contains(err.Error(), "cancelled by user") {
			fmt.Println("\nOperation cancelled by user")
			os.Exit(0)
		}
		return nil, fmt.Errorf("failed to get selection method: %v", err)
	}

	switch selectionMethod {
	case "select-resource":
		return getResourceSelection(ctx, cred, subscriptionID)
	case "manual-input":
		return getResourceManualInput()
	default:
		return nil, fmt.Errorf("invalid selection method: %s", selectionMethod)
	}
}

// getResourceManualInput prompts the user to manually input resource details.
func getResourceManualInput() (*config.TargetResource, error) {
	resourceID, err := utils.ReadInput("Enter resource ID")
	if err != nil {
		if strings.Contains(err.Error(), "cancelled by user") {
			fmt.Println("\nOperation cancelled by user")
			os.Exit(0)
		}
		return nil, fmt.Errorf("failed to read resource ID: %v", err)
	}

	return &config.TargetResource{
		ID: resourceID,
	}, nil
}

// getResourceSelection retrieves available virtual machines and lets the user select one.
func getResourceSelection(ctx context.Context, cred *azidentity.DefaultAzureCredential, subscriptionID string) (*config.TargetResource, error) {
	debugPrintf("Fetching virtual machines from subscription: %s...\n", subscriptionID)

	client, err := armresources.NewClient(subscriptionID, cred, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create resources client: %v", err)
	}

	filter := "resourceType eq 'Microsoft.Compute/virtualMachines'"
	debugPrintf("Using filter: %s\n", filter)

	var resources []*armresources.GenericResourceExpanded
	pageNum := 1

	pager := client.NewListPager(&armresources.ClientListOptions{
		Filter: &filter,
	})

	for pager.More() {
		debugPrintf("Fetching page %d...\n", pageNum)
		page, err := pager.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get resources page: %v", err)
		}
		debugPrintf("Found %d resources on page %d\n", len(page.Value), pageNum)
		for _, res := range page.Value {
			debugPrintf("Found resource: Name=%s, Type=%s, ID=%s\n",
				*res.Name,
				*res.Type,
				*res.ID)
		}
		resources = append(resources, page.Value...)
		pageNum++
	}

	if len(resources) == 0 {
		return nil, fmt.Errorf("no virtual machines found in subscription")
	}

	var items []string
	resourceMap := make(map[string]*armresources.GenericResourceExpanded)

	for _, res := range resources {
		if res.Name == nil || res.ID == nil {
			continue
		}

		// Extract resource group from ID
		parts := strings.Split(*res.ID, "/")
		var resourceGroup string
		for i, part := range parts {
			if strings.EqualFold(part, "resourceGroups") && i+1 < len(parts) {
				resourceGroup = parts[i+1]
				break
			}
		}

		// Include name, resource group, and location in the display
		item := fmt.Sprintf("%s | Group: %s | Region: %s",
			*res.Name,
			resourceGroup,
			*res.Location)
		items = append(items, item)
		resourceMap[item] = res
	}

	selected, err := utils.SelectWithMenu(items, "Select virtual machine (type to filter)")
	if err != nil {
		if strings.Contains(err.Error(), "cancelled by user") {
			fmt.Println("\nOperation cancelled by user")
			os.Exit(0)
		}
		return nil, fmt.Errorf("failed to select resource: %v", err)
	}

	selectedResource := resourceMap[selected]
	return &config.TargetResource{
		ID:   *selectedResource.ID,
		Name: *selectedResource.Name,
		Type: "Microsoft.Compute/virtualMachines",
	}, nil
}
