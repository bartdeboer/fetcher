package fetcher

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/bartdeboer/fetcher/internal/extractor"
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

	if err := extractor.ExtractArchive(archiveFilename, extractDir); err != nil {
		return fmt.Errorf("error extracting archive: %v", err)
	}

	// Define the destination directory (GOPATH/bin)
	destDir := filepath.Join(os.Getenv("GOPATH"), "bin")

	// Copy extracted files to the destination
	if err := CopyDir(extractDir, destDir); err != nil {
		return fmt.Errorf("error installing files: %v", err)
	}

	// Optionally delete the archive file
	// Uncomment the next line if you want to delete the archive after extraction
	// os.Remove(archiveFilename)

	return nil
}
