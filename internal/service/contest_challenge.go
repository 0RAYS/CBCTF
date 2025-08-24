package service

import (
	"CBCTF/internal/config"
	"CBCTF/internal/db"
	f "CBCTF/internal/form"
	"CBCTF/internal/i18n"
	"CBCTF/internal/model"
	"errors"
	"slices"
	"time"

	"gorm.io/gorm"
)

func CreateContestChallenge(tx *gorm.DB, contest model.Contest, form f.CreateContestChallengeForm) ([]model.ContestChallenge, []string, bool, string) {
	contestChallengeL := make([]model.ContestChallenge, 0)
	failedL := make([]string, 0)
	contestChallengeRepo := db.InitContestChallengeRepo(tx)
	challengeRepo := db.InitChallengeRepo(tx)
	contestFlagRepo := db.InitContestFlagRepo(tx)
	for _, challengeRandID := range form.ChallengeRandIDL {
		challenge, ok, _ := challengeRepo.GetByRandID(challengeRandID, db.GetOptions{
			Preloads: map[string]db.GetOptions{"ChallengeFlags": {}},
		})
		if !ok {
			failedL = append(failedL, challengeRandID)
			continue
		}
		if !contestChallengeRepo.IsUniqueContestChallenge(contest.ID, challenge.ID) {
			continue
		}
		_ = tx.Transaction(func(tx2 *gorm.DB) error {
			options := db.CreateContestChallengeOptions{
				ContestID:   contest.ID,
				ChallengeID: challenge.ID,
				Name:        challenge.Name,
				Desc:        challenge.Desc,
				Type:        challenge.Type,
				Category:    challenge.Category,
				Hidden:      true,
			}
			if challenge.Type == model.QuestionChallengeType {
				options.Attempt = 1
			}
			contestChallenge, ok, msg := contestChallengeRepo.Create(options)
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
			contestChallenge, ok, msg = contestChallengeRepo.GetByID(contestChallenge.ID, db.GetOptions{
				Preloads: map[string]db.GetOptions{"Challenge": {}, "ContestFlags": {}},
			})
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

func GetContestChallengeImageList(tx *gorm.DB, contest model.Contest) ([]string, bool, string) {
	images := make([]string, 0)
	dynamicContestChallenges, _, ok, msg := db.InitContestChallengeRepo(tx).List(-1, -1, db.GetOptions{
		Conditions: map[string]any{"type": model.DynamicChallengeType, "contest_id": contest.ID},
		Selects:    []string{"id", "challenge_id"},
		Preloads:   map[string]db.GetOptions{"Challenge": {Selects: []string{"id", "generator_image"}}},
	})
	if !ok {
		return nil, false, msg
	}
	for _, contestChallenge := range dynamicContestChallenges {
		if !slices.Contains(images, contestChallenge.Challenge.GeneratorImage) {
			images = append(images, contestChallenge.Challenge.GeneratorImage)
		}
	}
	podsContestChallenge, _, ok, msg := db.InitContestChallengeRepo(tx).List(-1, -1, db.GetOptions{
		Conditions: map[string]any{"type": model.PodsChallengeType, "contest_id": contest.ID},
		Selects:    []string{"id", "challenge_id"},
		Preloads: map[string]db.GetOptions{
			"Challenge": {
				Selects:  []string{"id"},
				Preloads: map[string]db.GetOptions{"Dockers": {Selects: []string{"id", "challenge_id", "image"}}},
			},
		},
	})
	if !ok {
		return nil, false, msg
	}
	for _, contestChallenge := range podsContestChallenge {
		for _, docker := range contestChallenge.Challenge.Dockers {
			if !slices.Contains(images, docker.Image) {
				images = append(images, docker.Image)
			}
		}
	}
	if !slices.Contains(images, config.Env.K8S.Frpc.Image) {
		images = append(images, config.Env.K8S.Frpc.Image)
	}
	if !slices.Contains(images, config.Env.K8S.TCPDumpImage) {
		images = append(images, config.Env.K8S.TCPDumpImage)
	}
	return images, true, msg
}
