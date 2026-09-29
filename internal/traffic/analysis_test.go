package traffic

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/gopacket/gopacket/layers"
)

func TestAnalysisSplitFlagAndKnownValue(t *testing.T) {
	dir := t.TempDir()
	writeTrafficTestCapture(t, filepath.Join(dir, "pod-web.pcap"), "HTTP/1.1 200 OK\r\nContent-Length: 24\r\n\r\nflag{cross_", "packet} secret42")
	report, err := AnalyzeDir(context.Background(), dir, []string{"secret42"})
	if err != nil {
		t.Fatal(err)
	}
	if len(report.Flags) != 2 {
		t.Fatalf("flags: %+v", report)
	}
	if report.Flags[0].Value != "flag{cross_packet}" || !report.Flags[1].Verified {
		t.Fatalf("findings: %+v", report.Flags)
	}
	if report.Flags[0].Evidence.Capture != "pod-web.pcap" || report.Flags[0].Evidence.EndTime.Before(report.Flags[0].Evidence.Time) {
		t.Fatal("missing evidence")
	}
}

func TestContentEncodings(t *testing.T) {
	flag := "flag{encoded_secret}"
	var compressed bytes.Buffer
	w := gzip.NewWriter(&compressed)
	_, _ = w.Write([]byte(flag))
	_ = w.Close()
	cases := map[string][]byte{
		"base64":    []byte(base64.StdEncoding.EncodeToString([]byte(flag))),
		"hex":       []byte(hex.EncodeToString([]byte(flag))),
		"url":       []byte("flag%7Bencoded_secret%7D"),
		"html":      []byte("flag&#123;encoded_secret&#125;"),
		"json":      []byte(`"flag\u007bencoded_secret\u007d"`),
		"gzip":      compressed.Bytes(),
		"http_gzip": append([]byte(fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Encoding: gzip\r\nContent-Length: %d\r\n\r\n", compressed.Len())), compressed.Bytes()...),
		"chunked":   []byte("HTTP/1.1 200 OK\r\nTransfer-Encoding: chunked\r\n\r\n5\r\nflag{\r\nf\r\nencoded_secret}\r\n0\r\n\r\n"),
	}
	for name, data := range cases {
		t.Run(name, func(t *testing.T) {
			a := newAnalyzer(context.Background(), []string{flag})
			a.analyzeContent(data, Evidence{})
			found := false
			for _, finding := range a.report.Flags {
				if finding.Value == flag && finding.Verified {
					found = true
				}
			}
			if !found {
				t.Fatalf("not decoded: %+v", a.report)
			}
		})
	}
}

func TestTCPReorderingRetransmissionAndGap(t *testing.T) {
	a := newAnalyzer(context.Background(), nil)
	s := newStreamCollector(a.analyzeContent, a.warn)
	e := Evidence{Time: time.Unix(1, 0)}
	for _, segment := range []struct {
		seq  uint32
		data string
	}{{105, "secret}"}, {100, "flag{"}, {100, "flag{"}} {
		tcp := &layers.TCP{Seq: segment.seq}
		tcp.Payload = []byte(segment.data)
		s.add(tcp, e)
	}
	s.flush()
	if len(a.report.Flags) != 1 || a.report.Flags[0].Value != "flag{secret}" {
		t.Fatalf("reassembly: %+v", a.report)
	}
	a = newAnalyzer(context.Background(), nil)
	s = newStreamCollector(a.analyzeContent, a.warn)
	first := &layers.TCP{Seq: 100}
	first.Payload = []byte("flag{")
	s.add(first, e)
	last := &layers.TCP{Seq: 110}
	last.Payload = []byte("secret}")
	s.add(last, e)
	s.flush()
	if len(a.report.Flags) != 0 || !a.warnings["tcp_capture_gap"] {
		t.Fatalf("fabricated flag across gap: %+v", a.report)
	}
}

func TestContentDecompressionLimit(t *testing.T) {
	var compressed bytes.Buffer
	w := gzip.NewWriter(&compressed)
	_, _ = w.Write([]byte(strings.Repeat("a", maxStreamBytes+1)))
	_ = w.Close()
	a := newAnalyzer(context.Background(), nil)
	a.analyzeContent(compressed.Bytes(), Evidence{})
	if !a.report.Truncated || !a.warnings["decompression_limit"] {
		t.Fatal("decompression limit not reported")
	}
}

func TestSplitProxyHeader(t *testing.T) {
	path := filepath.Join(t.TempDir(), "frpc.pcap")
	writeTrafficTestCapture(t, path, "PROXY TCP6 2606:4700::", "1111 fd00::2 12345 80\r\n")
	accesses, err := extractFrpcProxyAccesses(context.Background(), path)
	if err != nil || len(accesses) != 1 || accesses[0].IP != "2606:4700::1111" {
		t.Fatalf("split proxy evidence: %+v %v", accesses, err)
	}
}
