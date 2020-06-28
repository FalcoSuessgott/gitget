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
	ssh2 "github.com/go-git/go-git/v5/plumbing/transport/ssh"
	giturls "github.com/whilp/git-urls"
	"golang.org/x/crypto/ssh"
)

type Repository struct {
	URL      string
	Repo     *git.Repository
	Path     string
	Branches []string
	Branch   string
	Files    []string
	Tree     gotree.Tree
}

func isGitURL(rawURL string) bool {
	parsedURL, err := giturls.Parse(rawURL)
	if err == nil && parsedURL.IsAbs() && parsedURL.Hostname() != "" {
		return true
	}

	return false
}

func isSSHURL(rawURL string) bool {
	url, err := giturls.Parse(rawURL)
	return err == nil && (url.Scheme == "git" || url.Scheme == "ssh")
}

func repoName(repoURL string) string {
	u, _ := giturls.Parse(repoURL)
	return u.Path[1:]
}

func getBranches(repo *git.Repository) ([]string, error) {
	var branches []string

	bs, _ := remoteBranches(repo.Storer)

	err := bs.ForEach(func(b *plumbing.Reference) error {
		name := strings.Split(b.Name().String(), "/")[3:]
		branches = append(branches, strings.Join(name, ""))
		return nil
	})

	if err != nil {
		return nil, err
	}

	return branches, nil
}

func cloneRepo(url string) (*git.Repository, string, error) {
	var r *git.Repository

	dir, err := ioutil.TempDir("", "tmp-dir")
	if err != nil {
		return nil, dir, err
	}

	if isSSHURL(url) {
		s := fmt.Sprintf("%s/.ssh/id_rsa", os.Getenv("HOME"))
		sshKey, _ := ioutil.ReadFile(s)
		signer, _ := ssh.ParsePrivateKey(sshKey)
		auth := &ssh2.PublicKeys{User: "git", Signer: signer}

		r, _ = git.PlainClone(dir, false, &git.CloneOptions{
			URL:      url,
			Progress: os.Stdout,
			Tags:     git.NoTags,
			Auth:     auth,
		})
	} else {
		r, err = git.PlainClone(dir, false, &git.CloneOptions{
			URL:      url,
			Tags:     git.NoTags,
			Progress: os.Stdout,
		})
	}

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
	w, err := repo.Worktree()

	if err != nil {
		return err
	}

	err = w.Checkout(&git.CheckoutOptions{
		Branch: plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", branch)),
		Force:  true,
	})

	return err
}

func NewRepository(url string) Repository {
	if !isGitURL(url) {
		fmt.Println("Invalid git url. Exiting.")
		os.Exit(1)
	}

	fmt.Printf("Fetching %s\n\n", url)

	repo, path, err := cloneRepo(url)

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
		URL:      url,
		Branch:   branch,
		Branches: branches,
		Files:    files,
		Path:     path,
		Repo:     repo,
		Tree:     tree,
	}
}
