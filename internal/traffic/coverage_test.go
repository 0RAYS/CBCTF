package traffic

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestBadCaptureDoesNotDiscardOtherFlagEvidence(t *testing.T) {
	dir := t.TempDir()
	writeTrafficTestCapture(t, filepath.Join(dir, "pod-good.pcap"), "flag{retained}")
	if err := os.WriteFile(filepath.Join(dir, "pod-bad.pcap"), []byte("bad"), 0600); err != nil {
		t.Fatal(err)
	}
	report, err := AnalyzeDir(context.Background(), dir, AnalysisOptions{})
	if err != nil || !report.Partial || !report.Truncated || len(report.SourceIssues) != 1 || len(report.Files) != 2 || len(report.Flags) != 1 || report.Flags[0].Value != "flag{retained}" {
		t.Fatalf("partial evidence lost or mislabeled: %+v %v", report, err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err = AnalyzeDir(ctx, dir, AnalysisOptions{}); !errors.Is(err, context.Canceled) {
		t.Fatalf("cancellation downgraded to file warning: %v", err)
	}
	if _, err = AnalyzeDir(context.Background(), filepath.Join(dir, "missing"), AnalysisOptions{}); err == nil {
		t.Fatal("missing directory downgraded to success")
	}
}
