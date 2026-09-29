package service

import (
	"testing"
	"time"

	"CBCTF/internal/config"
	"CBCTF/internal/db"
)

func TestTrafficOverlapSpecialAddressFilter(t *testing.T) {
	previous := config.Env
	config.Env = &config.Config{}
	t.Cleanup(func() { config.Env = previous })
	at := time.Unix(1, 0)
	rows := []db.TeamVictimIP{
		{SrcIP: "8.8.8.8", TeamID: 1, FirstTime: at, VictimIDs: []uint{1, 2}},
		{SrcIP: "8.8.8.8", TeamID: 2, FirstTime: at, VictimIDs: []uint{3}},
		{SrcIP: "10.0.0.2", TeamID: 1}, {SrcIP: "10.0.0.2", TeamID: 2},
		{SrcIP: "fe80::1", TeamID: 1}, {SrcIP: "fe80::1", TeamID: 2},
	}
	got := groupTrafficOverlaps(rows)
	if len(got) != 1 || got[0].IP != "8.8.8.8" || len(got[0].Teams) != 2 || len(got[0].Teams[0].VictimIDs) != 2 {
		t.Fatalf("overlaps: %+v", got)
	}
}
