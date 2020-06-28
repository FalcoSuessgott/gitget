package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/disiqueira/gotree"
	"github.com/go-git/go-git/v5"
)

func (r *Repository) indexTree() []string {
	tree := []string{}

	for i, element := range strings.Split(r.Tree.Print(), "\n") {
		if i == 0 {
			tree = append(tree, element)
			continue
		}

		tree = append(tree, fmt.Sprintf("[%02d] %s", i, element))
	}

	return tree[:len(tree)-1]
}

func buildDirectoryTree(url, path string) (gotree.Tree, error) {
	i := 0
	name := path[strings.LastIndex(path, "/")+1:]
	tree := gotree.New(repoName(url))

	if filepath.Walk(path, func(dir string, info os.FileInfo, err error) error {

		if err != nil {
			return err
		}

		if info.IsDir() && info.Name() == git.GitDirName {
			return filepath.SkipDir
		}

		if name == path[strings.LastIndex(path, "/")+1:] && !info.IsDir() {
			tree.Add(info.Name())
		}

		if info.IsDir() && info.Name() != name {
			tmpTree := buildSubdirectoryTree(dir)
			i += len(tmpTree.Items())
			tree.AddTree(tmpTree)
			return filepath.SkipDir
		}
		return nil
	}); err != nil {
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
