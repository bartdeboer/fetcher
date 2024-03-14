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

type implementation struct {
	creator func(url, token string) Provider
	hosts   []string
}

var impls map[string]implementation

func init() {
	impls = make(map[string]implementation)
}

func RegisterProvider(name string, creator func(url, token string) Provider, hosts []string) {
	impls[name] = implementation{
		creator,
		hosts,
	}
}

func NewProviderFromUrl(repoUrl, token string) (Provider, error) {
	parsedURL, err := url.Parse(repoUrl)
	if err != nil {
		return nil, err
	}
	for _, impl := range impls {
		for _, host := range impl.hosts {
			if strings.Contains(parsedURL.Host, host) {
				return impl.creator(repoUrl, token), nil
			}
		}
	}
	return nil, fmt.Errorf("unsupported Git service provider")
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

func (f *Fetcher) getRepo(name string) *Repo {
	for _, repo := range f.Repos {
		if repo.Url == name || strings.HasSuffix(repo.Url, "/"+name) {
			return repo
		}
	}
	return nil
}

func (f *Fetcher) GetRepo(name string) (*Repo, error) {
	foundRepo := f.getRepo(name)
	if foundRepo == nil {
		return nil, fmt.Errorf("repository not found: %s", name)
	}
	provider, err := NewProviderFromUrl(foundRepo.Url, foundRepo.Token)
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
	repo, err := f.GetRepo(repoName)
	if err != nil {
		return err
	}
	release, err := repo.LatestRelease()
	if err != nil {
		return err
	}
	for _, filename := range release.Files() {
		if err := release.FetchFile(filename); err != nil {
			return fmt.Errorf("error fetching file: %v", err)
		}
	}
	return nil
}

func (r *Repo) LatestRelease() (Release, error) {
	release, err := r.provider.LatestRelease(r.Url, r.Token)
	if err != nil {
		return nil, fmt.Errorf("error retrieving latest release: %v", err)
	}
	return release, nil
}
