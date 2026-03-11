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
	"errors"
	"fmt"
	"math/big"
	"time"

	"gorm.io/gorm"
)

func StartContestGenerators(tx *gorm.DB, contest model.Contest, form dto.StartGeneratorsForm) model.RetVal {
	if len(form.Challenges) == 0 {
		return model.SuccessRetVal()
	}
	challenges, _, ret := db.InitChallengeRepo(tx).List(-1, -1, db.GetOptions{
		Conditions: map[string]any{"type": model.DynamicChallengeType, "rand_id": form.Challenges},
	})
	if !ret.OK {
		return ret
	}
	for _, challenge := range challenges {
		_, ret = db.InitContestChallengeRepo(tx).Get(db.GetOptions{
			Conditions: map[string]any{"contest_id": contest.ID, "challenge_id": challenge.ID},
		})
		if !ret.OK {
			continue
		}
		_ = tx.Transaction(func(tx2 *gorm.DB) error {
			generator, ret := db.InitGeneratorRepo(tx2).Create(db.CreateGeneratorOptions{
				ChallengeID: challenge.ID,
				ContestID:   contest.ID,
				Name:        fmt.Sprintf("gen-%d-%d-%s", contest.ID, challenge.ID, utils.RandStr(6)),
			})
			if !ret.OK {
				return errors.New(ret.Msg)
			}
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
			_, ret = k8s.StartGenerator(ctx, challenge, generator)
			cancel()
			if !ret.OK {
				k8s.StopGenerator(ctx, generator)
				return errors.New(ret.Msg)
			}
			return nil
		})
	}
	return model.SuccessRetVal()
}

func StopContestGenerators(tx *gorm.DB, form dto.StopGeneratorsForm) model.RetVal {
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
