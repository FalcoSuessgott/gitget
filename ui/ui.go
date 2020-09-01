package ui

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/AlecAivazis/survey/v2"
)

const (
	pageSize = 25
)

var (
	errEleNotFound = errors.New("element not found")
)

// PromptList prompts the user to choose one option from the specified list.
func PromptList(msg, def string, options []string) string {
	// https://github.com/AlecAivazis/survey/issues/101
	// https://github.com/AlecAivazis/survey/issues/101#issuecomment-420923209
	fmt.Printf("\x1b[?7l")

	result := ""
	prompt := &survey.Select{
		Message: msg,
		Options: options,
		Default: def,
	}

	err := survey.AskOne(prompt, &result, survey.WithPageSize(pageSize))

	if err != nil {
		fmt.Println("Exiting.")
		os.Exit(0)
	}

	defer fmt.Printf("\x1b[?7h")

	return result
}

func getIndexFromSlice(slice []string, element string) (int, error) {
	for i, e := range slice {
		if e == element {
			return i, nil
		}
	}

	return -1, errEleNotFound
}

// MultiSelect prompts the user for multiple options from a specified list.
func MultiSelect(msg string, elements []string) []int {
	fmt.Printf("\x1b[?7l")

	indexes := []int{}
	selectedFiles := []string{}

	prompt := &survey.MultiSelect{
		Message: msg,
		Options: elements,
	}

	err := survey.AskOne(prompt, &selectedFiles, survey.WithPageSize(pageSize))

	if err != nil {
		fmt.Println("Exiting.")
		os.Exit(0)
	}

	defer fmt.Printf("\x1b[?7h")

	for _, file := range selectedFiles {
		id, err := getIndexFromSlice(elements, file)
		if err != nil {
			log.Print(err)
		}

		indexes = append(indexes, id)
	}

	return indexes
}
