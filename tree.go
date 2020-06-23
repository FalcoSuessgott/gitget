package main

import (
	"github.com/disiqueira/gotree"
	"fmt"
	"strings"
	"path/filepath"
	"os"
)

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

			if info.IsDir() && info.Name() == ".git" {
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
