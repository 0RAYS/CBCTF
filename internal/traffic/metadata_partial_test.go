package traffic

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/gopacket/gopacket/layers"
)

func TestMalformedProcessRowDoesNotSuppressCapture(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "pod.pcap")
	writeTrafficTestCapture(t, path, "flag{kept}")
	metadata := "bad json\n" + `{"protocol":"TCP","local_addr":"10.0.0.2:80","remote_addr":"8.8.8.8:12345","pid":42,"process_name":"web"}` + "\n{"
	if err := os.WriteFile(path+".connections.jsonl", []byte(metadata), 0600); err != nil {
		t.Fatal(err)
	}
	connections, issues, err := ReadPcapFile(context.Background(), path)
	if err != nil || len(connections) != 1 || len(issues) != 2 || connections[0].Process == nil || connections[0].Process.ProcessName != "web" {
		t.Fatalf("metadata error hid usable traffic: %+v %+v %v", connections, issues, err)
	}
	issues, err = EnrichPcap(context.Background(), path, path+".connections.jsonl", path+".enrich.pcap")
	if err != nil || len(issues) != 2 {
		t.Fatalf("metadata error blocked enrichment: %v %v", issues, err)
	}
}

func TestProxyCollectorDoesNotTrustHeadersAfterCaptureGap(t *testing.T) {
	var chunks []string
	collector := newStreamCollector(func(data []byte, _ Evidence) { chunks = append(chunks, string(data)) }, func(string) {})
	collector.prefixOnly = true
	first := &layers.TCP{Seq: 100}
	first.Payload = []byte("PROXY TCP4 8.8.8.8 1.1.1.1 1 2\r\n")
	later := &layers.TCP{Seq: 1000}
	later.Payload = []byte("PROXY TCP4 9.9.9.9 1.1.1.1 1 2\r\n")
	collector.add(first, Evidence{})
	collector.add(later, Evidence{})
	collector.flush()
	if len(chunks) != 1 || chunks[0] != string(first.Payload) {
		t.Fatalf("trusted an application fragment as proxy evidence: %v", chunks)
	}
}

func TestTruncatedFrpcFileKeepsEarlierCompleteProxyEvidence(t *testing.T) {
	path := filepath.Join(t.TempDir(), "frpc.pcap")
	writeTrafficTestCapture(t, path, "PROXY TCP4 8.8.8.8 1.1.1.1 12345 80\r\n", "request-body")
	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if err = os.Truncate(path, info.Size()-2); err != nil {
		t.Fatal(err)
	}
	accesses, err := extractFrpcProxyAccesses(context.Background(), path, map[uint16]bool{10000: true})
	if err == nil || len(accesses) != 1 || accesses[0].IP != "8.8.8.8" {
		t.Fatalf("good proxy evidence discarded: %+v %v", accesses, err)
	}
}
