package utils

import (
	"CBCTF/internal/log"
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func Zip(path, zipPath string) error {
	zipFile, err := os.Create(zipPath)
	if err != nil {
		return err
	}
	defer func(zipFile *os.File) {
		if err = zipFile.Close(); err != nil {
			log.Logger.Warningf("Failed to zip files: %s", err)
		}
	}(zipFile)

	zipWriter := zip.NewWriter(zipFile)
	defer func(zipWriter *zip.Writer) {
		if err = zipWriter.Close(); err != nil {
			log.Logger.Warningf("Failed to zip files: %s", err)
		}
	}(zipWriter)
	dir, err := os.ReadDir(path)
	if err != nil {
		return err
	}
	for _, file := range dir {
		err = func(path string) error {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer func(file *os.File) {
				if err = file.Close(); err != nil {
					log.Logger.Warningf("Failed to zip files: %s", err)
				}
			}(file)

			w, err := zipWriter.Create(filepath.Base(path))
			if err != nil {
				return err
			}

			_, err = io.Copy(w, file)
			return err
		}(fmt.Sprintf("%s/%s", path, file.Name()))

		if err != nil {
			return err
		}
	}
	return nil
}
