package swizzle

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const manifestName string = "swiz.zle"

// Game is the supported game and game version for a mod.
type Game struct {
	Executable string `json:"executable,omitempty" yaml:"executable,omitempty"`
	Version    SemVer `json:"version,omitempty" yaml:"version,omitempty"`
}

// Manifest defines the swiz.zle file format for a mod release.
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

// DownloadReleaseFiles fetches and writes all release archives to the system
// at the given folder path.
func (m *Manifest) DownloadReleaseFiles(ctx context.Context, path string) error {
	for _, f := range m.Files {
		err := f.Download(ctx, path, m)
		if err != nil {
			return err
		}
	}

	return m.WriteFile(path)
}

// WriteFile adds the manifest file to the system at the given folder path.
func (m *Manifest) WriteFile(path string) error {
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

/*
ParseManifest unmarshals data to a swizzle manifest. Swizzle manifest
files are either JSON or YAML.
*/
func ParseManifest(data []byte) (*Manifest, error) {
	manifest, ymlErr := parseYAMLManifest(data)
	if ymlErr != nil {
		var jsonErr error
		manifest, jsonErr = parseJSONManifest(data)
		if jsonErr != nil {
			return nil, fmt.Errorf("invalid syntax: %s : %s", ymlErr, jsonErr)
		}
	}

	return manifest, nil
}

// parseYAMLManifest attempts to unmarshal the data from yaml.
func parseYAMLManifest(data []byte) (*Manifest, error) {
	var manifest Manifest
	err := yaml.Unmarshal(data, &manifest)
	return &manifest, err
}

// parseJSONManifest attempts to unmarshal the data from json.
func parseJSONManifest(data []byte) (*Manifest, error) {
	var manifest Manifest
	err := json.Unmarshal(data, &manifest)
	return &manifest, err
}
