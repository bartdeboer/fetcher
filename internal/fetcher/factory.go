package fetcher

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/bartdeboer/fetcher/internal/providers/github"
	"github.com/bartdeboer/fetcher/internal/providers/provider"
)

func NewRepoFromUrl(repoUrl, token string) (provider.Provider, error) {
	parsedURL, err := url.Parse(repoUrl)
	if err != nil {
		return nil, err
	}
	switch {
	case strings.Contains(parsedURL.Host, "github.com"):
		return github.New(repoUrl, token), nil
	default:
		fmt.Println("Unsupported Git service provider")
	}
	return nil, nil
}
