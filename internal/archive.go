package internal

import (
	"path/filepath"
	"strings"
)

const TarballFileExtension FileExtension = ".tar.gz"
const RarFileExtension FileExtension = ".rar"
const UnknownFileExtension FileExtension = ""

type FileExtension string

type archiveData struct {
	location  string
	extension FileExtension
}

type Archive interface {
	Location() string
	Unpack(dst string, src string) error
}

type UnknownArchive struct {
	archiveData
}

func (a UnknownArchive) Location() string {
	return a.location
}

func (a UnknownArchive) Unpack(dst, src string) error {
	return nil
}

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
