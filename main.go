package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/bartdeboer/fetcher/internal/fetcher"
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

	f, err := fetcher.NewFetcherFromConfig()
	if err != nil {
		fmt.Printf("error launcing fetcher: %v\n", err)
		os.Exit(1)
	}

	switch command {
	case "tap":
		if err := f.SaveRepo(*frepo); err != nil {
			fmt.Printf("Error saving repository: %s %v\n", *frepo, err)
			os.Exit(1)
		}
	case "list-taps":
		f.ListRepos()
	case "download":
		if err := f.FetchAssets(*frepo); err != nil {
			fmt.Printf("Error retrieving latest release: %s %v\n", *frepo, err)
			os.Exit(1)
		}
	case "install":
		if err := f.InstallAssets(*frepo); err != nil {
			fmt.Printf("Error retrieving latest release: %s %v\n", *frepo, err)
			os.Exit(1)
		}
	case "list-installs":
		// fetcher.ListInstalls()
	default:
		fmt.Println("Unknown command")
		os.Exit(1)
	}
}
