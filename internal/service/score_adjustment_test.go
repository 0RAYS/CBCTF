package service

import (
	"testing"
)

func TestFinalizeTeamScoresAddsManualAdjustmentsToTheTotalOnly(t *testing.T) {
	scoreMap := map[uint]TeamScore{
		1: {Total: 1200, ByChallenge: map[uint]float64{10: 1200}},
		2: {Total: 300.555, ByChallenge: map[uint]float64{10: 300.555}},
	}
	finalizeTeamScores(scoreMap, map[uint]float64{1: 150, 2: -50})

	if scoreMap[1].Total != 1350 {
		t.Fatalf("team 1 total = %v, want 1350", scoreMap[1].Total)
	}
	if scoreMap[1].ByChallenge[10] != 1200 {
		t.Fatalf("team 1 challenge score = %v, want 1200", scoreMap[1].ByChallenge[10])
	}
	if scoreMap[2].Total != 250.55 {
		t.Fatalf("team 2 total = %v, want 250.55", scoreMap[2].Total)
	}
	if scoreMap[2].ByChallenge[10] != 300.55 {
		t.Fatalf("team 2 challenge score = %v, want 300.55", scoreMap[2].ByChallenge[10])
	}
}

func TestFinalizeTeamScoresKeepsSolvedScoreWhenNoAdjustmentExists(t *testing.T) {
	scoreMap := map[uint]TeamScore{1: {Total: 800, ByChallenge: map[uint]float64{10: 800}}}
	finalizeTeamScores(scoreMap, map[uint]float64{})

	if scoreMap[1].Total != 800 {
		t.Fatalf("team 1 total = %v, want 800", scoreMap[1].Total)
	}
}
