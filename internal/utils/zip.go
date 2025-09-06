package utils

import (
	"CBCTF/internal/log"
	"archive/zip"
	"io"
	"os"
	"path/filepath"
)

func Zip(path, zipPath string) error {
	entries, err := os.ReadDir(path)
	if err != nil {
		return err
	}
	zipFile, err := os.Create(zipPath)
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
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		filePath := filepath.Join(path, entry.Name())

		if filePath == zipPath {
			continue
		}
		err = func(filePath string) error {
			f, err := os.Open(filePath)
			if err != nil {
				return err
			}
			defer func() {
				if cerr := f.Close(); cerr != nil {
					log.Logger.Warningf("Failed to close file %s: %s", filePath, cerr)
				}
			}()
			w, err := zipWriter.Create(entry.Name())
			if err != nil {
				return err
			}
			_, err = io.Copy(w, f)
			return err
		}(filePath)
		if err != nil {
			return err
		}
	}
	return nil
}
