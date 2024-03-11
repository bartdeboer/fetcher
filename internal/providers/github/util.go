package github

import (
	"fmt"
	"strings"
)

func extractRepoName(url string) (string, error) {
	// Remove potential trailing ".git"
	url = strings.TrimSuffix(url, ".git")

	// Check if URL is SSH format
	if strings.HasPrefix(url, "git@github.com:") {
		parts := strings.Split(url, ":")
		if len(parts) != 2 {
			return "", fmt.Errorf("invalid SSH GitHub URL")
		}
		return parts[1], nil
	} else if strings.HasPrefix(url, "https://github.com/") {
		// Check if URL is HTTPS format
		parts := strings.Split(url, "https://github.com/")
		if len(parts) != 2 {
			return "", fmt.Errorf("invalid HTTPS GitHub URL")
		}
		return parts[1], nil
	} else {
		return "", fmt.Errorf("unsupported URL format")
	}
}
