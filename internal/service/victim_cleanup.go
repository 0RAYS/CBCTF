package service

import (
	"context"
	"time"

	"gorm.io/gorm"

	"CBCTF/internal/db"
	"CBCTF/internal/log"
)

// Called after the flag transaction commits. Only enqueueing is synchronous;
// Kubernetes teardown runs in Asynq. Never hand a request-scoped transaction to
// an untracked goroutine, and finish this bounded side effect on disconnect.
func queueTeamVictimStop(tx *gorm.DB, teamID, challengeID uint) {
	ctx, cancel := context.WithTimeout(context.WithoutCancel(tx.Statement.Context), 10*time.Second)
	defer cancel()
	tx = tx.WithContext(ctx)
	victim, ret := db.InitVictimRepo(tx).HasAliveVictim(teamID, challengeID)
	if !ret.OK {
		return
	}
	if ret = ForceStopVictim(tx, victim); !ret.OK {
		log.Logger.Warningf("Failed to enqueue post-flag victim cleanup: victim_id=%d reason=%s", victim.ID, ret.Msg)
	}
}
