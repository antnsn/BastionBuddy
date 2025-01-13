package azure

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/antnsn/BastionBuddy/internal/utils"
)

// TargetResource represents an Azure resource that will be connected to.
type TargetResource struct {
	ID string
}

// GetTargetResource retrieves the target resource details either through
// manual input or by selecting from available resources.
func GetTargetResource(ctx context.Context, subID string, cred *azidentity.DefaultAzureCredential) (*TargetResource, error) {
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
		return getResourceSelection(ctx, subID, cred)
	case "manual-input":
		return getResourceManualInput()
	default:
		return nil, fmt.Errorf("invalid selection method: %s", selectionMethod)
	}
}

// getResourceManualInput prompts the user to manually input resource details.
func getResourceManualInput() (*TargetResource, error) {
	resourceID, err := utils.ReadInput("Enter resource ID")
	if err != nil {
		if strings.Contains(err.Error(), "cancelled by user") {
			fmt.Println("\nOperation cancelled by user")
			os.Exit(0)
		}
		return nil, fmt.Errorf("failed to read resource ID: %v", err)
	}
	if resourceID == "" {
		return nil, fmt.Errorf("resource ID cannot be empty")
	}

	return &TargetResource{
		ID: resourceID,
	}, nil
}

// getResourceSelection retrieves available virtual machines and lets the user select one.
func getResourceSelection(ctx context.Context, subID string, cred *azidentity.DefaultAzureCredential) (*TargetResource, error) {
	fmt.Println("Fetching virtual machines...")

	// Create resources client
	resourceClient, err := armresources.NewClient(subID, cred, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create resources client: %v", err)
	}

	// List all virtual machines using resource client
	filter := "resourceType eq 'Microsoft.Compute/virtualMachines'"
	fmt.Printf("Using filter: %s\n", filter)
	pager := resourceClient.NewListPager(&armresources.ClientListOptions{
		Filter: &filter,
	})

	var items []string
	resourceMap := make(map[string]string) // Map of "name (resourceGroup)" -> resourceID
	pageCount := 0
	for pager.More() {
		pageCount++
		fmt.Printf("Fetching page %d...\n", pageCount)
		page, err := pager.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to list virtual machines: %v", err)
		}

		fmt.Printf("Found %d resources on page %d\n", len(page.Value), pageCount)
		for _, resource := range page.Value {
			if resource.Name == nil {
				fmt.Println("Found resource with nil name, skipping")
				continue
			}

			fmt.Printf("Found resource: Name=%s, Type=%s, ID=%s\n", *resource.Name, *resource.Type, *resource.ID)

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
				fmt.Printf("Could not extract resource group from ID: %s\n", *resource.ID)
				continue
			}

			key := fmt.Sprintf("%s (%s)", *resource.Name, resourceGroup)
			items = append(items, key)
			resourceMap[key] = *resource.ID
		}
	}

	if len(items) == 0 {
		// Try listing without filter to see what resources are available
		fmt.Println("No virtual machines found, listing all resources to debug...")
		pager = resourceClient.NewListPager(nil)
		for pager.More() {
			page, err := pager.NextPage(ctx)
			if err != nil {
				return nil, fmt.Errorf("failed to list resources: %v", err)
			}
			for _, resource := range page.Value {
				if resource.Type != nil && resource.Name != nil {
					fmt.Printf("Available resource: Type=%s, Name=%s\n", *resource.Type, *resource.Name)
				}
			}
		}
		return nil, fmt.Errorf("no virtual machines found in subscription %s", subID)
	}

	// Let user select a virtual machine
	selected, err := utils.SelectWithMenu(items, "Select virtual machine:")
	if err != nil {
		if strings.Contains(err.Error(), "cancelled by user") {
			fmt.Println("\nOperation cancelled by user")
			os.Exit(0)
		}
		return nil, fmt.Errorf("failed to select virtual machine: %v", err)
	}

	// Look up the resource ID from our map
	resourceID, ok := resourceMap[selected]
	if !ok {
		return nil, fmt.Errorf("could not find resource ID for selection: %s", selected)
	}

	return &TargetResource{
		ID: resourceID,
	}, nil
}
