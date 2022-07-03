package swizzle

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/google/go-github/v45/github"
)

// ReleaseFile is an archived file for a mod release, hosted on GitHub releases.
type ReleaseFile struct {
	// Release file name.
	Name string `json:"name,omitempty" yaml:"name,omitempty"`

	// Source is the path to mod content inside the release archive. Default is to the
	// root of the archive.
	Source string `json:"source,omitempty" yaml:"source,omitempty"`

	// Destination is the folder path, relative to the game directory, where the
	// mod content should be installed. Default is to the root of the game directory.
	Destination string `json:"destination,omitempty" yaml:"destination,omitempty"`

	archive Archive
	asset   *github.ReleaseAsset
}

// setReleaseAssset matches a ReleaseFile with the set of github release assets.
func (f *ReleaseFile) setReleaseAsset(assets []*github.ReleaseAsset) {
	for _, a := range assets {
		if a.GetName() == f.Name {
			f.asset = a
		}
	}
}

// download fetches and writes the release file to the system at the given folder path.
//
// Downloaded file names are appended with the repo name and the release version to
// prevent naming conflicts with other release files.
func (f *ReleaseFile) download(ctx context.Context, path string, m *Manifest) error {
	if m == nil {
		return fmt.Errorf("nil manifest")
	}

	if f.asset == nil {
		return fmt.Errorf("nil release archive")
	}

	cleanpath := filepath.Clean(path)
	if _, err := os.Stat(cleanpath); err != nil {
		err = os.MkdirAll(cleanpath, 0755)
		if err != nil {
			return err
		}
	}

	resp, err := m.Repo.FetchReleaseAsset(ctx, f.asset)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	f.archive = NewArchive(
		fmt.Sprintf("%s-%s-%s", m.Repo.Name(), m.Version.Get(nil).String(), f.Name),
		cleanpath,
	)
	out, err := os.Create(f.archive.Location())
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = out.ReadFrom(resp.Body)
	if err != nil {
		return err
	}

	return nil
}
