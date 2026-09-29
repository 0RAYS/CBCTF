package utils

import (
	"archive/zip"
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestZipAtomicArchiveAndCancellation(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "capture.pcap"), []byte("capture"), 0600); err != nil {
		t.Fatal(err)
	}
	output := filepath.Join(dir, "traffics.zip")
	if err := ZipWithContext(context.Background(), dir, output, nil); err != nil {
		t.Fatal(err)
	}
	archive, err := zip.OpenReader(output)
	if err != nil {
		t.Fatal(err)
	}
	files := 0
	for _, f := range archive.File {
		if !f.FileInfo().IsDir() {
			files++
			if filepath.Base(f.Name) != "capture.pcap" {
				t.Fatal("archive included itself or temp file")
			}
		}
	}
	_ = archive.Close()
	if files != 1 {
		t.Fatalf("files: %d", files)
	}
	before, err := os.ReadFile(output)
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err = ZipWithContext(ctx, dir, output, nil); err == nil {
		t.Fatal("cancel ignored")
	}
	after, err := os.ReadFile(output)
	if err != nil || string(before) != string(after) {
		t.Fatal("cancel corrupted existing archive")
	}
	if err = ZipWithContext(context.Background(), filepath.Join(dir, "missing"), output, nil); err == nil {
		t.Fatal("walk error ignored")
	}
}
