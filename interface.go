package fetcher

type Provider interface {
	LatestRelease(repo, token string) (Release, error)
	CanHandleHost(host string) bool
}

type Release interface {
	TagName() string
	FetchFile(name, token string) error
	Files() []string
}
