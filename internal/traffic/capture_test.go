package traffic

import (
	"context"
	"net"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/gopacket/gopacket"
	"github.com/gopacket/gopacket/layers"
	"github.com/gopacket/gopacket/pcapgo"
)

func TestTrafficAddressPolicy(t *testing.T) {
	for _, ip := range []string{"127.0.0.1", "0.1.2.3", "10.0.0.1", "100.64.0.1", "169.254.1.1", "224.0.0.1", "255.255.255.255", "198.18.1.1", "192.0.2.1", "::", "::1", "fe80::1%eth0", "fd00::1", "ff02::1", "2001:db8::1", "64:ff9b::a00:1", "2002:a00:1::", "fec0::1", "::ffff:127.0.0.1"} {
		if IsPublicTrafficIP(ip) {
			t.Errorf("special IP accepted: %s", ip)
		}
	}
	for _, ip := range []string{"8.8.8.8", "::ffff:8.8.8.8", "2606:4700:4700::1111"} {
		if !IsPublicTrafficIP(ip) {
			t.Errorf("public IP rejected: %s", ip)
		}
	}
	if NormalizeTrafficIP("::ffff:8.8.8.8") != "8.8.8.8" {
		t.Fatal("mapped IP not normalized")
	}
	if isIgnoredTrafficIP("10.0.0.2") || isIgnoredTrafficIP("fd00::2") {
		t.Fatal("target private IP filtered from topology")
	}
	if !isIgnoredTrafficIP("fe80::1") {
		t.Fatal("link local not filtered")
	}
}

func TestCollectTrafficAccesses(t *testing.T) {
	at := time.Unix(1700000000, 0)
	result := PcapDirResult{Connections: []Connection{
		{SrcIP: "8.8.8.8", DstIP: "10.0.0.2", Type: "TCP", SYN: true, Time: at},
		{SrcIP: "1.1.1.1", DstIP: "10.0.0.2", Type: "TCP", SYN: true, ACK: true, Time: at},
		{SrcIP: "10.0.0.3", DstIP: "10.0.0.2", Type: "TCP", SYN: true, Time: at},
	}, Accesses: []TrafficAccess{{IP: "::ffff:8.8.8.8", Time: at, Source: "tcp_syn"}, {IP: "2606:4700::1111", Time: at, Source: "proxy_protocol"}}}
	got := CollectTrafficAccesses(result, map[string]bool{"10.0.0.2": true})
	if len(got) != 2 || got[0].IP != "8.8.8.8" {
		t.Fatalf("unexpected accesses: %+v", got)
	}
}

func writeTrafficTestCapture(t *testing.T, path string, payloads ...string) {
	t.Helper()
	f, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	w := pcapgo.NewWriter(f)
	if err = w.WriteFileHeader(65535, layers.LinkTypeEthernet); err != nil {
		t.Fatal(err)
	}
	seq := uint32(100)
	for i, payload := range payloads {
		ip := &layers.IPv4{Version: 4, TTL: 64, Protocol: layers.IPProtocolTCP, SrcIP: net.ParseIP("8.8.8.8"), DstIP: net.ParseIP("10.0.0.2")}
		tcp := &layers.TCP{SrcPort: 12345, DstPort: 80, Seq: seq, ACK: true}
		if filepath.Base(path) == "frpc.pcap" {
			ip.SrcIP, ip.DstIP = net.ParseIP("127.0.0.1"), net.ParseIP("127.0.0.1")
			tcp.DstPort = 10000
		}
		if err = tcp.SetNetworkLayerForChecksum(ip); err != nil {
			t.Fatal(err)
		}
		buf := gopacket.NewSerializeBuffer()
		err = gopacket.SerializeLayers(buf, gopacket.SerializeOptions{FixLengths: true, ComputeChecksums: true}, &layers.Ethernet{SrcMAC: net.HardwareAddr{1, 2, 3, 4, 5, 6}, DstMAC: net.HardwareAddr{6, 5, 4, 3, 2, 1}, EthernetType: layers.EthernetTypeIPv4}, ip, tcp, gopacket.Payload(payload))
		if err != nil {
			t.Fatal(err)
		}
		data := buf.Bytes()
		if err = w.WritePacket(gopacket.CaptureInfo{Timestamp: time.Unix(1700000000, int64(i)*1000000), CaptureLength: len(data), Length: len(data)}, data); err != nil {
			t.Fatal(err)
		}
		seq += uint32(len(payload))
	}
}

func TestTrafficReaderErrorsAndEnrichment(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "pod-test.pcap")
	writeTrafficTestCapture(t, path, "hello")
	if err := EnrichPcap(context.Background(), path, path+".connections.jsonl", path+".enrich.pcap"); err != nil {
		t.Fatal(err)
	}
	result, err := ReadPcapDir(context.Background(), dir, nil)
	if err != nil || len(result.Connections) != 1 {
		t.Fatalf("duplicate enriched capture: %d, %v", len(result.Connections), err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err = ReadPcapFile(ctx, path); err == nil {
		t.Fatal("cancellation ignored")
	}
	if err = os.WriteFile(filepath.Join(dir, "broken.pcap"), []byte("broken"), 0600); err != nil {
		t.Fatal(err)
	}
	partial, err := ReadPcapDir(context.Background(), dir, nil)
	if err != nil || len(partial.SourceIssues) != 1 || len(partial.Connections) != 1 {
		t.Fatalf("capture failure must retain valid traffic and report coverage: %+v %v", partial, err)
	}
}

func TestProxyAccessIPv6(t *testing.T) {
	path := filepath.Join(t.TempDir(), "frpc.pcap")
	writeTrafficTestCapture(t, path, "PROXY TCP6 2606:4700::1111 fd00::2 12345 80\r\n")
	got, err := extractFrpcProxyAccesses(context.Background(), path, map[uint16]bool{10000: true})
	if err != nil || len(got) != 1 || got[0].IP != "2606:4700::1111" {
		t.Fatalf("IPv6 access: %+v %v", got, err)
	}
}
