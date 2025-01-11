package azure

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/yourusername/azbastion/internal/utils"
)

type Resource struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

func GetTargetResource() (string, error) {
	fmt.Println("Select how to specify the target resource:")
	selectionMethod, err := utils.SelectWithFzf([]string{"resource-id", "select-resource"}, "Selection method: ")
	if err != nil {
		return "", err
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
		return "", err
	}
	return resourceID, nil
}

func getResourceSelection() (string, error) {
	targetSubID, err := utils.SelectWithFzf([]string{"subscription1", "subscription2"}, "Select a target subscription:")
	if err != nil {
		return "", err
	}

	if err := utils.AzureSetSubscription(targetSubID); err != nil {
		return "", err
	}

	// Select resource group
	fmt.Println("Fetching resource groups...")
	rgOutput, err := utils.AzureCommand("group", "list", "--query", "[].name", "-o", "tsv")
	if err != nil {
		return "", err
	}

	rgList := strings.Split(strings.TrimSpace(string(rgOutput)), "\n")
	targetRG, err := utils.SelectWithFzf(rgList, "Select target resource group: ")
	if err != nil {
		return "", err
	}

	// Select resource
	fmt.Printf("Fetching resources in %s...\n", targetRG)
	resourceOutput, err := utils.AzureCommand("resource", "list", "-g", targetRG,
		"--query", "[?type=='Microsoft.Compute/virtualMachines'].{Name:name, Id:id}", "-o", "json")
	if err != nil {
		return "", err
	}

	var resources []Resource
	if err := json.Unmarshal(resourceOutput, &resources); err != nil {
		return "", err
	}

	var options []string
	for _, res := range resources {
		options = append(options, fmt.Sprintf("%s (%s)", res.Name, res.ID))
	}

	selected, err := utils.SelectWithFzf(options, "Select target resource: ")
	if err != nil {
		return "", err
	}

	return utils.ExtractIDFromParentheses(selected)
}
