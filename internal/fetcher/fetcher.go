package fetcher

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"strings"
)

const (
	githubAPI    = "https://api.github.com/repos/"
	reposFile    = "repos.json"
	installsFile = "installs.json"
)

type Release struct {
	TagName string `json:"tag_name"`
	Assets  []struct {
		Name               string `json:"name"`
		BrowserDownloadURL string `json:"browser_download_url"`
		Url                string `json:"url"`
	} `json:"assets"`
}

func LoadRepos() []string {
	data, err := os.ReadFile(reposFile)
	if err != nil {
		return []string{}
	}
	var repos []string
	json.Unmarshal(data, &repos)
	return repos
}

func SaveRepos(repos []string) {
	data, err := json.Marshal(repos)
	if err != nil {
		fmt.Println("Error saving repositories:", err)
		os.Exit(1)
	}
	os.WriteFile(reposFile, data, 0644)
}

func ListRepos(repos []string) {
	if len(repos) == 0 {
		fmt.Println("No tapped repositories")
		return
	}
	fmt.Println("Tapped repositories:")
	for _, repo := range repos {
		fmt.Println(repo)
	}
}

func GetLatestRelease(repo string) (*Release, error) {
	url := githubAPI + repo + "/releases/latest"

	fmt.Printf("Url: %s\n", url)

	// Create a new request
	req, err := http.NewRequest("GET", url, nil)
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

func DownloadFile(url, filename string) error {
	// Create a new request

	fmt.Printf("Url: %s\n", url)

	req, err := http.NewRequest("GET", url, nil)
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
	out, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer out.Close()

	// Copy the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}

func InstallRelease(repo string) {
	release, err := GetLatestRelease(repo)
	if err != nil {
		fmt.Println("Error fetching release:", err)
		os.Exit(1)
	}

	fmt.Println("Latest version:", release.TagName)
	installed := []string{}
	for _, asset := range release.Assets {
		if !(strings.Contains(asset.Name, runtime.GOOS) && strings.Contains(asset.Name, runtime.GOARCH)) {
			continue
		}
		installed = append(installed, asset.Name)
		fmt.Println("Downloading:", asset.Name)
		if err := DownloadFile(asset.Url, asset.Name); err != nil {
			fmt.Println("Error downloading file:", err)
			os.Exit(1)
		}
		if err := InstallFromArchive(asset.Name); err != nil {
			fmt.Println("Error installing file:", err)
			os.Exit(1)
		}
	}
	fmt.Println("Successfully installed all assets for", repo)
}

func ListInstalls() {
	installs := LoadInstalls()
	if len(installs) == 0 {
		fmt.Println("No installed software")
		return
	}
	fmt.Println("Installed software:")
	for _, install := range installs {
		fmt.Println(install)
	}
}

func LoadInstalls() []string {
	data, err := os.ReadFile(installsFile)
	if err != nil {
		return []string{}
	}
	var installs []string
	json.Unmarshal(data, &installs)
	return installs
}
