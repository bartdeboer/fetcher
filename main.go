package main

import (
	"encoding/json"
	"flag"
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
	} `json:"assets"`
}

func main() {

	var command string
	// Check if there is at least one argument (the command)
	if len(os.Args) > 1 && !strings.HasPrefix(os.Args[1], "-") {
		command = os.Args[1]
		// Remove the command from os.Args
		os.Args = append(os.Args[:1], os.Args[2:]...)
	}

	repo := flag.String("repo", "", "GitHub repo in the format 'owner/repo'")
	// command := flag.String("command", "", "Command to execute: tap, list-taps, install, list-installs")
	flag.Parse()

	repos := loadRepos()

	switch command {
	case "tap":
		if *repo == "" {
			fmt.Println("Repository is required")
			os.Exit(1)
		}
		repos = append(repos, *repo)
		saveRepos(repos)
	case "list-taps":
		listRepos(repos)
	case "install":
		if *repo == "" {
			fmt.Println("Repository is required")
			os.Exit(1)
		}
		installRelease(*repo)
	case "list-installs":
		listInstalls()
	default:
		fmt.Println("Unknown command")
		os.Exit(1)
	}
}

func loadRepos() []string {
	data, err := os.ReadFile(reposFile)
	if err != nil {
		return []string{}
	}
	var repos []string
	json.Unmarshal(data, &repos)
	return repos
}

func saveRepos(repos []string) {
	data, err := json.Marshal(repos)
	if err != nil {
		fmt.Println("Error saving repositories:", err)
		os.Exit(1)
	}
	os.WriteFile(reposFile, data, 0644)
}

func listRepos(repos []string) {
	if len(repos) == 0 {
		fmt.Println("No tapped repositories")
		return
	}
	fmt.Println("Tapped repositories:")
	for _, repo := range repos {
		fmt.Println(repo)
	}
}

func getLatestRelease(repo string) (*Release, error) {
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

func installRelease(repo string) {
	release, err := getLatestRelease(repo)
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
		if err := downloadFile(asset.BrowserDownloadURL, asset.Name); err != nil {
			fmt.Println("Error downloading file:", err)
			os.Exit(1)
		}
		if err := installFromArchive(asset.Name); err != nil {
			fmt.Println("Error installing file:", err)
			os.Exit(1)
		}
	}
	fmt.Println("Successfully installed all assets for", repo)
}

func listInstalls() {
	installs := loadInstalls()
	if len(installs) == 0 {
		fmt.Println("No installed software")
		return
	}
	fmt.Println("Installed software:")
	for _, install := range installs {
		fmt.Println(install)
	}
}

func loadInstalls() []string {
	data, err := os.ReadFile(installsFile)
	if err != nil {
		return []string{}
	}
	var installs []string
	json.Unmarshal(data, &installs)
	return installs
}
