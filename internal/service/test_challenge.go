package service

import (
	"CBCTF/internal/db"
	"CBCTF/internal/dto"
	"CBCTF/internal/i18n"
	"CBCTF/internal/k8s"
	"CBCTF/internal/model"
	"context"
	"time"

	"gorm.io/gorm"
)

// GenTestAttachment 不使用任务队列生成附件，直接生成
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
		if ret.Msg != i18n.Model.Generator.NotAvailable {
			return ret
		}
		generators, ret := StartContestGenerators(tx, 0, dto.StartGeneratorsForm{Challenges: []string{challenge.RandID}})
		if !ret.OK {
			return ret
		}
		if len(generators) < 1 {
			return model.RetVal{Msg: i18n.Model.Generator.NotAvailable}
		}
		generator = generators[0]
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	ret = k8s.GenAttachment(ctx, challenge, generator, 0, flags)
	db.InitGeneratorRepo(tx).UpdateStatus(generator.ID, ret.OK, time.Now())
	StopContestGenerators(tx, dto.StopGeneratorsForm{Generators: []uint{generator.ID}})
	return ret
}
