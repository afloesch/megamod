package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/go-resty/resty/v2"
	"gopkg.in/yaml.v3"
)

const yamlFileExt string = ".mm.yml"

const ESRBERating ESRBString = "E"
const ESRBEPlusRating ESRBString = "E10+"
const ESRBTeenRating ESRBString = "T"
const ESRBMatureRating ESRBString = "M"
const ESRBAdultRating ESRBString = "AO"

type ESRBString string

type Game struct {
	Executable string `json:"executable,omitempty" yaml:"executable,omitempty"`
	Version    string `json:"version,omitempty" yaml:"version,omitempty"`
}

type AgeRating struct {
	ESRB ESRBString `json:"esrb,omitempty" yaml:"esrb,omitempty"`
}

type Directory struct {
	Source string `json:"source,omitempty" yaml:"source,omitempty"`
	Target string `json:"target,omitempty" yaml:"target,omitempty"`
}

type Manifest struct {
	Dependency map[string]string `json:"dependency,omitempty" yaml:"dependency,omitempty"`
	Directory  *Directory        `json:"directory,omitempty" yaml:"directory,omitempty"`
	Game       Game              `json:"game,omitempty" yaml:"game,omitempty"`
	Name       string            `json:"name,omitempty" yaml:"name,omitempty"`
	AgeRating  AgeRating         `json:"rating,omitempty" yaml:"rating,omitempty"`
	URL        string            `json:"url" yaml:"url"`
	Version    string            `json:"version" yaml:"version"`

	archive Archive
	allDeps map[string]*Manifest
}

func (m *Manifest) AllDependencies() map[string]*Manifest {
	return m.allDeps
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

	resp, err := resty.New().R().
		SetContext(ctx).
		SetDoNotParseResponse(true).
		Get(m.URL)
	if err != nil {
		return err
	}
	defer resp.RawBody().Close()

	m.archive = NewArchive(GetArchivePath(m.Name, cleanpath, m.URL))
	out, err := os.Create(m.archive.Location())
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = out.ReadFrom(resp.RawBody())
	if err != nil {
		return err
	}

	return m.writeManifest()

	/*src := ""
	if m.Directory != nil && m.Directory.Source != "" {
		src = m.Directory.Source
	}

	dest := cleanpath
	if m.Directory != nil && m.Directory.Target != "" {
		dest = filepath.Join(dest, m.Directory.Target)
	}
	return m.archive.Unpack(dest, src)*/
}

func (m *Manifest) writeManifest() error {
	content, err := yaml.Marshal(m)
	if err != nil {
		return err
	}

	return os.WriteFile(
		fmt.Sprintf("%s%s", m.archive.Location(), yamlFileExt),
		content,
		0644,
	)
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

func (m *Manifest) FetchDepManifests(ctx context.Context) (map[string]*Manifest, error) {

	for k := range m.Dependency {
		var d *Manifest
		var err error

		if m.allDeps[k] == nil {
			d, err = FetchManifest(ctx, k)
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
