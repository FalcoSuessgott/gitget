package main 

import (
	"fmt"
	"os"
	"strconv"
	"strings"
    "github.com/whilp/git-urls"
)

var (
    err error
    r Repository
)

//TODO: tree view for created files
//TODO: format code
//TODO: get uri from copy buffer

func main(){
    link  := os.Args[1]

    _, err := giturls.Parse(link)

    if err != nil {
        fmt.Println("Invalid git url")
        os.Exit(1)
    }

    fmt.Printf("Fetching %s\n\n", link)

    r.URL = link
    r.Repo, r.Path, err  = cloneRepo(r.URL)
    
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    r.Branches, err = r.getBranches()

    if err != nil {
        fmt.Println(err)
    }

    choosenBranch := r.Branches[0]

    if len(r.Branches) > 1 { 
        choosenBranch = promptList("Choose the branch to be checked out", "master", r.Branches)
    } 

    fmt.Println("\nChecking out the only branch: " + r.Branches[0])
    r.checkoutBranch(choosenBranch)

    r.Files = r.listFiles(r.Path)
    r.Tree, err = r.buildDirectoryTree()

    if err != nil {
        fmt.Println(err)
    }

    selectedFiles := multiSelect("Select files and directories to be imported", r.indexTree())
    selectedFilesIndexes := []int{}

    for _, file := range selectedFiles {
        index, _ := strconv.Atoi(GetStringInBetween(file, "[", "]"))
        selectedFilesIndexes = append(selectedFilesIndexes, index)
    }

    for _, i := range selectedFilesIndexes{
        path := strings.Split(r.Files[i], "/")
        name := path[len(path)-1]

        if isFile(r.Files[i]); err != nil {
            content, err := r.getFileContent(r.Files[i])

            if err != nil {
                fmt.Printf("Error while reading file: %s.(Err: %v)\n",r.Files[i], err)
            }

            createFile(name, content)
        } else {
            err := CopyDir(r.Files[i],name)
            
            if err != nil {
                fmt.Println(err)
            }  
        }
    }
}   
