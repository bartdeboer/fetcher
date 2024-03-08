package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/bartdeboer/fetcher/internal/fetcher"
)

const (
	githubAPI    = "https://api.github.com/repos/"
	reposFile    = "repos.json"
	installsFile = "installs.json"
)

func main() {

	var command string
	// Check if there is at least one argument (the command)
	if len(os.Args) > 1 && !strings.HasPrefix(os.Args[1], "-") {
		command = os.Args[1]
		// Remove the command from os.Args
		os.Args = append(os.Args[:1], os.Args[2:]...)
	}

	repo := flag.String("repo", "", "GitHub repo in the format 'owner/repo'")
	// command := flag.String("command", "", "Command to execute: tap, list-taps, install, list-installs")
	flag.Parse()

	repos := fetcher.LoadRepos()

	switch command {
	case "tap":
		if *repo == "" {
			fmt.Println("Repository is required")
			os.Exit(1)
		}
		repos = append(repos, *repo)
		fetcher.SaveRepos(repos)
	case "list-taps":
		fetcher.ListRepos(repos)
	case "install":
		if *repo == "" {
			fmt.Println("Repository is required")
			os.Exit(1)
		}
		fetcher.InstallRelease(*repo)
	case "list-installs":
		fetcher.ListInstalls()
	default:
		fmt.Println("Unknown command")
		os.Exit(1)
	}
}
