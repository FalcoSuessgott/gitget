package main

import (
	"io/ioutil"
	"os"
	"strings"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

func getBranches(r *git.Repository) ([]string, error) {

	var branches []string

	branchList, err := r.Branches()
	
	if err != nil {
		return nil, err
	}

	branchList.ForEach(func(b *plumbing.Reference) error {
		name := strings.Split(b.Name().String(), "/")[2:]
		branches = append(branches, strings.Join(name, ""))
		return nil
	})

	return branches, nil
}

func cloneRepo(url string) (*git.Repository, error){
	dir, err := ioutil.TempDir("", "tmp-dir")

	if err != nil {
		return nil,err
	}

	defer os.RemoveAll(dir)
	
	r, err := git.PlainClone(dir, false, &git.CloneOptions{URL: url})
	
	if err != nil {
		return nil, err
	}
	
	return r, nil
}
