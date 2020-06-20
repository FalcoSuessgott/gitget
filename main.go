package main 

import (
    "fmt"
    "strings"
    "strconv"
)

var (
    err error
    r Repository
)

func main(){
    link  := "https://github.com/FalcoSuessgott/dotfiles"

    fmt.Printf("Fetching %s\n", link)

    r.URL = link
    r.Repo, r.Path, err  = cloneRepo(r.URL)
    
    if err != nil {
        fmt.Println(err)
    }

    r.Branches, err = r.getBranches()

    fmt.Print(r.Path)

    if err != nil {
        fmt.Println(err)
    }

    choosenBranch := r.Branches[0]

    if len(r.Branches) > 1 { 
        choosenBranch = promptList("Branches", "master", r.Branches)
    } 

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
        fmt.Println(string(content))
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