package service

import (
	f "CBCTF/internel/form"
	"CBCTF/internel/i18n"
	"CBCTF/internel/model"
	db "CBCTF/internel/repo"
	"errors"
	"gorm.io/gorm"
	"time"
)

func CreateContestChallenge(tx *gorm.DB, contest model.Contest, form f.CreateContestChallengeForm) ([]model.ContestChallenge, []string, bool, string) {
	contestChallengeL := make([]model.ContestChallenge, 0)
	failedL := make([]string, 0)
	contestChallengeRepo := db.InitContestChallengeRepo(tx)
	challengeRepo := db.InitChallengeRepo(tx)
	contestFlagRepo := db.InitContestFlagRepo(tx)
	for _, challengeRandID := range form.ChallengeRandIDL {
		challenge, ok, _ := challengeRepo.GetByRandID(challengeRandID, "ChallengeFlags")
		if !ok {
			failedL = append(failedL, challengeRandID)
			continue
		}
		if !contestChallengeRepo.IsUniqueContestChallenge(contest.ID, challenge.ID) {
			continue
		}
		_ = tx.Transaction(func(tx2 *gorm.DB) error {
			contestChallenge, ok, msg := contestChallengeRepo.Create(db.CreateContestChallengeOptions{
				ContestID:   contest.ID,
				ChallengeID: challenge.ID,
				Type:        challenge.Type,
				Name:        challenge.Name,
				Desc:        challenge.Desc,
				Hidden:      true,
			})
			if !ok {
				failedL = append(failedL, challengeRandID)
				return errors.New(msg)
			}
			for _, flag := range challenge.ChallengeFlags {
				contestFlagOptions := db.CreateContestFlagOptions{
					ContestID:          contest.ID,
					ContestChallengeID: contestChallenge.ID,
					ChallengeFlagID:    flag.ID,
					Value:              flag.Value,
					Score:              1000,
					CurrentScore:       1000,
					Decay:              100,
					MinScore:           100,
					ScoreType:          0,
					Solvers:            0,
					Last:               time.Now(),
				}
				_, ok, msg = contestFlagRepo.Create(contestFlagOptions)
				if !ok {
					failedL = append(failedL, challengeRandID)
					return errors.New(msg)
				}
			}
			contestChallenge, ok, msg = contestChallengeRepo.GetByID(contestChallenge.ID, "Challenge", "ContestFlags")
			if !ok {
				failedL = append(failedL, challengeRandID)
				return errors.New(msg)
			}
			contestChallengeL = append(contestChallengeL, contestChallenge)
			return nil
		})
	}
	return contestChallengeL, failedL, true, i18n.Success
}
