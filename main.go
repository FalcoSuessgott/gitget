package main 

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"github.com/disiqueira/gotree"
    "github.com/atotto/clipboard"
)

var (
    url string
    err error
    r Repository
)

//TODO: format code
//SSH clone

func main(){
    buf, _ := clipboard.ReadAll()

    if len(os.Args) == 1 && isValidGitURL(buf)  {
        fmt.Printf("Using git url from clipboard. ", strings.TrimSpace(buf))
        url = buf
    }

    if len(os.Args) > 2 {
        fmt.Println("Too many arguments. Exiting.")
        os.Exit(1)
    }

    url = os.Args[1]
    r := NewRepository(url)

    selectedFiles := multiSelect("Select files and directories to be imported", r.indexTree())
    selectedFilesIndexes := []int{}

    for _, file := range selectedFiles {
        index, _ := strconv.Atoi(GetStringInBetween(file, "[", "]"))
        selectedFilesIndexes = append(selectedFilesIndexes, index)
    }


    pwd, _ := os.Getwd()
    tree := gotree.New(pwd)

    for _, i := range selectedFilesIndexes{
        tree.AddTree(buildSubdirectoryTree(r.Files[i]))
        path := strings.Split(r.Files[i], "/")
        name := path[len(path)-1]

        if !isFile(r.Files[i]) {
            if CopyFile(r.Files[i], name); err != nil {
                fmt.Printf("Error while creating file: %s.(Err: %v)\n",r.Files[i], err)
            }
        } else {
            if CopyDir(r.Files[i],name); err != nil {
                fmt.Printf("Error while creating directory: %s.(Err: %v)\n",r.Files[i], err)
            }
        }
    }

    fmt.Println("\nFetched the following files and directories: ")
    fmt.Print(tree.Print())
}   
