package extractor

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// extractArchive determines the type of the archive and extracts it.
func ExtractArchive(filePath, destDir string) error {
	switch {
	case strings.HasSuffix(filePath, ".tar.gz"):
		return ExtractTarGz(filePath, destDir)
	case strings.HasSuffix(filePath, ".zip"):
		return ExtractZip(filePath, destDir)
	default:
		return fmt.Errorf("unsupported file extension")
	}
}

// extractTarGz handles the extraction of .tar.gz files.
func ExtractTarGz(filePath, destDir string) error {

	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("Error opening file: %v", err)
	}
	defer file.Close()

	gzipReader, err := gzip.NewReader(file)
	if err != nil {
		return fmt.Errorf("Error creating gzip.NewReader: %v", err)
	}
	defer gzipReader.Close()

	tarReader := tar.NewReader(gzipReader)
	for {
		header, err := tarReader.Next()
		switch {
		case err == io.EOF:
			return nil
		case err != nil:
			return fmt.Errorf("Error reading tarReader.Next: %v", err)
		case header == nil:
			continue
		}

		target := filepath.Join(destDir, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:

			fmt.Printf("Creating directory %s\n", target)

			if err := os.MkdirAll(target, os.FileMode(header.Mode)); err != nil {
				return fmt.Errorf("Error creating directory: %v", err)
			}

		case tar.TypeReg:

			fmt.Printf("Creating file %s\n", target)

			outFile, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return fmt.Errorf("Error creating file: %v", err)
			}

			if _, err := io.Copy(outFile, tarReader); err != nil {
				outFile.Close()
				return fmt.Errorf("Error extracting into file: %v", err)
			}
			outFile.Close()
		}
	}
}

// extractZip handles the extraction of .zip files.
func ExtractZip(filePath, destDir string) error {
	zipReader, err := zip.OpenReader(filePath)
	if err != nil {
		return err
	}
	defer zipReader.Close()

	for _, file := range zipReader.File {
		fPath := filepath.Join(destDir, file.Name)
		if file.FileInfo().IsDir() {
			os.MkdirAll(fPath, os.ModePerm)
			continue
		}

		if err := os.MkdirAll(filepath.Dir(fPath), os.ModePerm); err != nil {
			return err
		}

		outFile, err := os.OpenFile(fPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return err
		}

		rc, err := file.Open()
		if err != nil {
			outFile.Close()
			return err
		}

		_, err = io.Copy(outFile, rc)

		outFile.Close()
		rc.Close()

		if err != nil {
			return err
		}
	}
	return nil
}
