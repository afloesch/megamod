package internal

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const ZipFileExtension FileExtension = ".zip"

type ZipArchive struct {
	archiveData
}

func (a ZipArchive) Location() string {
	return a.location
}

func (a ZipArchive) Unpack(dst string, src string) error {
	f, err := zip.OpenReader(a.location)
	if err != nil {
		return err
	}
	defer f.Close()

	for _, f := range f.File {
		filePath := filepath.Clean(filepath.Join(dst, strings.Replace(f.Name, src, "", 1)))
		/*if !strings.HasPrefix(filePath, filepath.Clean(dst)+string(os.PathSeparator)) {
			return fmt.Errorf("invalid file path")
		}*/

		if f.FileInfo().IsDir() {
			os.MkdirAll(filePath, os.ModePerm)
			continue
		}

		if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			return err
		}

		dstFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}

		fileInArchive, err := f.Open()
		if err != nil {
			return err
		}

		if _, err := io.Copy(dstFile, fileInArchive); err != nil {
			return err
		}

		dstFile.Close()
		fileInArchive.Close()
	}

	return nil
}
