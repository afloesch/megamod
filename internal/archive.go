package internal

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"
)

// TarballFileExtension is the file extension for a tarball.
const TarballFileExtension FileExtension = ".tar.gz"

// RarFileExtension is the file extension for a Roshal archive.
const RarFileExtension FileExtension = ".rar"

// UnknownFileExtension is for an unknown archive file format.
const UnknownFileExtension FileExtension = ""

// FileExtension is a file extension string.
type FileExtension string

// archiveData is common to any archive format.
type archiveData struct {
	location  string
	extension FileExtension
}

// Archive defines the interface for working with different.
// file archive formats
type Archive interface {
	Location() string
	Unpack(dst string, src string) error
}

// UnknownArchive is an Archive format which cannot be identified.
type UnknownArchive struct {
	archiveData
}

func (a UnknownArchive) Location() string {
	return a.location
}

func (a UnknownArchive) Unpack(dst, src string) error {
	return nil
}

// NewArchive returns an Archive object for a file at a given path.
func NewArchive(path string) Archive {
	d := archiveData{
		location:  filepath.Clean(path),
		extension: UnknownFileExtension,
	}

	if strings.HasSuffix(d.location, string(ZipFileExtension)) {
		d.extension = ZipFileExtension
		return &ZipArchive{d}
	}

	if strings.HasSuffix(d.location, string(SevenZFileExtension)) {
		d.extension = SevenZFileExtension
		return &SevenZArchive{d}
	}

	/*if strings.HasSuffix(d.location, string(TarballFileExtension)) {
		d.extension = TarballFileExtension
	}

	if strings.HasSuffix(d.location, string(RarFileExtension)) {
		d.extension = RarFileExtension
	}*/

	return &UnknownArchive{d}
}

func getFilename(url string) string {
	parts := strings.Split(url, "/")
	return parts[len(parts)-1]
}

func sanitizeName(name string) string {
	clean := strings.ReplaceAll(name, " ", "-")
	clean = strings.ReplaceAll(clean, "/", "-")
	return strings.ReplaceAll(clean, "\\", "-")
}

func GetArchivePath(name, path, url string) string {
	return filepath.Join(path, fmt.Sprintf(
		"%v-%s-%s",
		time.Now().Unix(),
		sanitizeName(name),
		getFilename(url),
	))
}
