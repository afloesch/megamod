package internal

import (
	"context"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/go-resty/resty/v2"
)

type Repo string

func (r Repo) Organization() string {
	return strings.Split(r.String(), "/")[0]
}

func (r Repo) Name() string {
	parts := strings.Split(r.String(), "/")
	if len(parts) > 1 {
		return parts[1]
	}
	return r.Name()
}

func (r Repo) String() string {
	return string(r)
}

func (r Repo) FetchManifest(ctx context.Context, version SemVer) (*Manifest, error) {
	resp, err := r.FetchReleaseFile(ctx, version.Get(), manifestName)
	if err != nil {
		return nil, err
	}
	defer resp.RawBody().Close()

	data, err := ioutil.ReadAll(resp.RawBody())
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

func (r Repo) FetchReleaseFile(ctx context.Context, version *Version, file string) (*resty.Response, error) {
	return resty.
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
}
