package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/disiqueira/gotree"
	"github.com/urfave/cli/v2"
)

var (
	err error
)

func main() {
	app := &cli.App{
		Name:  "gitget",
		Usage: "Browse interactively through branches, files and directories of a git repository and download them",
		Action: func(c *cli.Context) error {
			err := parseArgs(os.Args)
			if err != nil {
				fmt.Println(err)
				cli.ShowAppHelp(c)
			}

			return nil
		},
		UsageText: "gitget GIT_URL (e.g gitget https://github.com/golang/example)",
		Version:   "1.1.0",
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func parseArgs(args []string) error {
	url := ""

	if len(args) > 2 {
		return fmt.Errorf("too many arguments")
	}

	buf, _ := clipboard.ReadAll()

	if len(args) == 1 && isGitURL(buf) {
		fmt.Println("Using git url from clipboard.")

		url = buf
	} else {
		if len(args) == 1 {
			return fmt.Errorf("no git url passed")
		}

		url = args[1]
	}

	r := NewRepository(url)
	selectedFiles := multiSelect("Select files and directories to be imported", r.indexTree())
	selectedFilesIndexes := []int{}

	for _, file := range selectedFiles {
		index, _ := strconv.Atoi(GetStringInBetween(file, "[", "]"))
		selectedFilesIndexes = append(selectedFilesIndexes, index)
	}

	pwd, _ := os.Getwd()
	tree := gotree.New(pwd)

	for _, i := range selectedFilesIndexes {
		tree.AddTree(buildSubdirectoryTree(r.Files[i]))
		path := strings.Split(r.Files[i], "/")
		name := path[len(path)-1]

		if !isFile(r.Files[i]) {
			if CopyFile(r.Files[i], name); err != nil {
				fmt.Printf("Error while creating file: %s.(Err: %v)\n", r.Files[i], err)
			}
		} else {
			if CopyDir(r.Files[i], name); err != nil {
				fmt.Printf("Error while creating directory: %s.(Err: %v)\n", r.Files[i], err)
			}
		}
	}

	os.RemoveAll(r.Path)
	fmt.Println("\nFetched the following files and directories: ")
	fmt.Println(tree.Print())

	return nil
}
