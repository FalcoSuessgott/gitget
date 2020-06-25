package main

import (
	"fmt"
	"io"
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

func isFile(path string) bool {
	file, _ := os.Stat(path)

    if file.IsDir() {
		return true
	}

	return false
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

func listFiles(path string) []string {

	files := []string{}

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

//https://gist.github.com/r0l1/92462b38df26839a3ca324697c8cba04
func CopyFile(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		if e := out.Close(); e != nil {
			err = e
		}
	}()

	_, err = io.Copy(out, in)
	if err != nil {
		return
	}

	err = out.Sync()
	if err != nil {
		return
	}

	si, err := os.Stat(src)
	if err != nil {
		return
	}
	err = os.Chmod(dst, si.Mode())
	if err != nil {
		return
	}

	return
}

//https://gist.github.com/r0l1/92462b38df26839a3ca324697c8cba04
func CopyDir(src string, dst string) (err error) {


	src = filepath.Clean(src)
	dst = filepath.Clean(dst)

	si, err := os.Stat(src)
	if err != nil {
		return err
	}
	if !si.IsDir() {
		return fmt.Errorf("source is not a directory")
	}

	_, err = os.Stat(dst)
	if err != nil && !os.IsNotExist(err) {
		return
	}
	if err == nil {
		return fmt.Errorf("destination already exists")
	}

	err = os.MkdirAll(dst, si.Mode())
	if err != nil {
		return
	}

	entries, err := ioutil.ReadDir(src)
	if err != nil {
		return
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			err = CopyDir(srcPath, dstPath)
			if err != nil {
				return
			}
		} else {
			// Skip symlinks.
			if entry.Mode()&os.ModeSymlink != 0 {
				continue
			}

			err = CopyFile(srcPath, dstPath)
			if err != nil {
				return
			}
		}
	}

	return
}