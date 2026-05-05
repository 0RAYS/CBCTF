package service

import (
	"CBCTF/internal/db"
	"CBCTF/internal/k8s"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"CBCTF/internal/view"
	"context"
	"time"

	"gorm.io/gorm"
)

// GenTestAttachment 不使用任务队列生成附件, 直接生成
func GenTestAttachment(tx *gorm.DB, challenge model.Challenge) model.RetVal {
	challengeFlags, _, ret := db.InitChallengeFlagRepo(tx).List(-1, -1, db.GetOptions{
		Conditions: map[string]any{"challenge_id": challenge.ID},
	})
	if !ret.OK {
		return ret
	}
	var flags []string
	for _, flag := range challengeFlags {
		flags = append(flags, flag.Value)
	}
	generator, ret := GetGenerator(tx, 0, challenge)
	if !ret.OK {
		return ret
	}
	log.Logger.Infof("Generating test attachment: challenge_id=%d generator_id=%d", challenge.ID, generator.ID)
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	ret = k8s.GenAttachment(ctx, challenge, generator, 0, flags)
	cancel()
	db.InitGeneratorRepo(tx).UpdateStatus(generator.ID, ret.OK, time.Now())
	if ret.OK {
		log.Logger.Infof("Test attachment generated: challenge_id=%d generator_id=%d", challenge.ID, generator.ID)
	}
	return ret
}

func GetTestChallengeStatus(tx *gorm.DB, challenge model.Challenge) view.ContestChallengeStatusView {
	return view.ContestChallengeStatusView{
		Attempts: 0,
		Init:     true,
		Solved:   false,
		Remote:   GetVictimStatus(tx, 0, challenge),
		FileName: buildContestChallengeFileName(tx, challenge, 0),
	}
}

func StopTestVictim(tx *gorm.DB, challenge model.Challenge) model.RetVal {
	victim, ret := db.InitVictimRepo(tx).HasAliveVictim(0, challenge.ID)
	if !ret.OK {
		return ret
	}
	return StopVictim(tx, victim)
}
