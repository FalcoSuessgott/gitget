package main 

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/disiqueira/gotree"
	"github.com/whilp/git-urls"
)

var (
    err error
    r Repository
)

//TODO: format code


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
    r.Tree, err = buildDirectoryTree(r.URL, r.Path)

    if err != nil {
        fmt.Println(err)
    }

    selectedFiles := multiSelect("Select files and directories to be imported", r.indexTree())
    selectedFilesIndexes := []int{}

    for _, file := range selectedFiles {
        index, _ := strconv.Atoi(GetStringInBetween(file, "[", "]"))
        selectedFilesIndexes = append(selectedFilesIndexes, index)
    }


   tree := gotree.New(".")

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
