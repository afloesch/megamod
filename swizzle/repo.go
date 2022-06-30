package swizzle

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/go-resty/resty/v2"
)

/*
Repo is a GitHub repository name which hosts swizzle mod
releases. For example, afloesch/megamod is the name for
this project's repository.

Example:
	// create a repo from name
	r := repo("afloesch/megamod")

	// fetch manifest for v1.0.0
	m, err := r.FetchManifest(context.Background(), "v1.0.0")
	if err != nil {
		fmt.Errorf("missing release version", err)
	}
*/
type Repo string

// Organization returns the repo owner.
func (r Repo) Organization() string {
	return strings.Split(r.String(), "/")[0]
}

// Name returns the repo name.
func (r Repo) Name() string {
	parts := strings.Split(r.String(), "/")
	if len(parts) > 1 {
		return parts[1]
	}
	return r.Name()
}

// String return the full repo name as a string.
func (r Repo) String() string {
	return string(r)
}

// FetchManifest fetches a release swizzle Manifest.
func (r Repo) FetchManifest(ctx context.Context, version SemVer) (*Manifest, error) {
	resp, err := r.FetchReleaseFile(ctx, version.Get(), manifestName)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	mani, err := ParseManifest(data)
	if err != nil {
		return nil, err
	}

	mani.Repo = r
	mani.Version = version
	return mani, nil
}

// FetchReleaseFile makes a request for a release file for a particular version.
func (r Repo) FetchReleaseFile(ctx context.Context, version *Version, file string) (*http.Response, error) {
	res, err := resty.
		New().
		R().
		SetDoNotParseResponse(true).
		SetContext(ctx).
		Get(fmt.Sprintf(
			"https://github.com/%s/releases/download/%s/%s",
			r.String(),
			version.String(),
			file,
		))

	return res.RawResponse, err
}
