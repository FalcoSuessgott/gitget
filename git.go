package main

import (
	"os"
	"fmt"
	"strings"
	"path/filepath"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/go-git/go-git/v5/plumbing/storer"
)

func getBranches(r *git.Repository) ([]string, error) {

	var branches []string

	bs, _ := remoteBranches(r.Storer)
	bs.ForEach(func(b *plumbing.Reference) error {
		name := strings.Split(b.Name().String(), "/")[3:]
		branches = append(branches, strings.Join(name, ""))
		return nil
	})

	return branches, nil
}

func cloneRepo(url string) (*git.Repository, error){
	
	r, err := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{URL: url})
	
	if err != nil {
		return nil, err
	}

	return r,nil 
}

func remoteBranches(s storer.ReferenceStorer) (storer.ReferenceIter, error) {
	refs, err := s.IterReferences()

	if err != nil {
		return nil, err
	}

	return storer.NewReferenceFilteredIter(func(ref *plumbing.Reference) bool {
		return ref.Name().IsRemote()
	}, refs), nil
}

func checkoutBranch(r *git.Repository, branch string) error {

	w, err  := r.Worktree()

	if err != nil {
		return err
	}

	err = w.Checkout(&git.CheckoutOptions{
		Branch: plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", branch)),
        Force: true,
	})
	
	return err
}

func listFiles(r *git.Repository, dir string) ([]string, error){

	files := []string{}
	err := filepath.Walk(dir,
		func(dir string, info os.FileInfo, err error) error {

			if err != nil {
				return err
			}

			if info.IsDir() && info.Name() == git.GitDirName {
				return filepath.SkipDir
			}

			files = append(files, info.Name())
			return nil
		})

	if err != nil {
		return nil, err
	}	

	return files, nil
}