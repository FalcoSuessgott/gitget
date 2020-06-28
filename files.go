package main

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const (
	gitDir = ".git"
)

func isFile(path string) bool {
	file, _ := os.Stat(path)
	return file.IsDir()
}

func listFiles(path string) []string {
	files := []string{}
	err := filepath.Walk(path, func(dir string, info os.FileInfo, err error) error {
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
		return nil
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return nil
	}

	defer func() {
		if e := out.Close(); e != nil {
			err = e
		}
	}()

	_, err = io.Copy(out, in)
	if err != nil {
		return nil
	}

	err = out.Sync()
	if err != nil {
		return nil
	}

	si, err := os.Stat(src)
	if err != nil {
		return nil
	}

	err = os.Chmod(dst, si.Mode())
	if err != nil {
		return nil
	}

	return nil
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
		return errors.New("source is not a directory")
	}

	_, err = os.Stat(dst)
	if err != nil && !os.IsNotExist(err) {
		return nil
	}

	if err == nil {
		return errors.New("destination already exists")
	}

	err = os.MkdirAll(dst, si.Mode())
	if err != nil {
		return nil
	}

	entries, err := ioutil.ReadDir(src)
	if err != nil {
		return nil
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			err = CopyDir(srcPath, dstPath)
			if err != nil {
				return nil
			}
		} else {
			// Skip symlinks.
			if entry.Mode()&os.ModeSymlink != 0 {
				continue
			}

			err = CopyFile(srcPath, dstPath)
			if err != nil {
				return nil
			}
		}
	}

	return nil
}
