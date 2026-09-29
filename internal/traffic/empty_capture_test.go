package traffic

import (
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/gopacket/gopacket/layers"
	"github.com/gopacket/gopacket/pcapgo"
)

func writeEmptyTrafficCaptures(t *testing.T, dir string) []string {
	t.Helper()
	names := []string{"pod-empty.pcap", "pod-empty.pcapng", "frpc.pcap", "pod-header.pcap", "pod-header.pcapng"}
	for _, name := range names[:3] {
		if err := os.WriteFile(filepath.Join(dir, name), nil, 0600); err != nil {
			t.Fatal(err)
		}
	}
	writeTrafficTestCapture(t, filepath.Join(dir, names[3]))
	f, err := os.Create(filepath.Join(dir, names[4]))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	w, err := pcapgo.NewNgWriter(f, layers.LinkTypeEthernet)
	if err != nil {
		t.Fatal(err)
	}
	if err = w.Flush(); err != nil {
		t.Fatal(err)
	}
	return names
}

func TestEmptyCapturesAcrossAnalysisReplayAndEnrichment(t *testing.T) {
	dir := t.TempDir()
	names := writeEmptyTrafficCaptures(t, dir)
	ctx := context.Background()
	for _, name := range names {
		connections, _, err := ReadPcapFile(ctx, filepath.Join(dir, name))
		if err != nil || len(connections) != 0 {
			t.Fatalf("empty capture %s: %d connections, %v", name, len(connections), err)
		}
	}
	report, err := AnalyzeDir(ctx, dir, AnalysisOptions{})
	if err != nil || report.Packets != 0 || len(report.Flags) != 0 || report.Truncated || len(report.Warnings) != 0 {
		t.Fatalf("empty directory report: %+v, %v", report, err)
	}
	result, err := ReadPcapDir(ctx, dir, map[uint16]bool{10000: true})
	if err != nil || len(result.Connections) != 0 || len(result.Accesses) != 0 {
		t.Fatalf("empty replay/proxy evidence: %+v, %v", result, err)
	}
	if errs := EnrichPcapDirWithContext(ctx, dir); len(errs) != 0 {
		t.Fatalf("empty capture enrichment: %v", errs)
	}
	for _, name := range names {
		connections, _, err := ReadPcapFile(ctx, filepath.Join(dir, name+".enrich.pcap"))
		if err != nil || len(connections) != 0 {
			t.Fatalf("invalid empty enriched capture %s: %v", name, err)
		}
	}
	info, err := os.Stat(filepath.Join(dir, names[0]))
	if err != nil || info.Size() != 0 {
		t.Fatalf("original empty capture changed: %v", err)
	}
}

func TestEmptyCapturesDoNotHideValidTraffic(t *testing.T) {
	dir := t.TempDir()
	writeEmptyTrafficCaptures(t, dir)
	writeTrafficTestCapture(t, filepath.Join(dir, "pod-active.pcap"), "flag{still_analyzed}")
	report, err := AnalyzeDir(context.Background(), dir, AnalysisOptions{})
	if err != nil || report.Packets != 1 || len(report.Flags) != 1 || report.Flags[0].Value != "flag{still_analyzed}" {
		t.Fatalf("valid traffic skipped: %+v, %v", report, err)
	}
	result, err := ReadPcapDir(context.Background(), dir, map[uint16]bool{10000: true})
	if err != nil || len(result.Connections) != 1 {
		t.Fatalf("valid replay traffic skipped: %+v, %v", result, err)
	}
}

func TestEmptyCaptureHandlingPreservesCancellationAndHeaderErrors(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "pod.pcap")
	if err := os.WriteFile(path, nil, 0600); err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, _, err := ReadPcapFile(ctx, path); !errors.Is(err, context.Canceled) {
		t.Fatalf("empty capture swallowed cancellation: %v", err)
	}
	for _, data := range [][]byte{{0xd4}, {0xd4, 0xc3, 0xb2}, {0xd4, 0xc3, 0xb2, 0xa1}, {10, 13, 13, 10}, []byte("broken capture")} {
		if err := os.WriteFile(path, data, 0600); err != nil {
			t.Fatal(err)
		}
		if _, _, err := ReadPcapFile(context.Background(), path); err == nil {
			t.Fatalf("partial or corrupt header accepted as empty: %x", data)
		} else if len(data) < 4 && !errors.Is(err, io.ErrUnexpectedEOF) {
			t.Fatalf("partial magic should report truncation: %v", err)
		}
		report, err := AnalyzeDir(context.Background(), dir, AnalysisOptions{})
		if err != nil || !report.Partial || len(report.SourceIssues) != 1 {
			t.Fatalf("file error lost from partial report: %+v %v", report, err)
		}
	}
}
