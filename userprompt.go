package main

import (
	"fmt"
	"os"
	"github.com/AlecAivazis/survey/v2"
)

func promptList(msg, def string, options []string) string {
	// https://github.com/AlecAivazis/survey/issues/101
	// https://github.com/AlecAivazis/survey/issues/101#issuecomment-420923209
	fmt.Printf("\x1b[?7l")

	result := ""
	prompt := &survey.Select{
    	Message: msg,
		Options: options,
		Default: def,
	}

	err := survey.AskOne(prompt, &result,  survey.WithPageSize(15))

	if err != nil {
		fmt.Println("Exiting.")
		os.Exit(0)
	}

	defer fmt.Printf("\x1b[?7h")

	return result
}

func multiSelect(msg string, elements []string) []string{
	
	fmt.Printf("\x1b[?7l")
	selectedFiles := []string{}

	prompt := &survey.MultiSelect{
		Message: msg,
		Options: elements,
	}
	
	err := survey.AskOne(prompt, &selectedFiles, survey.WithPageSize(15))

	if err != nil {
		fmt.Println("Exiting.")
		os.Exit(0)
	}

	defer fmt.Printf("\x1b[?7h")

	return selectedFiles
}