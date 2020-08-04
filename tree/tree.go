package tree

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/disiqueira/gotree"
	"github.com/go-git/go-git/v5"
)

// BuildDirectoryTree builds the tree for the root directory.
func BuildDirectoryTree(url, path string) (gotree.Tree, error) {
	name := path[strings.LastIndex(path, "/")+1:]
	tree := gotree.New(url)
	err := filepath.Walk(path, func(dir string, info os.FileInfo, err error) error {
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
			tmpTree := BuildSubdirectoryTree(dir)
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

// BuildSubdirectoryTree builds the tree for any subdirectory.
func BuildSubdirectoryTree(dir string) gotree.Tree {
	dirName := getDirName(dir)
	tree := gotree.New(dirName)
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		// if directory, step into and build tree
		if info.IsDir() && dirName != info.Name() {
			tree.AddTree(BuildSubdirectoryTree(path))
			return filepath.SkipDir
		}

		// only add nodes to tree with the same depth
		if len(strings.Split(dir, "/"))+1 == len(strings.Split(path, "/")) &&
			info.Name() != dirName && !info.IsDir() {
			tree.Add(info.Name())
		}

		return nil
	})

	if err != nil {
		return nil
	}

	return tree
}

func getDirName(dir string) string {
	return dir[strings.LastIndex(dir, "/")+1:]
}

// NewTree returns a new gotree struct.
func NewTree(pwd string) gotree.Tree {
	return gotree.New(pwd)
}
