package service

import (
	"context"
	"time"

	"gorm.io/gorm"

	"CBCTF/internal/db"
	"CBCTF/internal/log"
)

// Called after the flag transaction commits. Only enqueueing is synchronous;
// Kubernetes teardown runs in Asynq, so the request-scoped transaction is never
// handed to an untracked goroutine and this side effect survives a disconnect.
func queueTeamVictimStop(tx *gorm.DB, teamID, contestChallengeID uint) {
	ctx, cancel := context.WithTimeout(context.WithoutCancel(tx.Statement.Context), 10*time.Second)
	defer cancel()
	tx = tx.WithContext(ctx)
	victim, ret := db.InitVictimRepo(tx).HasAliveVictimByContestChallenge(teamID, contestChallengeID)
	if !ret.OK {
		return
	}
	if ret = ForceStopVictim(tx, victim); !ret.OK {
		log.Logger.Warningf("Failed to enqueue post-flag victim cleanup: victim_id=%d reason=%v", victim.ID, ret)
	}
}
