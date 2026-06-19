package service

import (
	"CBCTF/internal/db"
	"CBCTF/internal/i18n"
	"CBCTF/internal/k8s"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"CBCTF/internal/redis"
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
	generator, lockToken, ret := GetGenerator(tx, 0, challenge)
	if !ret.OK {
		return ret
	}
	defer func() {
		if lockToken == "" {
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := redis.UnlockGeneratorAttachment(ctx, generator.ID, lockToken); err != nil {
			log.Logger.Warningf("Failed to unlock test generator attachment: generator_id=%d error=%v", generator.ID, err)
		}
	}()
	log.Logger.Infof("Generating test attachment: challenge_id=%d generator_id=%d", challenge.ID, generator.ID)
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	ret = k8s.GenAttachment(ctx, challenge, generator, 0, flags)
	cancel()
	generatorRepo := db.InitGeneratorRepo(tx)
	generatorRepo.UpdateStatus(generator.ID, ret.OK, time.Now())
	if !ret.OK && (ret.Msg == i18n.Model.NotFound || ret.Msg == i18n.K8S.NotFound) {
		unregisterCtx, unregisterCancel := context.WithTimeout(context.Background(), 5*time.Second)
		if err := redis.UnregisterGenerator(unregisterCtx, generator); err != nil {
			log.Logger.Warningf("Failed to unregister test generator: generator_id=%d error=%v", generator.ID, err)
		}
		unregisterCancel()
		if deleteRet := generatorRepo.Delete(generator.ID); !deleteRet.OK {
			return deleteRet
		}
	}
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
