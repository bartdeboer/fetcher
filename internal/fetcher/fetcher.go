package fetcher

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"runtime"
	"strings"
)

const configFile = "fetcher.json"

type Fetcher struct {
	Repos []*Repo `json:"repositories"`
}

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

func (f *Fetcher) saveState() error {
	data, err := json.MarshalIndent(f, "", "  ")
	if err != nil {
		return fmt.Errorf("error writing json: %v", err)
	}
	return os.WriteFile(configFile, data, 0644)
}

func (f *Fetcher) SaveRepo(repoUrl string) error {
	if repoUrl == "" {
		return fmt.Errorf("url is required")
	}
	_, err := url.Parse(repoUrl)
	if err != nil {
		return fmt.Errorf("error parsing url %s: %v", repoUrl, err)
	}
	repo := f.findRepo(repoUrl)
	if repo != nil {
		return fmt.Errorf("repo already exists: %s", repoUrl)
	}
	f.Repos = append(f.Repos, &Repo{
		Url: repoUrl,
	})
	fmt.Printf("Update %s to add your token\n", configFile)
	return f.saveState()
}

func (f *Fetcher) findRepo(name string) *Repo {
	for _, repo := range f.Repos {
		if repo.Url == name || strings.HasSuffix(repo.Url, "/"+name) {
			return repo
		}
	}
	return nil
}

func (f *Fetcher) FindRepo(name string) (*Repo, error) {
	foundRepo := f.findRepo(name)
	if foundRepo == nil {
		return nil, fmt.Errorf("repository not found: %s", name)
	}
	provider, err := NewRepoFromUrl(foundRepo.Url)
	if err != nil {
		return nil, fmt.Errorf("error creating provider for %s: %w", foundRepo.Url, err)
	}
	foundRepo.provider = provider
	return foundRepo, nil
}

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

func (f *Fetcher) FetchAssets(repoName string) error {
	repo, err := f.FindRepo(repoName)
	if err != nil {
		return err
	}
	release, err := repo.LatestRelease()
	if err != nil {
		return err
	}
	for _, asset := range release.Assets() {
		if err := asset.Fetch(repo.Token); err != nil {
			return fmt.Errorf("error fetching file: %v", err)
		}
	}
	return nil
}

func (f *Fetcher) InstallAssets(repoName string) error {
	repo, err := f.FindRepo(repoName)
	if err != nil {
		return err
	}
	release, err := repo.LatestRelease()
	if err != nil {
		return err
	}
	isInstalled := false
	for _, asset := range release.Assets() {
		if !(strings.Contains(asset.Name(), runtime.GOOS) && strings.Contains(asset.Name(), runtime.GOARCH)) {
			continue
		}
		if err := asset.Fetch(repo.Token); err != nil {
			return fmt.Errorf("error fetching file: %v", err)
		}
		if err := InstallFromArchive(asset.Name()); err != nil {
			return fmt.Errorf("error installing file: %v", err)
		}
		isInstalled = true
	}
	if isInstalled {
		repo.InstalledTagName = release.TagName()
		f.saveState()
	}
	return nil
}
