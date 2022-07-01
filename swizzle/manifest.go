package swizzle

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
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
	// AgeRating is a swizzle supported rating system value for content
	// age ratings. An unspecified age rating is assumed safe for all ages.
	AgeRating AgeRating `json:"ages,omitempty" yaml:"ages,omitempty"`

	// Dependency is the optional list of dependent mods for the swizzle manifest.
	Dependency map[Repo]SemVer `json:"dependency,omitempty" yaml:"dependency,omitempty"`

	// An optional short description for the mod.
	Description string `json:"description,omitempty" yaml:"description,omitempty"`

	// All bundled release assets for the mod release. Release files are optional.
	Files []ReleaseFile `json:"files,omitempty" yaml:"files,omitempty"`

	// Game is the supported game and game version for the mod release. All Game values
	// are optional but must be provided to enforce game version compatibility checks
	// between mods, so it can be considered best practice to provide this information.
	//
	// See the Game struct for more information.
	Game Game `json:"game,omitempty" yaml:"game,omitempty"`

	// License is an optional content license for the mod.
	License string `json:"license,omitempty" yaml:"license,omitempty"`

	// Name is optional. Swizzle will default to the repo name if the name
	// field is not set.
	Name string `json:"name,omitempty" yaml:"name,omitempty"`

	// The GitHub repository for the mod release. See the swizzle Repo docs
	// for more information.
	Repo Repo `json:"repo,omitempty" yaml:"repo,omitempty"`

	// Mod version. Must use semantic versioning.
	Version SemVer `json:"version,omitempty" yaml:"version,omitempty"`
}

// AddDependency gets the specified release manifest and adds it and all dependencies
// to the manifest.
func (m *Manifest) AddDependency(ctx context.Context, repo Repo, version SemVer) error {
	dep, err := repo.FetchManifest(ctx, version)
	if err != nil {
		return err
	}

	for k := range dep.Dependency {
		subVer := dep.Dependency[k]
		err = m.addDependency(k, subVer)
		if err != nil {
			return err
		}
	}

	return m.addDependency(repo, version)
}

func (m *Manifest) addDependency(repo Repo, version SemVer) error {

	if m.Dependency == nil {
		m.Dependency = map[Repo]SemVer{}
	}

	/*var exists bool
	for k := range m.Dependency {
		if k == repo {
			exists = true
		}
	}*/

	// temp logic
	m.Dependency[repo] = version

	return nil
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

	fname := fmt.Sprintf(
		"%s-%s.%s",
		m.Repo.Name(),
		m.Version.Get().String(),
		manifestName,
	)
	fpath := filepath.Clean(filepath.Join(path, fname))
	return m.WriteFile(fpath)
}

// WriteFile adds the manifest file to the file system at the given path.
func (m *Manifest) WriteFile(path string) error {
	content, err := yaml.Marshal(m)
	if err != nil {
		return err
	}

	fpath := filepath.Clean(path)
	return os.WriteFile(fpath, content, 0644)
}

// ReadFile parses a manifest file from the file system at the given path.
func (m *Manifest) ReadFile(path string) error {
	b, err := ioutil.ReadFile(filepath.Clean(path))
	if err != nil {
		return err
	}

	parsed, err := ParseManifest(b)
	*m = *parsed
	return err
}

/*
ParseManifest unmarshals data to a swizzle manifest. Swizzle manifest
files are either JSON or YAML, and must be named `swiz.zle`.
*/
func ParseManifest(data []byte) (*Manifest, error) {
	manifest, ymlErr := parseYAMLManifest(data)
	if ymlErr != nil {
		var jsonErr error
		manifest, jsonErr = parseJSONManifest(data)
		if jsonErr != nil {
			return nil, fmt.Errorf("invalid manifest: %s : %s", ymlErr, jsonErr)
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
