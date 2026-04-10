package utils

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
)

func Zip(src string, destZip string) error {
	zipFile, err := os.Create(destZip)
	if err != nil {
		return err
	}
	defer func() {
		_ = zipFile.Close()
	}()
	zipWriter := zip.NewWriter(zipFile)
	defer func() {
		_ = zipWriter.Close()
	}()

	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if filepath.Clean(path) == filepath.Clean(destZip) {
			return nil
		}
		relPath, err := filepath.Rel(filepath.Dir(src), path)
		if err != nil {
			return err
		}
		if relPath == "." {
			return nil
		}
		if info.IsDir() {
			_, err = zipWriter.Create(relPath + "/")
			return err
		}
		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer func(f *os.File) {
			_ = f.Close()
		}(f)
		w, err := zipWriter.Create(relPath)
		if err != nil {
			return err
		}
		_, err = io.Copy(w, f)
		return err
	})
}
