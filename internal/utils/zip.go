package utils

import (
	"CBCTF/internal/log"
	"archive/zip"
	"io"
	"os"
	"path/filepath"
)

func Zip(files []string, zipPath string) error {
	zipFile, err := os.Create(zipPath)
	if err != nil {
		log.Logger.Warningf("Failed to zip files: %s", err)
		return err
	}
	defer func(zipFile *os.File) {
		err := zipFile.Close()
		if err != nil {
			log.Logger.Warningf("Failed to zip files: %s", err)
		}
	}(zipFile)

	zipWriter := zip.NewWriter(zipFile)
	defer func(zipWriter *zip.Writer) {
		err := zipWriter.Close()
		if err != nil {
			log.Logger.Warningf("Failed to zip files: %s", err)
		}
	}(zipWriter)

	for _, path := range files {
		err := func(path string) error {
			file, err := os.Open(path)
			if err != nil {
				log.Logger.Warningf("Failed to zip files: %s", err)
				return err
			}
			defer func(file *os.File) {
				err := file.Close()
				if err != nil {
					log.Logger.Warningf("Failed to zip files: %s", err)
				}
			}(file)

			w, err := zipWriter.Create(filepath.Base(path))
			if err != nil {
				log.Logger.Warningf("Failed to zip files: %s", err)
				return err
			}

			_, err = io.Copy(w, file)
			return err
		}(path)

		if err != nil {
			log.Logger.Warningf("Failed to zip files: %s", err)
			return err
		}
	}
	return nil
}
