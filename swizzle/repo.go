package swizzle

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/afloesch/semver"
	"github.com/go-resty/resty/v2"
	"github.com/google/go-github/v45/github"
)

/*
Repo is a GitHub repository name which hosts swizzle mod
releases.

Example:
	// Create swizzle repo by name
	r := repo("afloesch/megamod")

	// Fetch release for v1.0.0
	rel, err := r.Release(context.TODO, "v1.0.0")
	if err != nil {
		fmt.Printf("missing release version: %s", err)
	}

	// Get release manifest
	manifest, _ := r.Manifest(context.TODO, rel)
*/
type Repo string

// Organization returns the repo owner.
func (r Repo) Organization() string {
	return strings.Split(r.String(), "/")[0]
}

// Name returns the repo name.
func (r Repo) Name() string {
	s := r.String()
	parts := strings.Split(s, "/")
	if len(parts) > 1 {
		return parts[1]
	}
	return s
}

// String return the full repo name as a string.
func (r Repo) String() string {
	return string(r)
}

// LatestManifest attempts to find the latest swizzle release from a
// github repository.
func (r Repo) LatestManifest(ctx context.Context) (*Manifest, error) {
	rel, err := r.Releases(ctx)
	if err != nil {
		return nil, err
	}

	var release *github.RepositoryRelease
	var asset *github.ReleaseAsset
	for _, r := range rel {
		for _, a := range r.Assets {
			if a.GetName() == manifestName {
				asset = a
				release = r
			}
		}
	}

	if asset == nil {
		return nil, fmt.Errorf("manifest not found")
	}

	return r.manifest(ctx, asset, release)
}

// Manifest fetches a swizzle Manifest from a release.
func (r Repo) Manifest(ctx context.Context, release *github.RepositoryRelease) (*Manifest, error) {
	if release == nil {
		return nil, fmt.Errorf("nil release")
	}

	var asset *github.ReleaseAsset
	for _, a := range release.Assets {
		if a.GetName() == manifestName {
			asset = a
		}
	}

	if asset == nil {
		return nil, fmt.Errorf("manifest not found")
	}

	return r.manifest(ctx, asset, release)
}

func (r Repo) manifest(ctx context.Context, asset *github.ReleaseAsset, release *github.RepositoryRelease) (*Manifest, error) {
	resp, err := r.FetchReleaseAsset(ctx, asset)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("invalid manifest file data: %s", err)
	}

	mani, err := ParseManifest(data)
	if err != nil {
		return nil, err
	}

	for _, f := range mani.Files {
		f.setReleaseAsset(release.Assets)
	}

	mani.release = release
	mani.releaseAsset = asset
	mani.Repo = r
	mani.Version = semver.String(release.GetTagName())
	return mani, nil
}

// Release fetches a repository release.
func (r Repo) Release(ctx context.Context, version string) (*github.RepositoryRelease, error) {
	ver := semver.String(version).Get()

	rel, err := r.Releases(ctx)
	if err != nil {
		return nil, fmt.Errorf("invalid repo: %s", err)
	}

	for _, d := range rel {
		if d.GetTagName() == ver.String() {
			return d, nil
		}
	}

	return nil, fmt.Errorf("no release for version '%s'", ver.String())
}

// LatestRelease fetches the latest release for a repository.
func (r Repo) LatestRelease(ctx context.Context) (*github.RepositoryRelease, error) {
	client := github.NewClient(http.DefaultClient)
	rel, res, err := client.Repositories.GetLatestRelease(ctx, r.Organization(), r.Name())
	if err != nil {
		return nil, fmt.Errorf("invalid repo: %s", err)
	}
	defer res.Body.Close()

	if res.StatusCode > 299 {
		return nil, fmt.Errorf("invalid repo")
	}

	return rel, nil
}

// Releases fetches a list of releases from a repository.
func (r Repo) Releases(ctx context.Context) ([]*github.RepositoryRelease, error) {
	client := github.NewClient(http.DefaultClient)
	rel, res, err := client.Repositories.ListReleases(ctx, r.Organization(), r.Name(), nil)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode > 299 {
		return nil, fmt.Errorf("invalid repo")
	}

	return rel, nil
}

// FetchReleaseAsset fetches a release asset and returns the http.Response from the request.
func (r Repo) FetchReleaseAsset(ctx context.Context, asset *github.ReleaseAsset) (*http.Response, error) {
	if asset == nil {
		return nil, fmt.Errorf("nil asset")
	}

	res, err := resty.New().R().
		SetDoNotParseResponse(true).
		SetContext(ctx).
		Get(asset.GetBrowserDownloadURL())

	if res.StatusCode() == 404 {
		return nil, fmt.Errorf("release file '%s' not found", asset.GetBrowserDownloadURL())
	}

	return res.RawResponse, err
}
