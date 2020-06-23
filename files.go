package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"path/filepath"
)

const (
	gitDir = ".git"
)
func createFile(path string, content []byte) error {
	err := ioutil.WriteFile(path, content, os.FileMode(0644))

	if err != nil {
		return err
	}

	fmt.Printf("%s created.\n", path)

	return nil
}

func createDirectory(path string) error {
	err := os.MkdirAll(path, os.ModePerm)

	if err != nil {
		return err
	}

	fmt.Printf("%s created directory.\n", path)
	return nil
}

func fileOrDirExists(path string) (bool ,error){
	_, err := os.Stat(path)

    if err != nil {
        return false, err
	}

	return true, nil
}

func isFile(path string) (bool, error){
	file, err := os.Stat(path)

	if err != nil {
        return false, err
	}

    if file.IsDir() {
		return true, nil
	}

	return false, nil
}
func (r *Repository) getFileContent(path string) ([]byte, error) {

	file, _ := os.Open(path)
	
    defer file.Close()

  	b, err := ioutil.ReadAll(file)
  
	  if err != nil {
        return nil, err
	}

	return b, nil
}

func (r *Repository) listFiles(path string) []string {

	files := []string{}

	_, err := r.Repo.Worktree()

	if err != nil {
		fmt.Println(err)
	}

	err = filepath.Walk(path,
		func(dir string, info os.FileInfo, err error) error {

			if err != nil {
				return err
			}

			if info.IsDir() && info.Name() == gitDir {
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


func GetStringInBetween(str string, start string, end string) (result string) {
    s := strings.Index(str, start)
    if s == -1 {
        return
    }
    s += len(start)
    e := strings.Index(str, end)

    if e == -1 {
        return
    }

    return str[s:e]
}