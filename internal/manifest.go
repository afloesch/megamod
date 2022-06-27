package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	"gopkg.in/yaml.v3"
)

const ESRBERating ESRBString = "E"
const ESRBEPlusRating ESRBString = "E10+"
const ESRBTeenRating ESRBString = "T"
const ESRBMatureRating ESRBString = "M"
const ESRBAdultRating ESRBString = "AO"

type ESRBString string

type Game struct {
	Executable string `json:"executable" yaml:"executable"`
	Version    string `json:"version" yaml:"version"`
}

type AgeRating struct {
	ESRB ESRBString `json:"esrb" yaml:"esrb"`
}

type Directory struct {
	Source string `json:"source" yaml:"source"`
	Target string `json:"target" yaml:"target"`
}

type Manifest struct {
	Dependency map[string]string `json:"dependency" yaml:"dependency"`
	Directory  *Directory        `json:"directory" yaml:"directory"`
	Game       Game              `json:"game" yaml:"game"`
	Name       string            `json:"name" yaml:"name"`
	AgeRating  AgeRating         `json:"rating" yaml:"rating"`
	Schema     string            `json:"schema" yaml:"schema"`
	URL        string            `json:"url" yaml:"url"`
	Version    string            `json:"version" yaml:"version"`

	archive Archive
	allDeps map[string]*Manifest
}

func (m *Manifest) FetchRelease(ctx context.Context, path string) error {
	if m.URL == "" {
		return nil
	}

	cleanpath := filepath.Clean(path)
	if _, err := os.Stat(cleanpath); err != nil {
		err = os.MkdirAll(cleanpath, 0755)
		if err != nil {
			return err
		}
	}

	resp, err := resty.New().R().SetContext(ctx).SetDoNotParseResponse(true).
		Get(m.URL)
	if err != nil {
		return err
	}
	defer resp.RawBody().Close()

	parts := strings.Split(m.URL, "/")
	fname := parts[len(parts)-1]
	location := filepath.Join(cleanpath, fname)
	m.archive = NewArchive(location)

	out, err := os.Create(m.archive.Location())
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = out.ReadFrom(resp.RawBody())
	if err != nil {
		return err
	}

	src := ""
	if m.Directory != nil && m.Directory.Source != "" {
		src = m.Directory.Source
	}

	dest := cleanpath
	if m.Directory != nil && m.Directory.Target != "" {
		dest = filepath.Join(dest, m.Directory.Target)
	}
	return m.archive.Unpack(dest, src)
}

func (m *Manifest) FetchDeps(ctx context.Context, path string) error {
	for _, v := range m.allDeps {
		err := v.FetchRelease(ctx, path)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *Manifest) FlattenDeps(ctx context.Context) ([]*Manifest, error) {
	var deps []*Manifest

	for k := range m.Dependency {
		var d *Manifest
		if m.allDeps[k] == nil {
			d, err := FetchManifest(ctx, k)
			if err != nil {
				return nil, err
			}

			m.allDeps[k] = d
			deps = append(deps, d)
		}

		if d != nil && len(d.Dependency) > 0 {
			sub, err := d.FlattenDeps(ctx)
			if err != nil {
				return nil, err
			}

			deps = append(deps, sub...)
		}
	}

	return deps, nil
}

func FetchManifest(ctx context.Context, repo string) (*Manifest, error) {
	t := time.Now().Unix()
	resp, err := resty.New().R().SetContext(ctx).
		Get(fmt.Sprintf("https://raw.githubusercontent.com/%s?%v", repo, t))
	if err != nil {
		return nil, err
	}

	return ParseYAMLManifest(resp.Body())
}

func ParseManifestFromFile(fpath string) (*Manifest, error) {
	data, err := os.ReadFile(filepath.Clean(fpath))
	if err != nil {
		return nil, err
	}

	man, err := ParseYAMLManifest(data)
	if err != nil {
		return nil, err
	}

	return man, nil
}

func ParseYAMLManifest(data []byte) (*Manifest, error) {
	var manifest Manifest
	err := yaml.Unmarshal(data, &manifest)
	if err != nil {
		return nil, err
	}
	manifest.allDeps = map[string]*Manifest{}
	return &manifest, nil
}

func ParseJSONManifest(data []byte) (*Manifest, error) {
	var manifest Manifest
	err := json.Unmarshal(data, &manifest)
	if err != nil {
		return nil, err
	}
	manifest.allDeps = map[string]*Manifest{}
	return &manifest, nil
}
