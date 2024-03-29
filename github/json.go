package github

import (
	"encoding/json"
)

func (r *Release) UnmarshalJSON(data []byte) error {
	var tmp struct {
		TagName string  `json:"tag_name"`
		Assets  []Asset `json:"assets"`
	}
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	r.tagName = tmp.TagName
	for _, asset := range tmp.Assets {
		r.assets = append(r.assets, asset)
	}
	return nil
}
