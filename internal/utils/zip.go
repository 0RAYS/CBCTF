package utils

import (
	"archive/zip"
	"context"
	"io"
	"os"
	"path/filepath"
)

func Zip(src string, destZip string) error {
	return ZipWithContext(context.Background(), src, destZip, nil)
}

func ZipWithContext(ctx context.Context, src string, destZip string, include func(string) bool) error {
	zipFile, err := os.CreateTemp(filepath.Dir(destZip), ".archive-*")
	if err != nil {
		return err
	}
	defer func() {
		_ = zipFile.Close()
		_ = os.Remove(zipFile.Name())
	}()
	zipWriter := zip.NewWriter(zipFile)
	defer func() {
		_ = zipWriter.Close()
	}()

	err = filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if err = ctx.Err(); err != nil {
			return err
		}
		if filepath.Clean(path) == filepath.Clean(destZip) || filepath.Clean(path) == filepath.Clean(zipFile.Name()) {
			return nil
		}
		if !info.IsDir() && include != nil && !include(path) {
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
		_, err = io.Copy(w, contextReader{ctx: ctx, reader: f})
		return err
	})
	if err != nil {
		return err
	}
	if err = zipWriter.Close(); err != nil {
		return err
	}
	if err = zipFile.Close(); err != nil {
		return err
	}
	return os.Rename(zipFile.Name(), destZip)
}

type contextReader struct {
	ctx    context.Context
	reader io.Reader
}

func (r contextReader) Read(p []byte) (int, error) {
	if err := r.ctx.Err(); err != nil {
		return 0, err
	}
	return r.reader.Read(p)
}
