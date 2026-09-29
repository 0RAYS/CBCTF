package service

import (
	"math"
	"testing"
	"time"

	"CBCTF/internal/model"
	"CBCTF/internal/traffic"
)

func TestTrafficReplayWindows(t *testing.T) {
	connections := []traffic.Connection{{TimeShift: 0}, {TimeShift: 999 * time.Microsecond}, {TimeShift: time.Millisecond}, {TimeShift: 2 * time.Millisecond}}
	if got := sliceTrafficConnections(connections, 0, 1); len(got) != 2 {
		t.Fatalf("window boundary duplicated or omitted packets: %+v", got)
	}
	if got := sliceTrafficConnections(connections, 1, 2); len(got) != 1 {
		t.Fatalf("next window: %+v", got)
	}
	if got := sliceTrafficConnections(connections, 3, 3); len(got) != 0 {
		t.Fatal("empty window must stay empty")
	}
	start, end := clampTrafficWindow(1, math.MaxInt64, 10)
	if start != 1 || end != 10 {
		t.Fatalf("overflow: %d %d", start, end)
	}
}

func TestTrafficDirectionAndKnownIPs(t *testing.T) {
	if trafficDirection(false, false) != "external" || trafficDirection(true, true) != "internal" {
		t.Fatal("external peers classified as internal")
	}
	victim := model.Victim{Endpoints: model.Endpoints{{IP: "::ffff:10.0.0.2"}}, ExposedEndpoints: model.Endpoints{{IP: "8.8.8.8"}}}
	ips := victim.TrafficInternalIPs()
	if !ips["10.0.0.2"] || ips["8.8.8.8"] {
		t.Fatalf("incorrect target IP attribution: %+v", ips)
	}
}
