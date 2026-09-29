package traffic

import (
	"slices"
	"time"
)

// TrafficAccess is observed client evidence, never an arbitrary remote peer.
type TrafficAccess struct {
	IP      string    `json:"ip"`
	Time    time.Time `json:"time"`
	Source  string    `json:"source"`
	Capture string    `json:"capture"`
}

func CollectTrafficAccesses(result PcapDirResult, internalIPs map[string]bool) []TrafficAccess {
	accesses := append([]TrafficAccess{}, result.Accesses...)
	for _, conn := range result.Connections {
		if conn.ClientIP != "" {
			// Only the dedicated FRP capture is a trusted proxy evidence source.
			continue
		}
		if internalIPs[conn.DstIP] && !internalIPs[conn.SrcIP] && conn.Type == "TCP" && conn.SYN && !conn.ACK {
			accesses = append(accesses, TrafficAccess{IP: conn.SrcIP, Time: conn.Time, Source: "tcp_syn", Capture: conn.Capture})
		}
	}
	seen := make(map[string]bool)
	filtered := make([]TrafficAccess, 0, len(accesses))
	slices.SortStableFunc(accesses, func(a, b TrafficAccess) int { return a.Time.Compare(b.Time) })
	for _, access := range accesses {
		access.IP = NormalizeTrafficIP(access.IP)
		if !IsPublicTrafficIP(access.IP) || internalIPs[access.IP] {
			continue
		}
		// Keep one observation per IP/second/source, preserving contest-time filtering.
		key := access.IP + "|" + access.Time.UTC().Truncate(time.Second).Format(time.RFC3339) + "|" + access.Source
		if seen[key] {
			continue
		}
		seen[key] = true
		filtered = append(filtered, access)
	}
	return filtered
}
