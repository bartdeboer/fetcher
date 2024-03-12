package fetcher

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/bartdeboer/archiver"
)

func createExtractDir(archiveFilename string, useTemp bool) (string, error) {
	base := filepath.Base(archiveFilename)
	base = strings.TrimSuffix(base, ".tar.gz")
	base = strings.TrimSuffix(base, ".zip")
	var dir string
	if useTemp {
		dir = os.TempDir()
	} else {
		dir = "."
	}
	extractDir := filepath.Join(dir, base)
	if err := os.Mkdir(extractDir, 0755); err != nil {
		return "", fmt.Errorf("error creating destination directory %s: %v", extractDir, err)
	}
	return extractDir, nil
}

func InstallFromArchive(archiveFilename string) error {

	extractDir, err := createExtractDir(archiveFilename, true)
	if err != nil {
		return fmt.Errorf("error creating extract dir: %v", err)
	}

	defer os.RemoveAll(extractDir)

	if err := archiver.Extract(archiveFilename, extractDir); err != nil {
		return fmt.Errorf("error extracting archive: %v", err)
	}

	// Define the destination directory (GOPATH/bin)
	destDir := filepath.Join(os.Getenv("GOPATH"), "bin")

	// Copy extracted files to the destination
	if err := CopyDir(extractDir, destDir); err != nil {
		return fmt.Errorf("error installing files: %v", err)
	}

	// Optionally delete the archive file
	// os.Remove(archiveFilename)

	return nil
}
