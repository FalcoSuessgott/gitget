package cmd

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/FalcoSuessgott/gitget/fs"
	"github.com/FalcoSuessgott/gitget/repo"
	t "github.com/FalcoSuessgott/gitget/tree"
	"github.com/FalcoSuessgott/gitget/ui"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/atotto/clipboard"
)

var (
	rootCmd = &cobra.Command{
		Use:   "gitget",
		Short: "Browse interactively through branches, files and directories of a git repository and download them",
		Long:  "",
		Run: func(cmd *cobra.Command, args []string) {
			err := parseArgs(os.Args)
			if err != nil {
				fmt.Println(err)
				err := cmd.Help()
				if err != nil {
					log.Fatal(err)
				}
			}
		}}
	errNoGitURL   = errors.New("no git url passed")
	errToManyArgs = errors.New("too many arguments")
)

const (
	allowedArgsLength = 2
)

// Execute invokes the command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringP("remote", "r", "origin", "name of the remote")
	viper.SetDefault("remote", "origin")
}

func parseArgs(args []string) error {
	url := ""

	if len(args) > allowedArgsLength {
		return errToManyArgs
	}

	buf, _ := clipboard.ReadAll()

	if len(args) == 1 && repo.IsGitURL(buf) {
		fmt.Println("Using git url from clipboard.")

		url = buf
	} else {
		if len(args) == 1 {
			return errNoGitURL
		}

		url = args[1]
	}

	r := repo.NewRepository(url)
	gitTree := strings.Split(r.Tree.Print(), "\n")
	indexes := ui.MultiSelect("Select files and directories to be imported", gitTree)
	pwd, _ := os.Getwd()
	tree := t.NewTree(pwd)

	for _, i := range indexes {
		tree.AddTree(t.BuildSubdirectoryTree(r.Files[i]))
		path := strings.Split(r.Files[i], "/")
		name := path[len(path)-1]

		if !fs.IsFile(r.Files[i]) {
			err := fs.CopyFile(r.Files[i], name)
			if err != nil {
				fmt.Printf("Error while creating file: %s.(Err: %v)\n", r.Files[i], err)
			}
		} else {
			err := fs.CopyDir(r.Files[i], name)
			if err != nil {
				fmt.Printf("Error while creating directory: %s.(Err: %v)\n", r.Files[i], err)
			}
		}
	}

	os.RemoveAll(r.Path)
	fmt.Println("\nFetched the following files and directories: ")
	fmt.Println(tree.Print())

	return nil
}
