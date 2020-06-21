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

func main(){
    link  := os.Args[1]

    _, err := giturls.Parse(link)

    if err != nil {
        fmt.Println("Invalid git url")
        os.Exit(1)
    }

    fmt.Printf("Fetching %s\n", link)

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
        choosenBranch = promptList("Branches", "master", r.Branches)
    } 

    fmt.Println("Checking out the only branch: " + r.Branches[0])
    r.checkoutBranch(choosenBranch)

    r.Files = r.listFiles()
    r.Tree, err = r.buildDirectoryTree()

    if err != nil {
        fmt.Println(err)
    }

    selectedFiles := multiSelect("Which files to import", r.indexTree())
    selectedFilesIndexes := []int{}

    for _, file := range selectedFiles {
        index, _ := strconv.Atoi(GetStringInBetween(file, "[", "]"))
        selectedFilesIndexes = append(selectedFilesIndexes, index)
    }

    for _, i := range selectedFilesIndexes{
        content, _ := r.getFileContent(r.Files[i])
        path := strings.Split(r.Files[i], "/")
        fileName := path[len(path)-1]
        createFile(fileName, content)
    }
}   

func GetStringInBetween(str string, start string, end string) (result string) {
    s := strings.Index(str, start)
    if s == -1 {
        return
    }
    s += len(start)
    e := strings.Index(str, end)
    if e == -1 {
        return
    }
    return str[s:e]
}