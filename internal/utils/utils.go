package utils

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func SelectWithFzf(items []string, prompt string) (string, error) {
	cmd := exec.Command("fzf", "--prompt", prompt)
	cmd.Stderr = os.Stderr

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return "", err
	}

	go func() {
		defer stdin.Close()
		for _, item := range items {
			fmt.Fprintln(stdin, item)
		}
	}()

	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(output)), nil
}

func ReadInput(prompt string) (string, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(prompt)
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(input), nil
}

func AzureCommand(args ...string) ([]byte, error) {
	cmd := exec.Command("az", args...)
	return cmd.Output()
}

func AzureSetSubscription(subscriptionID string) error {
	return exec.Command("az", "account", "set", "--subscription", subscriptionID).Run()
}

func ExtractIDFromParentheses(input string) (string, error) {
	start := strings.LastIndex(input, "(") + 1
	end := strings.LastIndex(input, ")")
	if start <= 0 || end <= 0 || start >= end {
		return "", fmt.Errorf("invalid format: unable to extract ID from %s", input)
	}
	return input[start:end], nil
}

func CheckDependencies() error {
	deps := []string{"az", "fzf", "jq"}
	var missing []string
	for _, dep := range deps {
		if _, err := exec.LookPath(dep); err != nil {
			missing = append(missing, dep)
		}
	}
	if len(missing) > 0 {
		return fmt.Errorf("missing required dependencies: %s", strings.Join(missing, ", "))
	}
	return nil
}
