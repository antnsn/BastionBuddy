package azure

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/antnsn/BastionBuddy/internal/utils"
)

// GetBastionDetails prompts the user to select or input bastion details
func GetBastionDetails(subscriptionID string) (string, string, error) {
	fmt.Println("Select how to specify the Bastion host:")
	selectionMethod, err := utils.SelectWithMenu([]string{"manual-input", "select-bastion"}, "Selection method:")
	if err != nil {
		return "", "", fmt.Errorf("failed to select method: %v", err)
	}

	switch selectionMethod {
	case "manual-input":
		return getBastionManualInput()
	case "select-bastion":
		return getBastionSelection(subscriptionID)
	default:
		return "", "", fmt.Errorf("invalid selection method: %s", selectionMethod)
	}
}

func getBastionManualInput() (string, string, error) {
	bastionName, err := utils.ReadInput("Enter Bastion host name: ")
	if err != nil {
		return "", "", fmt.Errorf("failed to read bastion name: %v", err)
	}
	if bastionName == "" {
		return "", "", fmt.Errorf("bastion name cannot be empty")
	}

	bastionRG, err := utils.ReadInput("Enter Bastion resource group: ")
	if err != nil {
		return "", "", fmt.Errorf("failed to read resource group: %v", err)
	}
	if bastionRG == "" {
		return "", "", fmt.Errorf("resource group cannot be empty")
	}

	return bastionName, bastionRG, nil
}

func getBastionSelection(subscriptionID string) (string, string, error) {
	fmt.Println("Fetching Bastion hosts...")
	data, err := cacheInstance.Get(fmt.Sprintf("bastions_%s.json", subscriptionID), func() ([]byte, error) {
		return utils.AzureCommand("network", "bastion", "list")
	})
	if err != nil {
		return "", "", fmt.Errorf("failed to list bastion hosts: %v", err)
	}

	var bastions []BastionHost
	if err := json.Unmarshal(data, &bastions); err != nil {
		return "", "", fmt.Errorf("failed to parse bastion hosts: %v", err)
	}

	if len(bastions) == 0 {
		return "", "", fmt.Errorf("no Bastion hosts found or insufficient permissions")
	}

	var options []string
	for _, bastion := range bastions {
		options = append(options, fmt.Sprintf("%s (Resource Group: %s)",
			bastion.Name, bastion.ResourceGroup))
	}

	selected, err := utils.SelectWithMenu(options, "Select Bastion host:")
	if err != nil {
		return "", "", fmt.Errorf("failed to select bastion host: %v", err)
	}

	parts := strings.Split(selected, " (Resource Group: ")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid bastion selection format")
	}

	return parts[0], strings.TrimRight(parts[1], ")"), nil
}
