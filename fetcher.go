package fetcher

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"strings"
)

const configFile = "fetcher.json"

type Fetcher struct {
	Repos []*Repo `json:"repositories"`
}

type Repo struct {
	Url               string `json:"url"`
	InstalledFilename string `json:"installed_filename"`
	InstalledTagName  string `json:"installed_tag_name"`
	Token             string `json:"token"`
	provider          Provider
}

var providers map[string]Provider

func init() {
	providers = make(map[string]Provider)
}

func RegisterProvider(name string, provider Provider) {
	providers[name] = provider
}

func NewProviderFromUrl(repoUrl string) (Provider, error) {
	parsedURL, err := url.Parse(repoUrl)
	if err != nil {
		return nil, err
	}
	for _, provider := range providers {
		if provider.CanHandleHost(parsedURL.Host) {
			return provider, nil
		}
	}
	return nil, fmt.Errorf("unsupported Git service provider")
}

// Creates a new Fetcher instance (from config file)
func NewFetcherFromConfig() (*Fetcher, error) {
	var f *Fetcher
	data, err := os.ReadFile(configFile)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, fmt.Errorf("error reading file %s: %v", configFile, err)
		}
		return &Fetcher{}, nil
	}
	if err := json.Unmarshal(data, &f); err != nil {
		return nil, fmt.Errorf("error unmarshaling json: %v", err)
	}
	return f, nil
}

// Saves the fetcher's state to a configuration file.
func (f *Fetcher) saveState() error {
	data, err := json.MarshalIndent(f, "", "  ")
	if err != nil {
		return fmt.Errorf("error writing json: %v", err)
	}
	return os.WriteFile(configFile, data, 0644)
}

// Adds (taps) a new repo to the fetcher and saves the configuration.
func (f *Fetcher) SaveRepo(repoUrl string) error {
	if repoUrl == "" {
		return fmt.Errorf("url is required")
	}
	_, err := url.Parse(repoUrl)
	if err != nil {
		return fmt.Errorf("error parsing url %s: %v", repoUrl, err)
	}
	repo := f.getRepo(repoUrl)
	if repo != nil {
		return fmt.Errorf("repo already exists: %s", repoUrl)
	}
	f.Repos = append(f.Repos, &Repo{
		Url: repoUrl,
	})
	fmt.Printf("Update %s to add your token\n", configFile)
	return f.saveState()
}

// Finds a repo config by URL or name within the fetcher.
func (f *Fetcher) getRepo(name string) *Repo {
	for _, repo := range f.Repos {
		if repo.Url == name || strings.HasSuffix(repo.Url, "/"+name) {
			return repo
		}
	}
	return nil
}

// Retrieves and prepares a repo for use, including its provider.
func (f *Fetcher) GetRepo(name string) (*Repo, error) {
	foundRepo := f.getRepo(name)
	if foundRepo == nil {
		return nil, fmt.Errorf("repository not found: %s", name)
	}
	provider, err := NewProviderFromUrl(foundRepo.Url)
	if err != nil {
		return nil, fmt.Errorf("error creating provider for %s: %w", foundRepo.Url, err)
	}
	foundRepo.provider = provider
	return foundRepo, nil
}

// Lists all repositories managed by the fetcher.
func (f *Fetcher) ListRepos() {
	if len(f.Repos) == 0 {
		fmt.Println("No tapped repositories")
		return
	}
	fmt.Println("Tapped repositories:")
	for _, repo := range f.Repos {
		fmt.Println(repo.Url)
	}
}

// Downloads assets from a repo's latest release.
func (f *Fetcher) FetchAssets(repoName string) error {
	repo, err := f.GetRepo(repoName)
	if err != nil {
		return err
	}
	release, err := repo.LatestRelease()
	if err != nil {
		return err
	}
	for _, filename := range release.Files() {
		if err := release.FetchFile(filename, repo.Token); err != nil {
			return fmt.Errorf("error fetching file: %v", err)
		}
	}
	return nil
}

// Gets the latest release of a repository.
func (r *Repo) LatestRelease() (Release, error) {
	release, err := r.provider.LatestRelease(r.Url, r.Token)
	if err != nil {
		return nil, fmt.Errorf("error retrieving latest release: %v", err)
	}
	return release, nil
}
