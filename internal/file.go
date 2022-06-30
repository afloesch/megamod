package internal

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
)

type ReleaseFile struct {
	Filename string `json:"filename,omitempty" yaml:"filename,omitempty"`
	Source   string `json:"source,omitempty" yaml:"source,omitempty"`
	Dest     string `json:"destination,omitempty" yaml:"destination,omitempty"`

	archive Archive `json:"archive,omitempty" yaml:"archive,omitempty"`
}

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

	resp, err := fetchReleaseFile(ctx, m.Repo, m.Version.Get(), f.Filename)
	if err != nil {
		return err
	}
	defer resp.RawBody().Close()

	f.archive = NewArchive(
		fmt.Sprintf("%s-%s-%s", m.Repo.Name(), m.Version.Get().String(), f.Filename),
		cleanpath,
	)
	out, err := os.Create(f.archive.Location())
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = out.ReadFrom(resp.RawBody())
	if err != nil {
		return err
	}

	return nil
}
