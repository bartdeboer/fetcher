package main

import (
	"fmt"
	"io"
	"os"
	"path"
)

// CopyFile copies a single file from src to dest
func CopyFile(srcPath, destPath string) error {
	var err error
	var srcFile *os.File
	var destFile *os.File
	var fileInfo os.FileInfo

	if srcFile, err = os.Open(srcPath); err != nil {
		return err
	}
	defer srcFile.Close()

	if destFile, err = os.Create(destPath); err != nil {
		return err
	}
	defer destFile.Close()

	if _, err = io.Copy(destFile, srcFile); err != nil {
		return err
	}
	if fileInfo, err = os.Stat(srcPath); err != nil {
		return err
	}
	return os.Chmod(destPath, fileInfo.Mode())
}

// CopyDir recursively copies a dir from src to dest
func CopyDir(srcDir, destDir string) error {
	var err error
	var entries []os.DirEntry
	var dirInfo os.FileInfo

	if dirInfo, err = os.Stat(srcDir); err != nil {
		return err
	}

	if err = os.MkdirAll(destDir, dirInfo.Mode()); err != nil {
		return err
	}

	if entries, err = os.ReadDir(srcDir); err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := path.Join(srcDir, entry.Name())
		destPath := path.Join(destDir, entry.Name())

		if entry.IsDir() {
			if err = CopyDir(srcPath, destPath); err != nil {
				fmt.Println(err)
			}
		} else {
			if err = CopyFile(srcPath, destPath); err != nil {
				fmt.Println(err)
			}
		}
	}
	return nil
}
