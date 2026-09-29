package traffic

import (
	"testing"
	"time"
)

func TestTrafficFlowAndDuplicateObservations(t *testing.T) {
	a := Connection{SrcIP: "8.8.8.8", DstIP: "10.0.0.2", SrcPort: "1000", DstPort: "80", Type: "TCP", Capture: "a", Fingerprint: "same", Time: time.Unix(1, 0)}
	b := a
	b.SrcIP, b.DstIP, b.SrcPort, b.DstPort = a.DstIP, a.SrcIP, a.DstPort, a.SrcPort
	if a.FlowKey() != b.FlowKey() {
		t.Fatal("reverse direction counted as another flow")
	}
	b = a
	b.Capture = "b"
	b.Time = b.Time.Add(time.Microsecond)
	c := a
	c.Time = c.Time.Add(2 * time.Millisecond)
	if got := deduplicateTraffic([]Connection{a, b, c}); len(got) != 2 {
		t.Fatalf("duplicate observations/retransmissions: %+v", got)
	}
}
