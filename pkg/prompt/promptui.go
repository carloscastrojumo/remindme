package prompt

import (
	"fmt"
	"strings"

	"github.com/manifoldco/promptui"
)

func PromptForString(label string) string {
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

func PromptForStringArray(label string) []string {
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
