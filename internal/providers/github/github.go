package github

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/bartdeboer/fetcher/internal/providers/provider"
)

const apiBaseUrl = "https://api.github.com/repos/"

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
	assets  []provider.Asset
}

func (r *Release) TagName() string          { return r.tagName }
func (r *Release) Assets() []provider.Asset { return r.assets }

type Asset struct {
	name               string
	BrowserDownloadURL string
	url                string
}

func (a *Asset) Name() string { return a.name }
func (a *Asset) Url() string  { return a.url }

func (a *Asset) Fetch(token string) error {

	fmt.Printf("Url: %s\n", a.Url())

	req, err := http.NewRequest("GET", a.Url(), nil)
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

	out, err := os.Create(a.Name())
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func (g *Github) LatestRelease(repoUrl, token string) (provider.Release, error) {

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
