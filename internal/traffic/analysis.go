package traffic

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/gopacket/gopacket"
	"github.com/gopacket/gopacket/layers"
)

const (
	maxStreamBytes   = 1 << 20
	maxBufferedBytes = 32 << 20
	maxStreams       = 4096
	maxFindings      = 2000
)

type Evidence struct {
	Capture  string    `json:"capture"`
	Time     time.Time `json:"time"`
	EndTime  time.Time `json:"end_time"`
	SrcIP    string    `json:"src_ip"`
	DstIP    string    `json:"dst_ip"`
	SrcPort  string    `json:"src_port"`
	DstPort  string    `json:"dst_port"`
	Protocol string    `json:"protocol"`
}

type FlagFinding struct {
	Value    string   `json:"value"`
	Encoding string   `json:"encoding"`
	Verified bool     `json:"verified"`
	Evidence Evidence `json:"evidence"`
}

type HTTPRecord struct {
	Method      string   `json:"method,omitempty"`
	Host        string   `json:"host,omitempty"`
	URI         string   `json:"uri,omitempty"`
	Status      int      `json:"status,omitempty"`
	ContentType string   `json:"content_type,omitempty"`
	Evidence    Evidence `json:"evidence"`
}

type DNSRecord struct {
	Name     string   `json:"name"`
	Type     string   `json:"type"`
	Response bool     `json:"response"`
	Evidence Evidence `json:"evidence"`
}

type AnalysisReport struct {
	Flags     []FlagFinding `json:"flags"`
	HTTP      []HTTPRecord  `json:"http"`
	DNS       []DNSRecord   `json:"dns"`
	Warnings  []string      `json:"warnings"`
	Truncated bool          `json:"truncated"`
	Packets   int64         `json:"packets"`
	Bytes     int64         `json:"bytes"`
}

type analyzer struct {
	report    AnalysisReport
	known     []string
	seenFlags map[string]bool
	warnings  map[string]bool
	ctx       context.Context
}

func newAnalyzer(ctx context.Context, known []string) *analyzer {
	return &analyzer{ctx: ctx, known: known, seenFlags: make(map[string]bool), warnings: make(map[string]bool),
		report: AnalysisReport{Flags: []FlagFinding{}, HTTP: []HTTPRecord{}, DNS: []DNSRecord{}, Warnings: []string{}}}
}

func (a *analyzer) warn(code string) {
	if !a.warnings[code] {
		a.warnings[code] = true
		a.report.Warnings = append(a.report.Warnings, code)
	}
}

func packetEvidence(c Connection) Evidence {
	return Evidence{Capture: c.Capture, Time: c.Time, EndTime: c.Time, SrcIP: c.SrcIP, DstIP: c.DstIP,
		SrcPort: c.SrcPort, DstPort: c.DstPort, Protocol: c.Type}
}

// AnalyzeDir analyzes original captures only. Every decoded finding carries its
// capture, directional endpoints, and the time range of the contributing stream.
func AnalyzeDir(ctx context.Context, path string, knownFlags []string) (AnalysisReport, error) {
	files, err := os.ReadDir(path)
	if err != nil {
		return AnalysisReport{}, err
	}
	a := newAnalyzer(ctx, knownFlags)
	for _, file := range files {
		if file.IsDir() || !isOriginalTrafficCapture(file.Name()) || file.Name() == frpcPcapName {
			continue
		}
		streams := newStreamCollector(func(data []byte, evidence Evidence) { a.analyzeContent(data, evidence) }, func(code string) { a.report.Truncated = true; a.warn(code) })
		err = walkTrafficPackets(ctx, filepath.Join(path, file.Name()), func(packet gopacket.Packet, _ layers.LinkType) error {
			c, ok := extractTrafficConnection(packet, nil)
			if !ok {
				return nil
			}
			c.Capture = file.Name()
			evidence := packetEvidence(c)
			a.report.Packets++
			a.report.Bytes += int64(c.Size)
			if packet.Metadata().CaptureLength < packet.Metadata().Length {
				a.warn("capture_snaplen_truncated")
			}
			if ipv4, ok := packet.Layer(layers.LayerTypeIPv4).(*layers.IPv4); ok && (ipv4.FragOffset != 0 || ipv4.Flags&layers.IPv4MoreFragments != 0) {
				a.warn("ip_fragments_not_reassembled")
			}
			if packet.Layer(layers.LayerTypeIPv6Fragment) != nil {
				a.warn("ip_fragments_not_reassembled")
			}
			if tcp, ok := packet.Layer(layers.LayerTypeTCP).(*layers.TCP); ok {
				streams.add(tcp, evidence)
			} else if transport := packet.TransportLayer(); transport != nil {
				a.scanEncoded(transport.LayerPayload(), evidence, "raw")
			} else if network := packet.NetworkLayer(); network != nil {
				a.scanEncoded(network.LayerPayload(), evidence, "raw")
			}
			if dns, ok := packet.Layer(layers.LayerTypeDNS).(*layers.DNS); ok {
				for _, question := range dns.Questions {
					if len(a.report.DNS) >= maxFindings {
						a.report.Truncated = true
						a.warn("dns_limit")
						break
					}
					a.report.DNS = append(a.report.DNS, DNSRecord{Name: string(question.Name), Type: question.Type.String(), Response: dns.QR, Evidence: evidence})
					a.scanEncoded(question.Name, evidence, "dns")
				}
				for _, answer := range dns.Answers {
					for _, txt := range answer.TXTs {
						a.scanEncoded(txt, evidence, "dns_txt")
					}
				}
			}
			return ctx.Err()
		})
		if err != nil {
			return AnalysisReport{}, fmt.Errorf("analyze %s: %w", file.Name(), err)
		}
		streams.flush()
		if err = ctx.Err(); err != nil {
			return AnalysisReport{}, err
		}
	}
	sort.SliceStable(a.report.Flags, func(i, j int) bool { return a.report.Flags[i].Evidence.Time.Before(a.report.Flags[j].Evidence.Time) })
	sort.Strings(a.report.Warnings)
	return a.report, nil
}

func evidenceKey(e Evidence) string {
	return strings.Join([]string{e.Capture, e.SrcIP, e.SrcPort, e.DstIP, e.DstPort, e.Time.UTC().Format(time.RFC3339Nano)}, "|")
}
