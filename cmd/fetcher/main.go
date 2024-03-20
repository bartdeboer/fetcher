package main

import (
	"fmt"
	"os"

	"github.com/bartdeboer/fetcher"
	"github.com/bartdeboer/fetcher/github"
	"github.com/spf13/cobra"
)

func init() {
	fetcher.RegisterProvider("github", github.New())
}

func main() {

	f, err := fetcher.NewFetcherFromConfig()
	if err != nil {
		fmt.Printf("error launching fetcher: %v\n", err)
		os.Exit(1)
	}

	rootCmd := &cobra.Command{
		Use:   "fetcher [command]",
		Short: "Fetcher is a tool for fetching and installing releases",
	}

	cmdTap := &cobra.Command{
		Use:   "tap [repo]",
		Short: "Saves a repository",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			repo := args[0]
			if err := f.SaveRepo(repo); err != nil {
				fmt.Printf("Error saving repository %s: %v\n", repo, err)
				os.Exit(1)
			}
		},
	}

	cmdListTaps := &cobra.Command{
		Use:   "list",
		Short: "Lists all saved repositories",
		Run: func(cmd *cobra.Command, args []string) {
			f.ListRepos()
		},
	}

	cmdDownload := &cobra.Command{
		Use:   "download [repo]",
		Short: "Fetches the latest release assets for a repository",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			repo := args[0]
			if err := f.FetchAssets(repo); err != nil {
				fmt.Printf("Error retrieving latest release %s: %v\n", repo, err)
				os.Exit(1)
			}
		},
	}

	cmdInstall := &cobra.Command{
		Use:   "install [repo]",
		Short: "Installs the latest release assets for a repository",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			repo := args[0]
			if err := f.InstallAssets(repo); err != nil {
				fmt.Printf("Error installing release %s: %v\n", repo, err)
				os.Exit(1)
			}
		},
	}

	rootCmd.AddCommand(cmdTap, cmdListTaps, cmdDownload, cmdInstall)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
