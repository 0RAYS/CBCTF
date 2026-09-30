package traffic

import (
	"archive/zip"
	"context"
	"io"
	"os"
	"path/filepath"
	"testing"
)

func TestArchiveKeepsRawEvidenceWhenEnrichmentFails(t *testing.T) {
	dir := t.TempDir()
	writeTrafficTestCapture(t, filepath.Join(dir, "good.pcap"), "flag{evidence}")
	if err := os.WriteFile(filepath.Join(dir, "bad.pcap"), []byte("unreadable capture bytes"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "bad.pcap.enrich.pcap"), []byte("stale derived bytes"), 0600); err != nil {
		t.Fatal(err)
	}
	output := filepath.Join(dir, "traffics.zip")
	issues, err := Archive(context.Background(), dir, output)
	if err != nil || len(issues) != 1 {
		t.Fatalf("enrichment failure blocked raw archive: %v %v", issues, err)
	}
	archive, err := zip.OpenReader(output)
	if err != nil {
		t.Fatal(err)
	}
	defer func(archive *zip.ReadCloser) {
		_ = archive.Close()
	}(archive)
	names := make(map[string]bool)
	for _, entry := range archive.File {
		name := filepath.Base(entry.Name)
		names[name] = true
		if name == "bad.pcap" {
			r, err := entry.Open()
			if err != nil {
				t.Fatal(err)
			}
			data, err := io.ReadAll(r)
			_ = r.Close()
			if err != nil || string(data) != "unreadable capture bytes" {
				t.Fatal("original evidence was modified")
			}
		}
	}
	if !names["bad.pcap"] || !names["good.pcap"] || !names["good.pcap.enrich.pcap"] || names["bad.pcap.enrich.pcap"] {
		t.Fatalf("wrong archive selection: %v", names)
	}
}
