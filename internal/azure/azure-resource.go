package azure

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/antnsn/BastionBuddy/internal/utils"
)

// GetTargetResource prompts the user to select a target resource
func GetTargetResource() (string, error) {
	fmt.Println("Select how to specify the target resource:")
	selectionMethod, err := utils.SelectWithMenu([]string{"resource-id", "select-resource"}, "Selection method:")
	if err != nil {
		return "", fmt.Errorf("failed to select method: %v", err)
	}

	switch selectionMethod {
	case "resource-id":
		return getResourceManualInput()
	case "select-resource":
		return getResourceSelection()
	default:
		return "", fmt.Errorf("invalid selection method: %s", selectionMethod)
	}
}

func getResourceManualInput() (string, error) {
	resourceID, err := utils.ReadInput("Enter the target resource ID: ")
	if err != nil {
		return "", fmt.Errorf("failed to read resource ID: %v", err)
	}
	if resourceID == "" {
		return "", fmt.Errorf("resource ID cannot be empty")
	}
	return resourceID, nil
}

func getResourceSelection() (string, error) {
	// Select target subscription
	targetSubID, err := selectSubscription("target resource")
	if err != nil {
		return "", fmt.Errorf("failed to select target subscription: %v", err)
	}

	if err := utils.AzureSetSubscription(targetSubID); err != nil {
		return "", fmt.Errorf("failed to set target subscription: %v", err)
	}

	// Select resource group
	fmt.Println("Fetching resource groups...")
	rgOutput, err := utils.AzureCommand("group", "list", "--query", "[].name", "-o", "tsv")
	if err != nil {
		return "", fmt.Errorf("failed to list resource groups: %v", err)
	}

	rgList := strings.Split(strings.TrimSpace(string(rgOutput)), "\n")
	if len(rgList) == 0 || (len(rgList) == 1 && rgList[0] == "") {
		return "", fmt.Errorf("no resource groups found in subscription")
	}

	targetRG, err := utils.SelectWithMenu(rgList, "Select target resource group:")
	if err != nil {
		return "", fmt.Errorf("failed to select resource group: %v", err)
	}

	// Select target resource
	fmt.Printf("Fetching resources in %s...\n", targetRG)
	resourceOutput, err := utils.AzureCommand("resource", "list", "-g", targetRG,
		"--query", "[?type=='Microsoft.Compute/virtualMachines'].{Name:name, Id:id}", "-o", "json")
	if err != nil {
		return "", fmt.Errorf("failed to list resources: %v", err)
	}

	var resources []struct {
		Name string `json:"name"`
		ID   string `json:"id"`
	}
	if err := json.Unmarshal(resourceOutput, &resources); err != nil {
		return "", fmt.Errorf("failed to parse resources: %v", err)
	}

	if len(resources) == 0 {
		return "", fmt.Errorf("no virtual machines found in resource group %s", targetRG)
	}

	var options []string
	for _, res := range resources {
		options = append(options, fmt.Sprintf("%s (%s)", res.Name, res.ID))
	}

	selected, err := utils.SelectWithMenu(options, "Select target resource:")
	if err != nil {
		return "", fmt.Errorf("failed to select resource: %v", err)
	}

	return utils.ExtractIDFromParentheses(selected)
}
