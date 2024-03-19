package fetcher

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/bartdeboer/archiver"
)

// Creates a extract directory for the archive
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

// Installs Go binaries from the archive.
func installFromArchive(archiveFilename string) error {

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
	if err := copyDir(extractDir, destDir); err != nil {
		return fmt.Errorf("error installing files: %v", err)
	}

	// Optionally delete the archive file
	// os.Remove(archiveFilename)

	return nil
}

// Installs the assets from the latest release for the current platform.
func (f *Fetcher) InstallAssets(repoName string) error {
	repo, err := f.GetRepo(repoName)
	if err != nil {
		return err
	}
	release, err := repo.LatestRelease()
	if err != nil {
		return err
	}
	var installedFile string
	for _, filename := range release.Files() {
		if !(strings.Contains(filename, runtime.GOOS+"_"+runtime.GOARCH)) {
			continue
		}
		if err := release.FetchFile(filename); err != nil {
			return fmt.Errorf("error fetching file: %v", err)
		}
		if err := installFromArchive(filename); err != nil {
			return fmt.Errorf("error installing file: %v", err)
		}
		installedFile = filename
		break
	}
	if installedFile != "" {
		repo.InstalledTagName = release.TagName()
		repo.InstalledFilename = installedFile
		f.saveState()
	}
	return nil
}
