package utils

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/manifoldco/promptui"
)

func SelectWithMenu(items []string, prompt string) (string, error) {
	if len(items) == 0 {
		return "", fmt.Errorf("no items to select from")
	}

	selectPrompt := promptui.Select{
		Label: prompt,
		Items: items,
		Size:  10, // Show 10 items at a time with scrolling
	}

	_, result, err := selectPrompt.Run()
	if err != nil {
		return "", fmt.Errorf("selection failed: %v", err)
	}

	return result, nil
}

func ReadInput(prompt string) (string, error) {
	inputPrompt := promptui.Prompt{
		Label:    prompt,
		Validate: func(input string) error {
			if strings.TrimSpace(input) == "" {
				return fmt.Errorf("input cannot be empty")
			}
			return nil
		},
	}

	result, err := inputPrompt.Run()
	if err != nil {
		return "", fmt.Errorf("input failed: %v", err)
	}

	return strings.TrimSpace(result), nil
}

func AzureCommand(args ...string) ([]byte, error) {
	cmd := exec.Command("az", args...)
	return cmd.Output()
}

func AzureInteractiveCommand(args ...string) error {
	cmd := exec.Command("az", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func AzureSetSubscription(subscriptionID string) error {
	return exec.Command("az", "account", "set", "--subscription", subscriptionID).Run()
}

func ExtractIDFromParentheses(input string) (string, error) {
	start := strings.Index(input, "(")
	end := strings.Index(input, ")")
	if start == -1 || end == -1 || start >= end {
		return "", fmt.Errorf("invalid format: %s", input)
	}
	return input[start+1 : end], nil
}

func CheckDependencies() error {
	dependencies := []string{"az"}
	var missing []string

	for _, dep := range dependencies {
		_, err := exec.LookPath(dep)
		if err != nil {
			missing = append(missing, dep)
		}
	}

	if len(missing) > 0 {
		return fmt.Errorf("missing required dependencies: %s", strings.Join(missing, ", "))
	}

	return nil
}
