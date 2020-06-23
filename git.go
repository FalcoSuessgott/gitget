package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/disiqueira/gotree"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/storer"
)

type Repository struct {
	URL string
    Repo *git.Repository
	Path string
	Branches []string
	Files []string
	Tree gotree.Tree
}

func (r *Repository) getBranches() ([]string, error) {

	var branches []string

	bs, _ := remoteBranches(r.Repo.Storer)

	bs.ForEach(func(b *plumbing.Reference) error {
		name := strings.Split(b.Name().String(), "/")[3:]
		branches = append(branches, strings.Join(name, ""))
		return nil
	})

	return branches, nil
}

func cloneRepo(url string) (*git.Repository, string, error) {
	dir, err := ioutil.TempDir("", "tmp-dir")

	if err != nil {
		return nil, dir, err
	}


	r, err := git.PlainClone(dir, false, &git.CloneOptions{
		URL: url,
		Tags: git.NoTags,
		Progress: os.Stdout,
	})
	
	if err != nil {
		return nil, dir, err
	}

	return r, dir, nil
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

func (r *Repository) checkoutBranch(branch string) error {

	w, err  := r.Repo.Worktree()

	if err != nil {
		return err
	}

	err = w.Checkout(&git.CheckoutOptions{
		Branch: plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", branch)),
        Force: true,
	})
	
	return err
}

