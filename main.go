package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
)

const githubAPI = "https://api.github.com/repos/"

type Release struct {
	TagName string `json:"tag_name"`
	Assets  []struct {
		Name               string `json:"name"`
		BrowserDownloadURL string `json:"browser_download_url"`
	} `json:"assets"`
}

func main() {
	repo := flag.String("repo", "", "GitHub repo in the format 'owner/repo'")
	flag.Parse()

	if *repo == "" {
		fmt.Println("Repository is required")
		os.Exit(1)
	}

	release, err := getLatestRelease(*repo)
	if err != nil {
		fmt.Println("Error fetching release:", err)
		os.Exit(1)
	}

	fmt.Println("Latest version:", release.TagName)
	for _, asset := range release.Assets {
		fmt.Println("Downloading:", asset.Name)
		if err := downloadFile(asset.BrowserDownloadURL, asset.Name); err != nil {
			fmt.Println("Error downloading file:", err)
			os.Exit(1)
		}
	}
}

func getLatestRelease(repo string) (*Release, error) {
	url := githubAPI + repo + "/releases/latest"
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API error: %s", resp.Status)
	}

	var release Release
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, err
	}

	return &release, nil
}

func downloadFile(url, filename string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}
