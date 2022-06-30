package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/go-resty/resty/v2"
	"gopkg.in/yaml.v3"
)

const manifestName string = "swiz.zle"

type Game struct {
	Executable string `json:"executable,omitempty" yaml:"executable,omitempty"`
	Version    string `json:"version,omitempty" yaml:"version,omitempty"`
}

type Manifest struct {
	AgeRating   AgeRating         `json:"ages,omitempty" yaml:"ages,omitempty"`
	Dependency  map[string]string `json:"dependency,omitempty" yaml:"dependency,omitempty"`
	Description string            `json:"description,omitempty" yaml:"description,omitempty"`
	Files       []ReleaseFile     `json:"files,omitempty" yaml:"files,omitempty"`
	Game        Game              `json:"game,omitempty" yaml:"game,omitempty"`
	License     string            `json:"license,omitempty" yaml:"license,omitempty"`
	Name        string            `json:"name,omitempty" yaml:"name,omitempty"`
	Repo        Repo              `json:"repo,omitempty" yaml:"repo,omitempty"`
	Version     SemVer            `json:"version,omitempty" yaml:"version,omitempty"`
}

func (m *Manifest) DownloadReleaseFiles(ctx context.Context, path string) error {
	for _, f := range m.Files {
		err := f.Download(ctx, path, m)
		if err != nil {
			return err
		}
	}

	return m.writeManifest(path)
}

func (m *Manifest) writeManifest(path string) error {
	content, err := yaml.Marshal(m)
	if err != nil {
		return err
	}

	fname := fmt.Sprintf(
		"%s-%s.%s",
		m.Repo.Name(),
		m.Version.Get().String(),
		manifestName,
	)
	fpath := filepath.Clean(filepath.Join(path, fname))
	return os.WriteFile(fpath, content, 0644)
}

func fetchReleaseFile(ctx context.Context, repo Repo, version *Version, file string) (*resty.Response, error) {
	return resty.
		New().
		R().
		SetDoNotParseResponse(true).
		SetContext(ctx).
		Get(fmt.Sprintf(
			"https://github.com/%s/releases/download/%s/%s",
			repo.String(),
			version.String(),
			file,
		))
}

func FetchManifest(ctx context.Context, repo Repo, version SemVer) (*Manifest, error) {
	resp, err := fetchReleaseFile(ctx, repo, version.Get(), manifestName)
	if err != nil {
		return nil, err
	}
	defer resp.RawBody().Close()

	data, err := ioutil.ReadAll(resp.RawBody())
	if err != nil {
		return nil, err
	}

	mani, ymlErr := ParseYAMLManifest(data)
	if ymlErr != nil {
		var jsonErr error
		mani, jsonErr = ParseJSONManifest(data)
		if jsonErr != nil {
			return nil, fmt.Errorf("invalid syntax: %s : %s", ymlErr, jsonErr)
		}
	}

	mani.Repo = repo
	mani.Version = version
	return mani, nil
}

/*func (m *Manifest) FetchDepManifests(ctx context.Context) (map[string]*Manifest, error) {

	for k := range m.Dependency {
		var d *Manifest
		var err error

		if m.allDeps[k] == nil {
			d, err = FetchManifest(ctx, k, ver)
			if err != nil {
				return nil, err
			}
			m.allDeps[k] = d
		}

		if d != nil && len(d.Dependency) > 0 {
			sub, err := d.FetchDepManifests(ctx)
			if err != nil {
				return nil, err
			}

			for k, v := range sub {
				if m.allDeps[k] == nil {
					m.allDeps[k] = v
				}
			}
		}
	}

	return m.allDeps, nil
}*/

func ParseYAMLManifest(data []byte) (*Manifest, error) {
	var manifest Manifest
	err := yaml.Unmarshal(data, &manifest)
	return &manifest, err
}

func ParseJSONManifest(data []byte) (*Manifest, error) {
	var manifest Manifest
	err := json.Unmarshal(data, &manifest)
	return &manifest, err
}
