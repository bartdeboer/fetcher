package provider

type Provider interface {
	LatestRelease(repo, token string) (Release, error)
}

type Release interface {
	TagName() string
	Assets() []Asset
	// FetchAssets() error
}

type Asset interface {
	Name() string
	Fetch(token string) error
}
