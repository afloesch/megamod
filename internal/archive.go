package internal

import (
	"fmt"
	"path/filepath"
	"strings"
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
	location string
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
	return fmt.Errorf("unknown archive format")
}

// NewArchive returns an Archive object for a file at a given path.
func NewArchive(filename, path string) Archive {
	d := archiveData{
		location: filepath.Clean(filepath.Join(path, filename)),
	}

	if strings.HasSuffix(d.location, string(ZipFileExtension)) {
		return &ZipArchive{d}
	}

	if strings.HasSuffix(d.location, string(SevenZFileExtension)) {
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
