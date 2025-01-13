// Package utils provides common utility functions for the BastionBuddy application,
// including menu handling, user input, and Azure CLI command execution.
package utils

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/manifoldco/promptui"
)

// SelectWithMenu presents an interactive menu to the user with the given items and prompt.
// It returns the selected item and any error that occurred.
func SelectWithMenu(items []string, prompt string) (string, error) {
	if len(items) == 0 {
		return "", errors.New("no items to select from")
	}

	if len(items) == 1 {
		return items[0], nil
	}

	searcher := func(input string, index int) bool {
		item := strings.ToLower(items[index])
		input = strings.ToLower(input)
		return strings.Contains(item, input)
	}

	selector := promptui.Select{
		Label:    prompt,
		Items:    items,
		Size:     10,
		Searcher: searcher,
	}

	_, result, err := selector.Run()
	if err == promptui.ErrInterrupt || err == promptui.ErrEOF {
		return "", errors.New("selection cancelled by user")
	}
	return result, err
}

// ReadInput prompts the user for input with the given prompt text.
// It returns the user's input and any error that occurred.
func ReadInput(prompt string) (string, error) {
	prompter := promptui.Prompt{
		Label: prompt,
	}

	return prompter.Run()
}

// AzureCommand executes an Azure CLI command with the given arguments.
// It returns the command output and any error that occurred.
func AzureCommand(args ...string) ([]byte, error) {
	cmd := exec.Command("az", args...)
	return cmd.Output()
}

// AzureInteractiveCommand executes an Azure CLI command that requires user interaction.
// It returns any error that occurred during command execution.
func AzureInteractiveCommand(args ...string) error {
	cmd := exec.Command("az", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// AzureSetSubscription sets the active Azure subscription.
// It returns any error that occurred during the operation.
func AzureSetSubscription(subscriptionID string) error {
	_, err := AzureCommand("account", "set", "--subscription", subscriptionID)
	return err
}

// ExtractIDFromParentheses extracts an ID from a string that contains it within parentheses.
// For example, "Resource Group (12345)" returns "12345".
func ExtractIDFromParentheses(input string) (string, error) {
	re := regexp.MustCompile(`\((.*?)\)`)
	matches := re.FindStringSubmatch(input)
	if len(matches) < 2 {
		return "", fmt.Errorf("no ID found in parentheses in input: %s", input)
	}
	return matches[1], nil
}

// CheckDependencies verifies that all required external dependencies are available.
// Currently checks for the Azure CLI (az) command.
func CheckDependencies() error {
	if _, err := exec.LookPath("az"); err != nil {
		return fmt.Errorf("Azure CLI (az) is not installed. Please install it from https://docs.microsoft.com/cli/azure/install-azure-cli")
	}
	return nil
}
