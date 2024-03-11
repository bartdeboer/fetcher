package factory

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/bartdeboer/fetcher/internal/providers"
	"github.com/bartdeboer/fetcher/internal/providers/github"
)

func NewRepoFromUrl(repoUrl string) (providers.Provider, error) {
	parsedURL, err := url.Parse(repoUrl)
	if err != nil {
		return nil, err
	}
	switch {
	case strings.Contains(parsedURL.Host, "github.com"):
		return github.NewGithub(repoUrl), nil
	default:
		fmt.Println("Unsupported Git service provider")
	}
	return nil, nil
}
