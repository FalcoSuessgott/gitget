package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/disiqueira/gotree"
)

var (
	err error
)

func main() {
	url := ""

	if len(os.Args) > 2 {
		fmt.Println("Too many arguments. Exiting.")
		os.Exit(1)
	}

	buf, _ := clipboard.ReadAll()

	if len(os.Args) == 1 && isGitURL(buf) {
		fmt.Println("Using git url from clipboard. ")

		url = buf
	} else {
		if len(os.Args) == 1 {
			fmt.Println("No git url passed. Exiting.")
			usage()
			os.Exit(1)
		}
		url = os.Args[1]
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
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage:\n\tgitget GIT_URL\n")
}
