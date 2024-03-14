package fetcher

import (
	"fmt"
	"io"
	"os"
	"path"
)

// CopyFile copies a single file from src to dest
func copyFile(srcPath, destPath string) error {
	var err error
	var srcFile *os.File
	var destFile *os.File
	var fileInfo os.FileInfo

	if srcFile, err = os.Open(srcPath); err != nil {
		return fmt.Errorf("error copying file: %v", err)
	}
	defer srcFile.Close()

	if destFile, err = os.Create(destPath); err != nil {
		return fmt.Errorf("error copying file: %v", err)
	}
	defer destFile.Close()

	if _, err = io.Copy(destFile, srcFile); err != nil {
		return fmt.Errorf("error copying file: %v", err)
	}
	if fileInfo, err = os.Stat(srcPath); err != nil {
		return fmt.Errorf("error copying file: %v", err)
	}
	return os.Chmod(destPath, fileInfo.Mode())
}

// CopyDir recursively copies a dir from src to dest
func copyDir(srcDir, destDir string) error {
	var err error
	var entries []os.DirEntry
	var dirInfo os.FileInfo

	if dirInfo, err = os.Stat(srcDir); err != nil {
		return fmt.Errorf("error finding source directory %s: %v", srcDir, err)
	}

	if err = os.MkdirAll(destDir, dirInfo.Mode()); err != nil {
		return fmt.Errorf("error creating destination directory %s: %v", destDir, err)
	}

	if entries, err = os.ReadDir(srcDir); err != nil {
		return fmt.Errorf("error reading entries: %v", err)
	}

	for _, entry := range entries {
		srcPath := path.Join(srcDir, entry.Name())
		destPath := path.Join(destDir, entry.Name())

		if entry.IsDir() {
			if err = copyDir(srcPath, destPath); err != nil {
				fmt.Printf("Error copying directory: %v", err)
			}
		} else {
			if err = copyFile(srcPath, destPath); err != nil {
				fmt.Printf("Error copying file: %v", err)
			}
		}
	}
	return nil
}
