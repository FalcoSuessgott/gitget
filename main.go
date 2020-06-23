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
//TODO: recursion for directories
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
    
        if isFile(r.Files[i]); err != nil {
            content, err := r.getFileContent(r.Files[i])
            path := strings.Split(r.Files[i], "/")
            fileName := path[len(path)-1]

            if err != nil {
                fmt.Printf("Error while reading file: %s.(Err: %v)\n",r.Files[i], err)
            }

            createFile(fileName, content)
        } else {
            downloadDirectory(r.Files[i])   
        }
    }
}   

func downloadDirectory(dir string) {
    path := strings.Split(dir, "/")
    dirName := path[len(path)-1]
      
    createDirectory(dirName) 
    subFiles := r.listFiles(dir)

    for j, f := range subFiles {
        subPath := strings.Split(f, "/")
        subFileName := subPath[len(subPath)-1]
        
        if j == 0 {
            continue
        }
    
        if isFile(f); err != nil {
            content, err := r.getFileContent(subFiles[j])
        
            if err != nil {
                fmt.Println(err.Error() + " while creating ")
            }
        
            createFile(dirName + "/" + subFileName, content)
        } else {
            downloadDirectory(f)
        }
    }
}