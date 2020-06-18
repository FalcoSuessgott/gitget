package main 

import (
    "fmt"
)

func main(){
    repo  := "https://github.com/markbates/pkger.git"

    fmt.Printf("Fetching %s\n", repo)

    r, err := cloneRepo(repo)

    if err != nil {
        fmt.Println(err)
    }

    branches, err := getBranches(r)

    if err != nil {
        fmt.Println(err)
    }

    fmt.Print(branches)
    promptList("Branches", "master", branches)
}