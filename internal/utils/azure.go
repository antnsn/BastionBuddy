package utils

import (
	"fmt"
	"os/exec"
	"strings"
)

// CheckAzureLogin checks if the user is logged into Azure CLI
func CheckAzureLogin() (bool, error) {
	cmd := exec.Command("az", "account", "show")
	output, err := cmd.CombinedOutput()
	if err != nil {
		// Check if az CLI is installed
		if _, lookErr := exec.LookPath("az"); lookErr != nil {
			return false, fmt.Errorf("Azure CLI (az) is not installed: %v", lookErr)
		}
		// Check if the error is due to not being logged in
		outputStr := string(output)
		if strings.Contains(outputStr, "az login") || strings.Contains(outputStr, "not logged in") {
			return false, nil
		}
		return false, fmt.Errorf("error checking Azure login status: %v\nOutput: %s", err, outputStr)
	}
	return true, nil
}

// LoginToAzure attempts to log in to Azure CLI
func LoginToAzure() error {
	fmt.Println("Not logged into Azure. Starting login process...")
	cmd := exec.Command("az", "login")
	cmd.Stdin = nil // Ensure the command opens a browser window
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to login to Azure: %v\nOutput: %s", err, string(output))
	}
	return nil
}
