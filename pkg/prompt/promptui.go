package prompt

import (
	"errors"
	"fmt"
	"strings"

	"github.com/manifoldco/promptui"
)

// ForString prompts the user for a string
func ForString(label string) (string, error) {
	return run(promptui.Prompt{Label: label})
}

// ForRequiredString prompts the user until a non-blank string is entered
func ForRequiredString(label string) (string, error) {
	return run(promptui.Prompt{Label: label, Validate: requireValue})
}

// ForRequiredStringArray prompts the user until a comma separated list with at least one value is entered
func ForRequiredStringArray(label string) ([]string, error) {
	result, err := run(promptui.Prompt{Label: label, Validate: requireList})
	if err != nil {
		return nil, err
	}
	return splitList(result), nil
}

func run(p promptui.Prompt) (string, error) {
	result, err := p.Run()
	if err != nil {
		return "", fmt.Errorf("%v prompt: %w", p.Label, err)
	}
	return result, nil
}

func requireValue(input string) error {
	if strings.TrimSpace(input) == "" {
		return errors.New("a value is required")
	}
	return nil
}

func requireList(input string) error {
	if len(splitList(input)) == 0 {
		return errors.New("at least one value is required")
	}
	return nil
}

func splitList(input string) []string {
	var values []string
	for _, value := range strings.Split(input, ",") {
		if value = strings.TrimSpace(value); value != "" {
			values = append(values, value)
		}
	}
	return values
}
