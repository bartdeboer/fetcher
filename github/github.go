package github

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/bartdeboer/fetcher"
)

const apiBaseUrl = "https://api.github.com/repos/"

var hosts = []string{"github.com"}

type Repo struct{}

func New() *Repo {
	return &Repo{}
}

type Release struct {
	tagName string
	assets  []Asset
}

type Asset struct {
	Name               string `json:"name"`
	Url                string `json:"url"`
	BrowserDownloadURL string `json:"browser_download_url"`
}

func (r *Release) TagName() string {
	return r.tagName
}

func (r *Release) Files() []string {
	var files []string
	for _, asset := range r.assets {
		files = append(files, asset.Name)
	}
	return files
}

func (r *Release) findAsset(name string) (*Asset, error) {
	for i, asset := range r.assets {
		if asset.Name == name {
			return &r.assets[i], nil
		}
	}
	return nil, fmt.Errorf("asset %s not found", name)
}

// Retrieves the release asset
func (r *Release) FetchFile(name, token string) error {

	asset, err := r.findAsset(name)
	if err != nil {
		return err
	}

	fmt.Printf("Url: %s\n", asset.Url)

	req, err := http.NewRequest("GET", asset.Url, nil)
	if err != nil {
		return err
	}

	if token == "" {
		token = os.Getenv("GITHUB_TOKEN")
	}

	if token != "" {
		req.Header.Set("Authorization", "token "+token)
	}

	req.Header.Set("Accept", "application/octet-stream")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	for name, values := range resp.Header {
		for _, value := range values {
			fmt.Printf("%s: %s\n", name, value)
		}
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download file: %d %s", resp.StatusCode, resp.Status)
	}

	out, err := os.Create(asset.Name)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func (g *Repo) CanHandleHost(host string) bool {
	for i, _ := range hosts {
		if strings.Contains(host, hosts[i]) {
			return true
		}
	}
	return false
}

// Finds and returns the latest release
func (g *Repo) LatestRelease(repoUrl, token string) (fetcher.Release, error) {

	repoName, err := extractRepoName(repoUrl)
	if err != nil {
		return nil, err
	}

	apiUrl := apiBaseUrl + repoName + "/releases/latest"

	fmt.Printf("Url: %s\n", apiUrl)

	req, err := http.NewRequest("GET", apiUrl, nil)
	if err != nil {
		return nil, err
	}

	if token == "" {
		token = os.Getenv("GITHUB_TOKEN")
	}

	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

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
