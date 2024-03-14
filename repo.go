package fetcher

import (
	"fmt"
)

type Repo struct {
	Url               string `json:"url"`
	InstalledFilename string `json:"installed_filename"`
	InstalledTagName  string `json:"installed_tag_name"`
	Token             string `json:"token"`
	provider          Provider
}

func (r *Repo) LatestRelease() (Release, error) {
	release, err := r.provider.LatestRelease(r.Url, r.Token)
	if err != nil {
		return nil, fmt.Errorf("error retrieving latest release: %v", err)
	}
	return release, nil
}
