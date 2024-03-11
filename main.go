package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/bartdeboer/fetcher/internal/fetcher"
)

const (
	reposFile = "repos.json"
)

func main() {

	var command string

	// Check if there is at least one argument (the command)
	if len(os.Args) > 1 && !strings.HasPrefix(os.Args[1], "-") {
		command = os.Args[1]
		// Remove the command from os.Args
		os.Args = append(os.Args[:1], os.Args[2:]...)
	}

	frepo := flag.String("repo", "", "GitHub repo in the format 'owner/repo'")
	flag.Parse()

	f, err := fetcher.NewFetcherFromConfig(reposFile)
	if err != nil {
		fmt.Printf("error launcing fetcher: %v\n", err)
		os.Exit(1)
	}

	switch command {
	case "tap":
		if *frepo == "" {
			fmt.Println("Repository is required")
			os.Exit(1)
		}
		f.SaveRepo(*frepo)
	case "list-taps":
		f.ListRepos()
	case "download":
		if *frepo == "" {
			fmt.Println("Repository is required")
			os.Exit(1)
		}
		repo := f.FindRepo(*frepo)
		if repo == nil {
			fmt.Printf("Could not find repo: %s\n", *frepo)
			os.Exit(1)
		}
		latestRelease, err := repo.LatestRelease()
		if err == nil {
			fmt.Printf("Error retrieving latest release: %s\n", *frepo)
			os.Exit(1)
		}
		latestRelease.FetchAssets()
	case "install":
		if *frepo == "" {
			fmt.Println("Repository is required")
			os.Exit(1)
		}
		repo := f.FindRepo(*frepo)
		if repo == nil {
			fmt.Printf("Could not find repo: %s\n", *frepo)
			os.Exit(1)
		}
		latestRelease, err := repo.LatestRelease()
		if err == nil {
			fmt.Printf("Error retrieving latest release: %s\n", *frepo)
			os.Exit(1)
		}
		latestRelease.FetchAssets()
		repo.InstallAssets(latestRelease.Assets())
	case "list-installs":
		// fetcher.ListInstalls()
	default:
		fmt.Println("Unknown command")
		os.Exit(1)
	}
}
