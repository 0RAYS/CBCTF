package traffic

import (
	"crypto/sha256"
	"fmt"
	"net"
	"time"
)

// FlowKey identifies a bidirectional transport conversation. Packet count and
// flow count are intentionally separate; reused tuples are one conversation.
func (c Connection) FlowKey() string {
	a, b := net.JoinHostPort(c.SrcIP, c.SrcPort), net.JoinHostPort(c.DstIP, c.DstPort)
	if a > b {
		a, b = b, a
	}
	return c.Type + "|" + a + "|" + b
}

func packetFingerprint(data []byte) string { return fmt.Sprintf("%x", sha256.Sum256(data)) }

// Remove byte-identical network packets observed by different capture pods
// within 1ms. Retransmissions in the same capture are retained.
func deduplicateTraffic(connections []Connection) []Connection {
	seen := make(map[string]Connection)
	result := make([]Connection, 0, len(connections))
	for _, c := range connections {
		previous, ok := seen[c.Fingerprint]
		if c.Fingerprint != "" && ok && c.Capture != previous.Capture && c.Time.Sub(previous.Time) <= time.Millisecond {
			continue
		}
		seen[c.Fingerprint] = c
		result = append(result, c)
	}
	return result
}
