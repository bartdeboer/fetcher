package fetcher

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/bartdeboer/fetcher/internal/providers"
)

const (
	githubAPI    = "https://api.github.com/repos/"
	reposFile    = "repos.json"
	installsFile = "installs.json"
)

type Fetcher struct {
	reposFilename string
	repos         []*Repo
}

type Repo struct {
	Url                string `json:"url"`
	InstalledAssetName string `json:"installed_asset_name"`
	InstalledTagName   string `json:"installed_tag_name"`
	provider           providers.Provider
}

func NewFetcherFromConfig(filename string) (*Fetcher, error) {
	f := &Fetcher{
		reposFilename: filename,
	}
	if filename == "" {
		return nil, fmt.Errorf("no config file")
	}
	err := f.loadRepos(filename)
	if err != nil {
		return nil, fmt.Errorf("error loading repos: %v", err)
	}
	return f, nil
}

func (f *Fetcher) loadRepos(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("error reading file: %s %v", filename, err)
	}
	json.Unmarshal(data, &f.repos)
	return nil
}

func (f *Fetcher) saveRepos() error {
	data, err := json.Marshal(f.repos)
	if err != nil {
		return fmt.Errorf("error writing json: %v", err)
	}
	return os.WriteFile(reposFile, data, 0644)
}

func (f *Fetcher) SaveRepo(repoUrl string) error {
	_, err := url.Parse(repoUrl)
	if err != nil {
		return fmt.Errorf("error parsing url: %s %v", repoUrl, err)
	}
	f.repos = append(f.repos, &Repo{
		Url: repoUrl,
	})
	return f.saveRepos()
}

func (f *Fetcher) FindRepo(name string) *Repo {
	for _, repo := range f.repos {
		if repo.Url == name || strings.HasSuffix(repo.Url, "/"+name) {
			return repo
		}
	}
	return nil
}

func (f *Fetcher) ListRepos() {
	if len(f.repos) == 0 {
		fmt.Println("No tapped repositories")
		return
	}
	fmt.Println("Tapped repositories:")
	for _, repo := range f.repos {
		fmt.Println(repo.Url)
	}
}

func (r *Repo) LatestRelease() (providers.Release, error) {
	// TODO get release tag
	release, err := r.provider.LatestRelease(r.Url)
	if err != nil {
		return nil, fmt.Errorf("error retrieving latest release: %v", err)
	}
	r.InstalledTagName = release.TagName()
	// TODO: Save updated repo
	return r.provider.LatestRelease(r.Url)
}

func (r *Repo) InstallAssets(release providers.Release) error {
	// TODO update repo with installed release tag
	for _, asset := range release.Assets() {
		if err := InstallFromArchive(asset.Name()); err != nil {
			fmt.Println("Error installing file:", err)
			os.Exit(1)
		}
	}
	return nil
}
