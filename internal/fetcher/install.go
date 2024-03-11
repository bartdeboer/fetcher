package fetcher

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/bartdeboer/fetcher/internal/extractor"
)

func extractFromArchive(archiveFilename string, useTemp bool) (string, error) {

	// Extract the basename of the archive
	base := filepath.Base(archiveFilename)
	base = strings.TrimSuffix(base, ".tar.gz")
	base = strings.TrimSuffix(base, ".zip")

	var dir string
	if useTemp {
		dir = os.TempDir()
	} else {
		dir = "."
	}

	// Create a temporary directory
	destDir := filepath.Join(dir, base)
	if err := os.Mkdir(destDir, 0755); err != nil {
		return "", fmt.Errorf("error extracting archive: %v", err)
	}
	defer os.RemoveAll(destDir)

	// Extract the archive into the temporary directory
	if err := extractor.ExtractArchive(archiveFilename, destDir); err != nil {
		return "", fmt.Errorf("error extracting archive: %v", err)
	}

	return destDir, nil
}

func InstallFromArchive(archiveFilename string) error {

	extractDir, err := extractFromArchive(archiveFilename, true)
	if err != nil {
		return fmt.Errorf("error installing archive: %v", err)
	}

	// Define the destination directory (GOPATH/bin)
	destDir := filepath.Join(os.Getenv("GOPATH"), "bin")

	// Copy extracted files to the destination
	if err := CopyDir(extractDir, destDir); err != nil {
		return fmt.Errorf("error installing archive: %v", err)
	}

	// Optionally delete the archive file
	// Uncomment the next line if you want to delete the archive after extraction
	// os.Remove(archiveFilename)

	return nil
}
