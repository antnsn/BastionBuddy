package azure

import (
	"fmt"
	"os"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/antnsn/BastionBuddy/internal/utils"
)

// BastionHost represents an Azure Bastion host.
type BastionHost struct {
	Name          string
	ResourceGroup string
}

// GetBastionDetails retrieves the Bastion host details either through
// manual input or by selecting from available hosts.
func GetBastionDetails() (*BastionHost, error) {
	selectionMethod, err := utils.SelectWithMenu([]string{"select-host", "manual-input"}, "How would you like to specify the Bastion host?")
	if err != nil {
		if strings.Contains(err.Error(), "cancelled by user") {
			fmt.Println("\nOperation cancelled by user")
			os.Exit(0)
		}
		return nil, fmt.Errorf("failed to get selection method: %v", err)
	}

	switch selectionMethod {
	case "select-host":
		return GetBastionSelection()
	case "manual-input":
		return GetBastionManualInput()
	default:
		return nil, fmt.Errorf("invalid selection method: %s", selectionMethod)
	}
}

// GetBastionManualInput prompts the user to manually input Bastion host details.
func GetBastionManualInput() (*BastionHost, error) {
	name, err := utils.ReadInput("Enter Bastion host name")
	if err != nil {
		if strings.Contains(err.Error(), "cancelled by user") {
			fmt.Println("\nOperation cancelled by user")
			os.Exit(0)
		}
		return nil, fmt.Errorf("failed to read Bastion host name: %v", err)
	}
	if name == "" {
		return nil, fmt.Errorf("Bastion host name cannot be empty")
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

	return &BastionHost{
		Name:          name,
		ResourceGroup: resourceGroup,
	}, nil
}

// GetBastionSelection retrieves available Bastion hosts and lets the user select one.
func GetBastionSelection() (*BastionHost, error) {
	fmt.Println("Fetching Bastion hosts...")

	// Create resources client
	resourceClient, err := armresources.NewClient(subID, cred, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create resources client: %v", err)
	}

	// List all Bastion hosts using resource client
	filter := "resourceType eq 'Microsoft.Network/bastionHosts'"
	fmt.Printf("Using filter: %s\n", filter)
	pager := resourceClient.NewListPager(&armresources.ClientListOptions{
		Filter: &filter,
	})

	var items []string
	pageCount := 0
	for pager.More() {
		pageCount++
		fmt.Printf("Fetching page %d...\n", pageCount)
		page, err := pager.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to list Bastion hosts: %v", err)
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

			items = append(items, fmt.Sprintf("%s (%s)", *resource.Name, resourceGroup))
		}
	}

	if len(items) == 0 {
		// Try listing without filter to see what resources are available
		fmt.Println("No Bastion hosts found, listing all resources to debug...")
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
		return nil, fmt.Errorf("no Bastion hosts found in subscription %s", subID)
	}

	// Let user select a Bastion host
	selected, err := utils.SelectWithMenu(items, "Select Bastion host:")
	if err != nil {
		if strings.Contains(err.Error(), "cancelled by user") {
			fmt.Println("\nOperation cancelled by user")
			os.Exit(0)
		}
		return nil, fmt.Errorf("failed to select Bastion host: %v", err)
	}

	// Extract resource group from selection
	resourceGroup, err := utils.ExtractIDFromParentheses(selected)
	if err != nil {
		return nil, fmt.Errorf("failed to extract resource group: %v", err)
	}

	// Extract name from selection
	name := selected[:len(selected)-len(resourceGroup)-3]
	name = strings.TrimSpace(name)

	return &BastionHost{
		Name:          name,
		ResourceGroup: resourceGroup,
	}, nil
}
