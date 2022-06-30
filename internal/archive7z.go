package internal

import (
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/bodgit/sevenzip"
)

const SevenZFileExtension FileExtension = ".7z"

type SevenZArchive struct {
	archiveData
}

func (a SevenZArchive) Location() string {
	return a.location
}

func (a SevenZArchive) Unpack(dst string, src string) error {
	f, err := sevenzip.OpenReader(a.location)
	if err != nil {
		return err
	}

	for _, file := range f.File {
		filePath := filepath.Clean(filepath.Join(dst, strings.Replace(file.Name, src, "", 1)))
		/*if !strings.HasPrefix(filePath, filepath.Clean(dst)+string(os.PathSeparator)) {
			return fmt.Errorf("invalid file path")
		}*/

		if file.FileInfo().IsDir() {
			os.MkdirAll(filePath, os.ModePerm)
			continue
		}

		if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			return err
		}

		dstFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return err
		}

		fileInArchive, err := file.Open()
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
