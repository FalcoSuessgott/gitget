package main 

import (
    "fmt"
)

func main(){
    repo  := "https://github.com/markbates/pkger"

    fmt.Printf("Fetching %s\n", repo)

    r, err := cloneRepo(repo)

    if err != nil {
        fmt.Println(err)
    }

    branches, err := getBranches(r)

    if err != nil {
        fmt.Println(err)
    }

   choosenBranch := promptList("Branches", "master", branches)

   checkoutBranch(r, choosenBranch)

   files, err := listFiles()

   if err != nil {
    fmt.Println(err)
   }

   promptList("Checkout files", "README.md", files)
}