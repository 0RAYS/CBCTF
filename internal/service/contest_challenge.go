package service

import (
	"CBCTF/internal/config"
	"CBCTF/internal/db"
	"CBCTF/internal/dto"
	"CBCTF/internal/i18n"
	"CBCTF/internal/model"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

func CreateContestChallenge(tx *gorm.DB, contest model.Contest, form dto.CreateContestChallengeForm) ([]model.ContestChallenge, []string, model.RetVal) {
	contestChallengeL := make([]model.ContestChallenge, 0)
	failedL := make([]string, 0)
	contestChallengeRepo := db.InitContestChallengeRepo(tx)
	challengeRepo := db.InitChallengeRepo(tx)
	for _, challengeRandID := range form.ChallengeIDs {
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
		if err := tx.Transaction(func(tx2 *gorm.DB) error {
			contestChallengeRepo := db.InitContestChallengeRepo(tx2)
			contestFlagRepo := db.InitContestFlagRepo(tx2)

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
			if !ret.OK {
				if err, ok := ret.Attr["Error"]; ok {
					return errors.New(err.(string))
				}
				return fmt.Errorf("%s", ret.Msg)
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
				if !ret.OK {
					if err, ok := ret.Attr["Error"]; ok {
						return errors.New(err.(string))
					}
					return fmt.Errorf("%s", ret.Msg)
				}
			}
			contestChallenge, ret = contestChallengeRepo.GetByID(contestChallenge.ID, db.GetOptions{
				Preloads: map[string]db.GetOptions{"Challenge": {}, "ContestFlags": {}},
			})
			if !ret.OK {
				if err, ok := ret.Attr["Error"]; ok {
					return errors.New(err.(string))
				}
				return fmt.Errorf("%s", ret.Msg)
			}
			contestChallengeL = append(contestChallengeL, contestChallenge)
			return nil
		}); err != nil {
			failedL = append(failedL, challengeRandID)
		}
	}
	return contestChallengeL, failedL, model.SuccessRetVal()
}

func GetContestChallengeImageList(tx *gorm.DB, contest model.Contest) ([]string, model.RetVal) {
	imageSet := make(map[string]struct{})
	images := make([]string, 0)

	var generatorImages []string
	if res := tx.Table("contest_challenges").
		Distinct().
		Select("challenges.generator_image").
		Joins("INNER JOIN challenges ON contest_challenges.challenge_id = challenges.id AND challenges.deleted_at IS NULL").
		Where("contest_challenges.contest_id = ? AND contest_challenges.type = ? AND contest_challenges.deleted_at IS NULL", contest.ID, model.DynamicChallengeType).
		Where("challenges.generator_image <> ''").
		Order("challenges.generator_image ASC").
		Pluck("challenges.generator_image", &generatorImages); res.Error != nil {
		return nil, model.RetVal{Msg: i18n.Model.ContestChallenge.GetError, Attr: map[string]any{"Error": res.Error.Error()}}
	}
	for _, image := range generatorImages {
		if _, ok := imageSet[image]; ok {
			continue
		}
		imageSet[image] = struct{}{}
		images = append(images, image)
	}

	var dockerImages []string
	if res := tx.Table("contest_challenges").
		Distinct().
		Select("dockers.image").
		Joins("INNER JOIN challenges ON contest_challenges.challenge_id = challenges.id AND challenges.deleted_at IS NULL").
		Joins("INNER JOIN dockers ON dockers.challenge_id = challenges.id AND dockers.deleted_at IS NULL").
		Where("contest_challenges.contest_id = ? AND contest_challenges.type = ? AND contest_challenges.deleted_at IS NULL", contest.ID, model.PodsChallengeType).
		Where("dockers.image <> ''").
		Order("dockers.image ASC").
		Pluck("dockers.image", &dockerImages); res.Error != nil {
		return nil, model.RetVal{Msg: i18n.Model.ContestChallenge.GetError, Attr: map[string]any{"Error": res.Error.Error()}}
	}
	for _, image := range dockerImages {
		if _, ok := imageSet[image]; ok {
			continue
		}
		imageSet[image] = struct{}{}
		images = append(images, image)
	}

	for _, image := range []string{config.Env.K8S.Frp.FrpcImage, config.Env.K8S.Frp.NginxImage, config.Env.K8S.TCPDumpImage} {
		if image == "" {
			continue
		}
		if _, ok := imageSet[image]; ok {
			continue
		}
		imageSet[image] = struct{}{}
		images = append(images, image)
	}
	return images, model.SuccessRetVal()
}
