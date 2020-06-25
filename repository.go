package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
    "github.com/whilp/git-urls"
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
	Branch string
	Files []string
	Tree gotree.Tree
}

func isValidGitURL(url string) bool {
	_, err := giturls.Parse(url)

    if err != nil {
		return false
	}	
	
	return true
}

func getBranches(repo *git.Repository) ([]string, error) {

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

func checkoutBranch(repo *git.Repository, branch string) error {

	w, err  := repo.Worktree()

	if err != nil {
		return err
	}

	err = w.Checkout(&git.CheckoutOptions{
		Branch: plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", branch)),
        Force: true,
	})
	
	return err
}

func NewRepository(url string) Repository{

	var r Repository
	_, err = giturls.Parse(url)

    if err != nil {
        fmt.Println("Invalid git url. Exiting.")
        os.Exit(1)
    }

    fmt.Printf("Fetching %s\n\n", os.Args[1])
    
    r.URL = url
    r.Repo, r.Path, err  = cloneRepo(url)
    
    if err != nil {
        fmt.Println("Error while cloning. Exiting.")
        os.Exit(1)
    }

    r.Branches, err = getBranches(r.Repo)

    if err != nil {
        fmt.Println("Error while receiving Branches. Exiting.")
    }

	if len(r.Branches) == 1 {
		fmt.Println("\nChecking out the only branch: " + r.Branches[0])
		r.Branch = r.Branches[0]
	} else {
		r.Branch = promptList("Choose the branch to be checked out", "master", r.Branches)
	}

	if checkoutBranch(r.Repo, r.Branch); err != nil {
		fmt.Println("Error while checking out branch " + r.Branch + " .Exiting.")
	}

	r.Files = listFiles(r.Path)
    r.Tree, err = buildDirectoryTree(r.URL, r.Path)

    if err != nil {
        fmt.Println(err)
	}
	
	return r
}