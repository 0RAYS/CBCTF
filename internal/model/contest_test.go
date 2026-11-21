package model

import (
	"testing"
	"time"
)

func TestContestStatusWhilePaused(t *testing.T) {
	contest := Contest{Start: time.Now().Add(-time.Hour), Duration: 2 * time.Hour}
	if got := contest.Status(); got != ContestIsRunning {
		t.Fatalf("status = %s, want %s", got, ContestIsRunning)
	}

	contest.Paused = true
	contest.PausedAt = time.Now()
	if got := contest.Status(); got != ContestIsPaused {
		t.Fatalf("status = %s, want %s", got, ContestIsPaused)
	}
	if contest.IsRunning() {
		t.Fatal("a paused contest must not report running")
	}
}

func TestContestEndKeepsRemainingTimeWhilePaused(t *testing.T) {
	contest := Contest{Start: time.Now().Add(-time.Hour), Duration: 2 * time.Hour}
	contest.Paused = true
	contest.PausedAt = time.Now()

	remaining := contest.End().Sub(time.Now())
	if remaining < 59*time.Minute || remaining > 61*time.Minute {
		t.Fatalf("remaining = %s, want about one hour", remaining)
	}
}
