package traffic

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestAttackIndicatorsAndScanWindow(t *testing.T) {
	a := newAnalyzer(context.Background(), nil)
	a.internalIPs = map[string]bool{"10.0.0.2": true}
	a.scanEncoded([]byte("GET /?q=UNION%20SELECT%20password&file=../../etc/passwd HTTP/1.1"), Evidence{}, "raw")
	for port := 1; port <= 20; port++ {
		a.observe(Connection{SrcIP: "8.8.8.8", DstIP: "10.0.0.2", SrcPort: "1000", DstPort: fmt.Sprint(port), Type: "TCP", SYN: true, Size: 60, Time: time.Unix(60, 0)})
	}
	a.finishActivity()
	found := make(map[string]bool)
	for _, indicator := range a.report.Indicators {
		found[indicator.Rule] = true
	}
	for _, rule := range []string{"sql_injection", "path_traversal", "port_scan"} {
		if !found[rule] {
			t.Errorf("missing %s: %+v", rule, a.report.Indicators)
		}
	}
	if len(a.report.Sessions) != 20 {
		t.Fatalf("sessions: %+v", a.report.Sessions)
	}
	b := newAnalyzer(context.Background(), nil)
	b.internalIPs = a.internalIPs
	for port := 1; port <= 20; port++ {
		b.observe(Connection{SrcIP: "8.8.8.8", DstIP: "10.0.0.2", DstPort: fmt.Sprint(port), SYN: true, Time: time.Unix(int64(port)*60, 0)})
	}
	for _, indicator := range b.report.Indicators {
		if indicator.Rule == "port_scan" {
			t.Fatal("unrelated minutes merged into a port scan")
		}
	}
}
