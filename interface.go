package fetcher

type Provider interface {
	LatestRelease(repo, token string) (Release, error)
}

type Release interface {
	TagName() string
	FetchFile(name string) error
	Files() []string
}
