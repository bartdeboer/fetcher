package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func installFromArchive(archiveFilename string) error {
	// Extract the basename of the archive
	base := strings.TrimSuffix(filepath.Base(archiveFilename), filepath.Ext(archiveFilename))

	// Create a temporary directory
	tempDir := filepath.Join(os.TempDir(), base)
	if err := os.Mkdir(tempDir, 0755); err != nil {
		return fmt.Errorf("Error installing archive: %v", err)
	}
	defer os.RemoveAll(tempDir) // Ensure cleanup

	// Extract the archive into the temporary directory
	if err := extractArchive(archiveFilename, tempDir); err != nil {
		return fmt.Errorf("Error installing archive: %v", err)
	}

	// Define the destination directory (GOPATH/bin)
	destDir := filepath.Join(os.Getenv("GOPATH"), "bin")

	// Copy extracted files to the destination
	if err := CopyDir(tempDir, destDir); err != nil {
		return fmt.Errorf("Error installing archive: %v", err)
	}

	// Optionally delete the archive file
	// Uncomment the next line if you want to delete the archive after extraction
	// os.Remove(archiveFilename)

	return nil
}

func installToGOPATH(binaryPath string) error {
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		return fmt.Errorf("GOPATH is not set")
	}

	destination := filepath.Join(gopath, "bin", filepath.Base(binaryPath))

	// Open the source binary for reading
	sourceFile, err := os.Open(binaryPath)
	if err != nil {
		return fmt.Errorf("failed to open source file: %v", err)
	}
	defer sourceFile.Close()

	// Create the destination file
	destFile, err := os.Create(destination)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %v", err)
	}
	defer destFile.Close()

	// Copy the binary to the destination
	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return fmt.Errorf("failed to copy the binary: %v", err)
	}

	fmt.Printf("Successfully installed %s to %s\n", filepath.Base(binaryPath), destination)
	return nil
}

func installToGOPATH2(binaryPath string) error {
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		return fmt.Errorf("GOPATH is not set")
	}

	binDir := filepath.Join(gopath, "bin")
	binaryName := filepath.Base(binaryPath)
	destination := filepath.Join(binDir, binaryName)

	// Creating the destination directory if it doesn't exist
	if err := os.MkdirAll(binDir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create bin directory: %v", err)
	}

	// Opening the source binary file
	sourceFile, err := os.Open(binaryPath)
	if err != nil {
		return fmt.Errorf("failed to open source file: %v", err)
	}
	defer sourceFile.Close()

	// Creating the destination file
	destFile, err := os.Create(destination)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %v", err)
	}
	defer destFile.Close()

	// Copying the file
	if _, err := io.Copy(destFile, sourceFile); err != nil {
		return fmt.Errorf("failed to copy file: %v", err)
	}

	fmt.Printf("Successfully installed %s to %s\n", binaryName, destination)
	return nil
}
