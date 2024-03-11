package github

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"strings"

	"github.com/bartdeboer/fetcher/internal/providers"
)

const (
	githubAPI    = "https://api.github.com/repos/"
	apiBase      = "https://api.github.com/repos/"
	reposFile    = "repos.json"
	installsFile = "installs.json"
)

type Github struct {
	url string
}

func NewGithub(url string) *Github {
	return &Github{
		url,
	}
}

type Release struct {
	tagName string
	assets  []providers.Asset
}

type Asset struct {
	name               string
	BrowserDownloadURL string
	url                string
}

func (r *Release) TagName() string {
	return r.tagName
}

func (r *Release) Assets() []providers.Asset {
	return r.assets
}

func (a *Asset) Name() string {
	return a.name
}

func (a *Asset) Url() string {
	return a.url
}

func (r *Release) FetchAssets() error {
	for _, asset := range r.Assets() {
		if !(strings.Contains(asset.Name(), runtime.GOOS) && strings.Contains(asset.Name(), runtime.GOARCH)) {
			continue
		}
		err := asset.Fetch()
		if err != nil {
			return fmt.Errorf("error fetching asset: %v", err)
		}
	}
	return nil
}

func (a *Asset) Fetch() error {

	fmt.Printf("Url: %s\n", a.Url())

	req, err := http.NewRequest("GET", a.Url(), nil)
	if err != nil {
		return err
	}

	// Retrieve the GitHub token from the environment
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		return fmt.Errorf("GitHub token not found in environment")
	}

	// Set the Authorization header
	req.Header.Set("Authorization", "token "+token)
	req.Header.Set("Accept", "application/octet-stream")

	// Make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Print response headers
	for name, values := range resp.Header {
		for _, value := range values {
			fmt.Printf("%s: %s\n", name, value)
		}
	}

	// Check for non-OK status code
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download file: %d %s", resp.StatusCode, resp.Status)
	}

	// Create the file
	out, err := os.Create(a.Name())
	if err != nil {
		return err
	}
	defer out.Close()

	// Copy the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}

func (g *Github) LatestRelease(repoUrl string) (providers.Release, error) {

	repoName, err := extractRepoName(repoUrl)
	if err != nil {
		return nil, err
	}

	apiUrl := githubAPI + repoName + "/releases/latest"

	fmt.Printf("Url: %s\n", apiUrl)

	// Create a new request
	req, err := http.NewRequest("GET", apiUrl, nil)
	if err != nil {
		return nil, err
	}

	// Retrieve the GitHub token from the environment
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		return nil, fmt.Errorf("GitHub token not found in environment")
	}

	// Set the Authorization header
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	// Make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Print response headers
	for name, values := range resp.Header {
		for _, value := range values {
			fmt.Printf("%s: %s\n", name, value)
		}
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API error: %d %s", resp.StatusCode, resp.Status)
	}

	var release Release
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, err
	}

	return &release, nil
}
