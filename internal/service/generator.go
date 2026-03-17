package service

import (
	"CBCTF/internal/db"
	"CBCTF/internal/dto"
	"CBCTF/internal/i18n"
	"CBCTF/internal/model"
	"CBCTF/internal/task"
	"CBCTF/internal/utils"
	"crypto/rand"
	"database/sql"
	"fmt"
	"math/big"

	"gorm.io/gorm"
)

func StartGenerators(tx *gorm.DB, contestID uint, form dto.StartGeneratorsForm) model.RetVal {
	if len(form.Challenges) == 0 {
		return model.SuccessRetVal()
	}
	challengeCount := make(map[string]int)
	for _, challenge := range form.Challenges {
		challengeCount[challenge] = challengeCount[challenge] + 1
	}
	challenges, _, ret := db.InitChallengeRepo(tx).List(-1, -1, db.GetOptions{
		Conditions: map[string]any{"type": model.DynamicChallengeType, "rand_id": form.Challenges},
	})
	if !ret.OK {
		return ret
	}
	contestChallengeRepo := db.InitContestChallengeRepo(tx)
	generatorRepo := db.InitGeneratorRepo(db.DB)
	for _, challenge := range challenges {
		if contestID > 0 {
			_, ret = contestChallengeRepo.Get(db.GetOptions{
				Conditions: map[string]any{"contest_id": contestID, "challenge_id": challenge.ID},
			})
			if !ret.OK {
				continue
			}
		}
		for range challengeCount[challenge.RandID] {
			generator, ret := generatorRepo.Create(db.CreateGeneratorOptions{
				ChallengeID:   challenge.ID,
				ChallengeName: challenge.Name,
				ContestID:     sql.Null[uint]{V: contestID, Valid: contestID > 0},
				Name:          fmt.Sprintf("gen-%d-%d-%s", contestID, challenge.ID, utils.RandStr(6)),
			})
			if !ret.OK {
				continue
			}
			_, _ = task.EnqueueStartGeneratorTask(challenge, generator)
		}
	}
	return model.SuccessRetVal()
}

func StopGenerators(tx *gorm.DB, form dto.StopGeneratorsForm) model.RetVal {
	if len(form.Generators) == 0 {
		return model.SuccessRetVal()
	}
	generators, _, ret := db.InitGeneratorRepo(tx).List(-1, -1, db.GetOptions{
		Conditions: map[string]any{"id": form.Generators},
	})
	if !ret.OK {
		return ret
	}
	for _, generator := range generators {
		_, _ = task.EnqueueStopGeneratorTask(generator)
	}
	return model.SuccessRetVal()
}

func GetGenerator(tx *gorm.DB, contestID uint, challenge model.Challenge) (model.Generator, model.RetVal) {
	options := db.GetOptions{Conditions: map[string]any{"challenge_id": challenge.ID}}
	if contestID > 0 {
		options.Conditions["contest_id"] = contestID
	}
	generators, _, ret := db.InitGeneratorRepo(tx).List(-1, -1, options)
	if !ret.OK {
		return model.Generator{}, ret
	}
	if len(generators) == 0 {
		return model.Generator{}, model.RetVal{Msg: i18n.Model.Generator.NotAvailable}
	}
	index, _ := rand.Int(rand.Reader, big.NewInt(int64(len(generators))))
	return generators[index.Int64()], model.SuccessRetVal()
}
