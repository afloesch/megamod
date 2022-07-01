package swizzle

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/google/go-github/v45/github"
)

/*
Repo is a GitHub repository name which hosts swizzle mod
releases.

Example:
	// get a swizzle repo by name
	r := repo("afloesch/megamod")

	// fetch swiz.zle manifest file for v1.0.0
	m, err := r.FetchManifest(context.TODO, swizzle.SemVer("v1.0.0"))
	if err != nil {
		fmt.Printf("missing release version: %s", err)
	}
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

// FetchManifest fetches a release swizzle Manifest.
func (r Repo) FetchManifest(ctx context.Context, version SemVer) (*Manifest, error) {
	rel, err := r.FetchRelease(ctx, version.Get())
	if err != nil {
		return nil, err
	}

	if rel == nil {
		return nil, fmt.Errorf("fail")
	}

	var asset *github.ReleaseAsset
	for _, a := range rel.Assets {
		if a.GetName() == manifestName {
			asset = a
		}
	}

	resp, err := r.FetchReleaseFile(ctx, asset)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("invalid manifest: %s", err)
	}

	mani, err := ParseManifest(data)
	if err != nil {
		return nil, err
	}

	for _, f := range mani.Files {
		for _, a := range rel.Assets {
			if a.GetName() == f.Name {
				f.releaseAsset = a
			}
		}
	}

	mani.release = rel
	mani.Repo = r
	mani.Version = version
	return mani, nil
}

func (r Repo) FetchRelease(ctx context.Context, version *Version) (*github.RepositoryRelease, error) {
	rel, err := r.Releases(ctx)
	if err != nil {
		return nil, fmt.Errorf("invalid repo: %s", err)
	}

	for _, d := range rel {
		if d.GetTagName() == version.String() {
			return d, nil
		}
	}

	return nil, fmt.Errorf("no release for version '%s'", version.String())
}

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

// FetchReleaseFile makes a request for a release file for a particular version.
func (r Repo) FetchReleaseFile(ctx context.Context, asset *github.ReleaseAsset) (*http.Response, error) {

	if asset == nil {
		return nil, fmt.Errorf("release file '%s' not found", asset.GetBrowserDownloadURL())
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
