package utils

import (
	"CBCTF/internal/log"
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
		if cerr := zipFile.Close(); cerr != nil {
			log.Logger.Warningf("Failed to close zip file: %s", cerr)
		}
	}()
	zipWriter := zip.NewWriter(zipFile)
	defer func() {
		if cerr := zipWriter.Close(); cerr != nil {
			log.Logger.Warningf("Failed to close zip writer: %s", cerr)
		}
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
			if cerr := f.Close(); cerr != nil {
				log.Logger.Warningf("Failed to close zip file: %s", cerr)
			}
		}(f)
		w, err := zipWriter.Create(relPath)
		if err != nil {
			return err
		}
		_, err = io.Copy(w, f)
		return err
	})
}
