package swizzle

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
)

// ReleaseFile is an archived file for a mod release, hosted on GitHub releases.
type ReleaseFile struct {
	Filename string `json:"filename,omitempty" yaml:"filename,omitempty"`
	Source   string `json:"source,omitempty" yaml:"source,omitempty"`
	Dest     string `json:"destination,omitempty" yaml:"destination,omitempty"`

	archive Archive `json:"archive,omitempty" yaml:"archive,omitempty"`
}

// Download fetches and writes the release file to system at the given folder path.
func (f *ReleaseFile) Download(ctx context.Context, path string, m *Manifest) error {
	if m == nil {
		return fmt.Errorf("nil manifest")
	}

	cleanpath := filepath.Clean(path)
	if _, err := os.Stat(cleanpath); err != nil {
		err = os.MkdirAll(cleanpath, 0755)
		if err != nil {
			return err
		}
	}

	resp, err := m.Repo.FetchReleaseFile(ctx, m.Version.Get(), f.Filename)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	f.archive = NewArchive(
		fmt.Sprintf("%s-%s-%s", m.Repo.Name(), m.Version.Get().String(), f.Filename),
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
