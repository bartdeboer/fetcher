package providers

type Provider interface {
	LatestRelease(repo string) (Release, error)
}

type Release interface {
	TagName() string
	Assets() []Asset
	FetchAssets() error
}

type Asset interface {
	Name() string
	Fetch() error
}
