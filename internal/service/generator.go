package service

import (
	"CBCTF/internal/db"
	"CBCTF/internal/dto"
	"CBCTF/internal/i18n"
	"CBCTF/internal/k8s"
	"CBCTF/internal/model"
	"CBCTF/internal/utils"
	"context"
	"crypto/rand"
	"database/sql"
	"fmt"
	"math/big"
	"time"

	"gorm.io/gorm"
)

func StartGenerators(tx *gorm.DB, contestID uint, form dto.StartGeneratorsForm) model.RetVal {
	if len(form.Challenges) == 0 {
		return model.SuccessRetVal()
	}
	challenges, _, ret := db.InitChallengeRepo(tx).List(-1, -1, db.GetOptions{
		Conditions: map[string]any{"type": model.DynamicChallengeType, "rand_id": form.Challenges},
	})
	if !ret.OK {
		return ret
	}
	contestChallengeRepo := db.InitContestChallengeRepo(tx)
	for _, challenge := range challenges {
		if contestID > 0 {
			_, ret = contestChallengeRepo.Get(db.GetOptions{
				Conditions: map[string]any{"contest_id": contestID, "challenge_id": challenge.ID},
			})
			if !ret.OK {
				continue
			}
		}
		go func(contestID uint, challenge model.Challenge) {
			generatorRepo := db.InitGeneratorRepo(tx)
			generator, ret := generatorRepo.Create(db.CreateGeneratorOptions{
				ChallengeID: challenge.ID,
				ContestID:   sql.Null[uint]{V: contestID, Valid: contestID > 0},
				Name:        fmt.Sprintf("gen-%d-%d-%s", contestID, challenge.ID, utils.RandStr(6)),
			})
			if !ret.OK {
				return
			}
			ret = generatorRepo.Update(generator.ID, db.UpdateGeneratorOptions{Status: new(model.PendingGeneratorStatus)})
			if !ret.OK {
				return
			}
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
			_, ret = k8s.StartGenerator(ctx, challenge, generator)
			cancel()
			if !ret.OK {
				StopGenerators(tx, dto.StopGeneratorsForm{Generators: []uint{generator.ID}})
				return
			}
			ret = generatorRepo.Update(generator.ID, db.UpdateGeneratorOptions{Status: new(model.RunningGeneratorStatus)})
			if !ret.OK {
				return
			}
		}(contestID, challenge)
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
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
		ret = k8s.StopGenerator(ctx, generator)
		cancel()
		if !ret.OK {
			return ret
		}
		db.InitGeneratorRepo(tx).Delete(generator.ID)
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
