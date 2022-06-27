package internal

import (
	"path/filepath"
	"strings"
)

const TarballArchiveType ArchiveType = ".tar.gz"
const RarArchiveType ArchiveType = ".rar"
const UnknownArchiveType ArchiveType = ""

type ArchiveType string

type archiveData struct {
	name      string
	location  string
	extension ArchiveType
}

type Archive interface {
	Location() string
	Unpack(dst string) error
}

type UnknownArchive struct {
	archiveData
}

func (a UnknownArchive) Location() string {
	return a.location
}

func (a UnknownArchive) Unpack(dst string) error {
	return nil
}

func NewArchive(path string) Archive {
	d := archiveData{
		location:  filepath.Clean(path),
		extension: UnknownArchiveType,
	}

	if strings.HasSuffix(d.location, string(ZipArchiveType)) {
		d.extension = ZipArchiveType
		return &ZipArchive{d}
	}

	if strings.HasSuffix(d.location, string(SevenZArchiveType)) {
		d.extension = SevenZArchiveType
		return &SevenZArchive{d}
	}

	/*if strings.HasSuffix(d.location, string(TarballArchiveType)) {
		d.extension = TarballArchiveType
	}

	if strings.HasSuffix(d.location, string(RarArchiveType)) {
		d.extension = RarArchiveType
	}*/

	return &UnknownArchive{d}
}
