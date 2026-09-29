package traffic

import (
	"fmt"
	"regexp"
	"sort"
	"time"
)

type SessionSummary struct {
	ID        string   `json:"id"`
	Evidence  Evidence `json:"evidence"`
	Direction string   `json:"direction"`
	Packets   int64    `json:"packets"`
	Bytes     int64    `json:"bytes"`
	SYN       int64    `json:"syn"`
	RST       int64    `json:"rst"`
}

type AttackIndicator struct {
	Rule     string   `json:"rule"`
	Level    string   `json:"level"`
	Detail   string   `json:"detail"`
	Evidence Evidence `json:"evidence"`
}

type scanWindow struct {
	ports    map[string]bool
	evidence Evidence
}

var attackPatterns = []struct {
	rule    string
	pattern *regexp.Regexp
}{
	{"sql_injection", regexp.MustCompile(`(?i)\bunion\s+(?:all\s+)?select\b|\b(?:sleep|benchmark)\s*\(|['"]\s*or\s+['"]?1['"]?\s*=\s*['"]?1`)},
	{"path_traversal", regexp.MustCompile(`(?:\.\.[/\\]){2,}|/etc/(?:passwd|shadow)|/proc/self/(?:environ|maps)`)},
	{"command_execution", regexp.MustCompile(`(?i)/bin/(?:ba)?sh\b|\b(?:curl|wget)\s+https?://|\b(?:system|passthru|shell_exec)\s*\(|\$\([^)]{1,120}\)`)},
	{"template_injection", regexp.MustCompile(`(?i)\{\{[^\r\n]{0,120}(?:__class__|__globals__|config|7\s*\*\s*7)[^\r\n]{0,120}\}\}`)},
}

func (a *analyzer) indicator(rule, level, detail string, e Evidence) {
	key := rule + "|" + evidenceKey(e)
	if a.seenIndicators[key] {
		return
	}
	if len(a.report.Indicators) >= maxFindings {
		a.report.Truncated = true
		a.warn("indicator_limit")
		return
	}
	a.seenIndicators[key] = true
	if len(detail) > 256 {
		detail = detail[:256]
	}
	a.report.Indicators = append(a.report.Indicators, AttackIndicator{Rule: rule, Level: level, Detail: detail, Evidence: e})
}

func (a *analyzer) scanAttackContent(data []byte, e Evidence) {
	for _, rule := range attackPatterns {
		if matched := rule.pattern.Find(data); len(matched) > 0 {
			a.indicator(rule.rule, "suspicious", string(matched), e)
		}
	}
}

func (a *analyzer) observe(c Connection) {
	e := packetEvidence(c)
	key := c.Capture + "|" + c.FlowKey()
	session := a.sessions[key]
	if session == nil {
		if len(a.sessions) >= maxFindings {
			a.report.Truncated = true
			a.warn("session_limit")
			return
		}
		direction := "external"
		src, dst := a.internalIPs[c.SrcIP], a.internalIPs[c.DstIP]
		switch {
		case src && dst:
			direction = "internal"
		case src:
			direction = "egress"
		case dst:
			direction = "ingress"
		}
		session = &SessionSummary{ID: key, Evidence: e, Direction: direction}
		a.sessions[key] = session
	}
	session.Packets++
	session.Bytes += int64(c.Size)
	if c.Time.Before(session.Evidence.Time) {
		session.Evidence.Time = c.Time
	}
	if c.Time.After(session.Evidence.EndTime) {
		session.Evidence.EndTime = c.Time
	}
	if c.SYN && !c.ACK {
		session.SYN++
	}
	if c.RST {
		session.RST++
	}
	if a.internalIPs[c.SrcIP] && !a.internalIPs[c.DstIP] {
		a.outbound[key] += int64(c.Size)
	}
	if c.SYN && !c.ACK && a.internalIPs[c.DstIP] {
		key := c.Capture + "|" + c.SrcIP + ">" + c.DstIP + "|" + c.Time.UTC().Truncate(time.Minute).Format(time.RFC3339)
		window := a.scans[key]
		if window == nil {
			if len(a.scans) >= maxStreams {
				a.report.Truncated = true
				a.warn("scan_window_limit")
				return
			}
			window = &scanWindow{ports: make(map[string]bool), evidence: e}
			a.scans[key] = window
		}
		if len(window.ports) < 20 {
			window.ports[c.DstPort] = true
		}
		if len(window.ports) == 20 {
			a.indicator("port_scan", "suspicious", "at least 20 distinct destination ports in one UTC minute", window.evidence)
		}
	}
}

func (a *analyzer) finishActivity() {
	for key, session := range a.sessions {
		a.report.Sessions = append(a.report.Sessions, *session)
		if a.outbound[key] >= 1<<20 {
			a.indicator("large_outbound_transfer", "info", fmt.Sprintf("%d captured bytes leaving target network", a.outbound[key]), session.Evidence)
		}
		if session.SYN >= 10 {
			a.indicator("repeated_connection_attempts", "suspicious", fmt.Sprintf("%d SYN packets (may include retransmissions)", session.SYN), session.Evidence)
		}
		if session.RST >= 5 {
			a.indicator("connection_resets", "info", fmt.Sprintf("%d TCP resets", session.RST), session.Evidence)
		}
	}
	sort.Slice(a.report.Sessions, func(i, j int) bool {
		if a.report.Sessions[i].Bytes != a.report.Sessions[j].Bytes {
			return a.report.Sessions[i].Bytes > a.report.Sessions[j].Bytes
		}
		return a.report.Sessions[i].ID < a.report.Sessions[j].ID
	})
	sort.Slice(a.report.Indicators, func(i, j int) bool {
		x, y := a.report.Indicators[i], a.report.Indicators[j]
		if !x.Evidence.Time.Equal(y.Evidence.Time) {
			return x.Evidence.Time.Before(y.Evidence.Time)
		}
		return x.Rule+"|"+evidenceKey(x.Evidence) < y.Rule+"|"+evidenceKey(y.Evidence)
	})
}
