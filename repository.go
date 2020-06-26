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
	if giturls.Transports.Valid(url) {
		return true
	}
	return false
}

func getGitURI(url string) string{

	parsedURL, _ := giturls.Parse(url)

	return parsedURL.String()
}

func getBranches(repo *git.Repository) ([]string, error) {

	var branches []string

	bs, _ := remoteBranches(repo.Storer)

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

	_, err = giturls.Parse(url)

    if err != nil {
        fmt.Println("Invalid git url. Exiting.")
        os.Exit(1)
    }

    fmt.Printf("Fetching %s\n\n", url)
    
    repo, path, err  := cloneRepo(url)
    
    if err != nil {
        fmt.Println("Error while cloning. Exiting.")
        os.Exit(1)
    }

    branches, err := getBranches(repo)

    if err != nil {
        fmt.Println("Error while receiving Branches. Exiting.")
    }

	branch := ""
	if len(branches) == 1 {
		fmt.Println("\nChecking out the only branch: " + branches[0])
		branch = branches[0]
	} else {
		branch = promptList("Choose the branch to be checked out", "master", branches)
	}

	if checkoutBranch(repo, branch); err != nil {
		fmt.Println("Error while checking out branch " + branch + " .Exiting.")
	}

	files := listFiles(path)
    tree, err := buildDirectoryTree(url, path)

    if err != nil {
        fmt.Println(err)
	}
	
	return Repository{
		URL: url,
		Branch: branch,
		Branches: branches,
		Files: files,
		Path: path,
		Repo: repo,
		Tree: tree,
	}
}