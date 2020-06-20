package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
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

	//defer os.RemoveAll(dir)

	r, err := git.PlainClone(dir, false, &git.CloneOptions{URL: url})
	
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

func (r *Repository) getFileContent(path string) ([]byte, error) {
	file, err := os.Open(path)

    if err != nil {
        return nil, err
	}
	
    defer file.Close()

  	b, err := ioutil.ReadAll(file)
  
	  if err != nil {
        return nil, err
	}

	return b, nil
}

func (r *Repository) listFiles() []string {

	files := []string{}

	_, err := r.Repo.Worktree()

	if err != nil {
		fmt.Println(err)
	}

	err = filepath.Walk(r.Path,
		func(dir string, info os.FileInfo, err error) error {

			if err != nil {
				return err
			}

			if info.IsDir() && info.Name() == git.GitDirName {
				return filepath.SkipDir
			}

			files = append(files, dir)

			return nil
		})

	if err != nil {
		fmt.Println(err)
	}

	return files
}

func (r *Repository) indexTree() []string{

	tree := []string{}

	for i, element := range strings.Split(r.Tree.Print(),"\n") {
		tree = append(tree, fmt.Sprintf("[%d] %s", i, element))
	}

	return tree
}

func (r *Repository) buildDirectoryTree() (gotree.Tree, error) {

	_, err := r.Repo.Worktree()

	if err != nil {
		return nil, err
	}

	i := 1
	shortPath := strings.Split(r.URL, "/")[4:]
	name := r.Path[strings.LastIndex(r.Path, "/")+1:]
	tree := gotree.New(strings.Join(shortPath, "/"))

	err = filepath.Walk(r.Path,
		func(dir string, info os.FileInfo, err error) error {

			if err != nil {
				return err
			}

			if info.IsDir() && info.Name() == git.GitDirName {
				return filepath.SkipDir
			}

			if name == r.Path[strings.LastIndex(r.Path, "/")+1:] && !info.IsDir() {
				tree.Add(info.Name())
			}

			if info.IsDir() && info.Name() != name {
				tmpTree := buildSubdirectoryTree(dir)
				i += len(tmpTree.Items())
				tree.AddTree(tmpTree)
				return filepath.SkipDir
			}
			return nil
		})

	if err != nil {
		return nil, err
	}

	return tree, nil
}


func buildSubdirectoryTree(dir string) gotree.Tree {

	dirName := getDirName(dir)
	tree := gotree.New(dirName)

	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {

		// if directory, step into and build tree
		if info.IsDir() && dirName != info.Name() {
			tree.AddTree(buildSubdirectoryTree(path))
			return filepath.SkipDir
		}

		// only add nodes to tree with the same depth
		if len(strings.Split(dir, "/"))+1 == len(strings.Split(path, "/")) &&
			info.Name() != dirName && !info.IsDir() {
			tree.Add(info.Name())
		}

		return nil
	})
	return tree
}

func getDirName(dir string) string {
	return dir[strings.LastIndex(dir, "/")+1:]
}