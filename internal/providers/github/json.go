package github

import (
	"encoding/json"

	"github.com/bartdeboer/fetcher/internal/providers/provider"
)

func (r *Release) UnmarshalJSON(data []byte) error {
	var tmp struct {
		TagName string   `json:"tag_name"`
		Assets  []*Asset `json:"assets"`
	}
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	r.tagName = tmp.TagName
	interfaceAssets := make([]provider.Asset, len(tmp.Assets))
	for i, asset := range tmp.Assets {
		interfaceAssets[i] = asset
	}
	r.assets = interfaceAssets
	return nil
}

func (a *Asset) UnmarshalJSON(data []byte) error {
	var tmp struct {
		Name               string `json:"name"`
		Url                string `json:"url"`
		BrowserDownloadURL string `json:"browser_download_url"`
	}
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	a.name = tmp.Name
	a.url = tmp.Url
	a.BrowserDownloadURL = tmp.BrowserDownloadURL
	return nil
}
