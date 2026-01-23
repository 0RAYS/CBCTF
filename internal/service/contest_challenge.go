package service

import (
	"CBCTF/internal/config"
	"CBCTF/internal/db"
	"CBCTF/internal/dto"
	"CBCTF/internal/model"
	"errors"
	"slices"
	"time"

	"gorm.io/gorm"
)

func CreateContestChallenge(tx *gorm.DB, contest model.Contest, form dto.CreateContestChallengeForm) ([]model.ContestChallenge, []string, model.RetVal) {
	contestChallengeL := make([]model.ContestChallenge, 0)
	failedL := make([]string, 0)
	contestChallengeRepo := db.InitContestChallengeRepo(tx)
	challengeRepo := db.InitChallengeRepo(tx)
	contestFlagRepo := db.InitContestFlagRepo(tx)
	for _, challengeRandID := range form.ChallengeRandIDL {
		challenge, ret := challengeRepo.GetByRandID(challengeRandID, db.GetOptions{
			Preloads: map[string]db.GetOptions{"ChallengeFlags": {}},
		})
		if !ret.OK {
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
				Description: challenge.Description,
				Type:        challenge.Type,
				Category:    challenge.Category,
				Hidden:      true,
			}
			if challenge.Type == model.QuestionChallengeType {
				options.Attempt = 1
			}
			contestChallenge, ret := contestChallengeRepo.Create(options)
			if err, ok := ret.Attr["Error"]; ok && !ret.OK {
				failedL = append(failedL, challengeRandID)
				return errors.New(err.(string))
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
				_, ret = contestFlagRepo.Create(contestFlagOptions)
				if err, ok := ret.Attr["Error"]; ok && !ret.OK {
					failedL = append(failedL, challengeRandID)
					return errors.New(err.(string))
				}
			}
			contestChallenge, ret = contestChallengeRepo.GetByID(contestChallenge.ID, db.GetOptions{
				Preloads: map[string]db.GetOptions{"Challenge": {}, "ContestFlags": {}},
			})
			if err, ok := ret.Attr["Error"]; ok && !ret.OK {
				failedL = append(failedL, challengeRandID)
				return errors.New(err.(string))
			}
			contestChallengeL = append(contestChallengeL, contestChallenge)
			return nil
		})
	}
	return contestChallengeL, failedL, model.SuccessRetVal()
}

func GetContestChallengeImageList(tx *gorm.DB, contest model.Contest) ([]string, model.RetVal) {
	images := make([]string, 0)
	dynamicContestChallenges, _, ret := db.InitContestChallengeRepo(tx).List(-1, -1, db.GetOptions{
		Conditions: map[string]any{"type": model.DynamicChallengeType, "contest_id": contest.ID},
		Selects:    []string{"id", "challenge_id"},
		Preloads:   map[string]db.GetOptions{"Challenge": {Selects: []string{"id", "generator_image"}}},
	})
	if !ret.OK {
		return nil, ret
	}
	for _, contestChallenge := range dynamicContestChallenges {
		if !slices.Contains(images, contestChallenge.Challenge.GeneratorImage) {
			images = append(images, contestChallenge.Challenge.GeneratorImage)
		}
	}
	podsContestChallenge, _, ret := db.InitContestChallengeRepo(tx).List(-1, -1, db.GetOptions{
		Conditions: map[string]any{"type": model.PodsChallengeType, "contest_id": contest.ID},
		Selects:    []string{"id", "challenge_id"},
		Preloads: map[string]db.GetOptions{
			"Challenge": {
				Selects:  []string{"id"},
				Preloads: map[string]db.GetOptions{"Dockers": {Selects: []string{"id", "challenge_id", "image"}}},
			},
		},
	})
	if !ret.OK {
		return nil, ret
	}
	for _, contestChallenge := range podsContestChallenge {
		for _, docker := range contestChallenge.Challenge.Dockers {
			if !slices.Contains(images, docker.Image) {
				images = append(images, docker.Image)
			}
		}
	}
	if !slices.Contains(images, config.Env.K8S.Frpc.FrpcImage) {
		images = append(images, config.Env.K8S.Frpc.FrpcImage)
	}
	if !slices.Contains(images, config.Env.K8S.Frpc.NginxImage) {
		images = append(images, config.Env.K8S.Frpc.NginxImage)
	}
	if !slices.Contains(images, config.Env.K8S.TCPDumpImage) {
		images = append(images, config.Env.K8S.TCPDumpImage)
	}
	return images, model.SuccessRetVal()
}
