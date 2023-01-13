package prompt

import (
	"fmt"
	"strings"

	"github.com/manifoldco/promptui"
)

// ForString prompts the user for a string
func ForString(label string) string {
	prompt := promptui.Prompt{
		Label: label,
	}

	result, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return ""
	}

	return result
}

// ForStringArray prompts the user for a string array
func ForStringArray(label string) []string {
	prompt := promptui.Prompt{
		Label: label,
	}

	result, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return []string{}
	}

	return strings.Split(result, ",")
}
