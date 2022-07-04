package swizzle

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/afloesch/semver"
	"github.com/google/go-github/v45/github"
	"gopkg.in/yaml.v3"
)

const manifestName string = "swiz.zle"

// Game is the supported game and game version for a mod release. Specifying
// the game values are optional, but encouraged to enable version and game
// compatibility checks between mods.
type Game struct {
	// Executable is the supported game executable file name, including the
	// file extension.
	Executable string `json:"executable,omitempty" yaml:"executable,omitempty"`

	// The supported game version/(s) for the mod based on a semver semantic
	// string format. An empty version is equivalent to >=v0.0.0 and so will
	// match all game versions.
	Version semver.String `json:"version,omitempty" yaml:"version,omitempty"`
}

// Manifest defines the swiz.zle file format for a mod release. All mods released
// for swizzle host mod releases publicly on GitHub, and include the release swiz.zle
// file in the GitHub release assets.
type Manifest struct {
	// AgeRating is a swizzle supported rating system value for content
	// age ratings. An unspecified age rating is assumed safe for all ages.
	AgeRating AgeRating `json:"ages,omitempty" yaml:"ages,omitempty"`

	// Dependency is the optional list of dependent mods for the swizzle manifest.
	Dependency map[Repo]semver.String `json:"dependency,omitempty" yaml:"dependency,omitempty"`

	// An optional short description for the mod.
	Description string `json:"description,omitempty" yaml:"description,omitempty"`

	// All bundled release assets for the mod release. Release files are optional.
	Files []*ReleaseFile `json:"files,omitempty" yaml:"files,omitempty"`

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
	Version semver.String `json:"version,omitempty" yaml:"version,omitempty"`

	release      *github.RepositoryRelease
	releaseAsset *github.ReleaseAsset
}

// New creates an empty swizzle manifest.
func New() *Manifest {
	return &Manifest{
		Dependency: map[Repo]semver.String{},
	}
}

func (m *Manifest) SetGame(executable, version string) *Manifest {
	m.Game = Game{
		Executable: executable,
		Version:    semver.String(version),
	}
	return m
}

func (m *Manifest) SetRepo(repo string) *Manifest {
	m.Repo = Repo(repo)
	return m
}

func (m *Manifest) SetVersion(version string) *Manifest {
	m.Version = semver.String(version)
	return m
}

// AddDependency gets the specified release manifest and adds it and all dependencies
// to the manifest.
func (m *Manifest) AddDependency(ctx context.Context, repo string, version string) error {
	r := Repo(repo)
	v := semver.String(version).Get()
	rel, err := r.Release(ctx, version)
	if err != nil {
		return err
	}

	dep, err := r.Manifest(ctx, rel)
	if err != nil {
		return err
	}

	for k := range dep.Dependency {
		subVer := dep.Dependency[k]
		err = m.addDependency(k, subVer.Get())
		if err != nil {
			return err
		}
	}

	return m.addDependency(r, v)
}

func (m *Manifest) addDependency(repo Repo, version *semver.Version) error {
	if m.Dependency == nil {
		m.Dependency = map[Repo]semver.String{}
	}

	var exists bool
	for k := range m.Dependency {
		if k == repo {
			exists = true
			currVer := m.Dependency[repo].Get()
			if ok := currVer.OpCompare(version); ok {
				m.Dependency[repo] = version.ToString()
			} else {
				return fmt.Errorf(
					"'%s' version '%s' is incompatible with '%s'",
					repo.String(),
					currVer.String(),
					version.String(),
				)
			}
		}
	}

	if !exists {
		m.Dependency[repo] = version.ToString()
	}

	return nil
}

// DownloadReleaseFile fetches and writes all release archives to the system
// at the given folder path.
func (m *Manifest) DownloadReleaseFile(
	ctx context.Context,
	file *ReleaseFile,
	path string,
) (chan bool, chan float64, chan error) {
	done := make(chan bool)
	progCh := make(chan float64)
	errCh := make(chan error)
	go file.download(ctx, path, m, done, progCh, errCh)
	return done, progCh, errCh

	/*var prog float64
	var err error
	for {
		select {
		case <-done:
			fmt.Println()
			return nil
		case prog = <-progCh:
			fmt.Printf(fmt.Sprintf("\r%v percent of %v", math.Floor(prog), file.Size()))
		case err = <-errCh:
			fmt.Println()
			return err
		}
	}*/

	/*fname := fmt.Sprintf(
		"%s-%s.%s",
		m.Repo.Name(),
		m.Version.Get().String(),
		manifestName,
	)
	fpath := filepath.Clean(filepath.Join(path, fname))
	return m.WriteFile(fpath)*/
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
func (m *Manifest) ReadFile(path string) (*Manifest, error) {
	b, err := ioutil.ReadFile(filepath.Clean(path))
	if err != nil {
		return m, err
	}

	parsed, err := ParseManifest(b)
	*m = *parsed
	return m, err
}

/*
ParseManifest unmarshals data to a swizzle manifest. Swizzle manifest
files are either JSON or YAML, and must be named "swiz.zle".
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

	if manifest.Dependency == nil {
		manifest.Dependency = map[Repo]semver.String{}
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
