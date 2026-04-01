package service

import (
	"CBCTF/internal/config"
	"CBCTF/internal/db"
	"CBCTF/internal/dto"
	"CBCTF/internal/model"
	"sort"
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
		ret = db.WithTransactionDB(tx, func(tx2 *gorm.DB) model.RetVal {
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
				return ret
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
					return ret
				}
			}
			contestChallenge, ret = contestChallengeRepo.GetByID(contestChallenge.ID, db.GetOptions{
				Preloads: map[string]db.GetOptions{"Challenge": {}, "ContestFlags": {}},
			})
			if !ret.OK {
				return ret
			}
			contestChallengeL = append(contestChallengeL, contestChallenge)
			return model.SuccessRetVal()
		})
		if !ret.OK {
			failedL = append(failedL, challengeRandID)
		}
	}
	return contestChallengeL, failedL, model.SuccessRetVal()
}

func GetContestChallengeImageList(tx *gorm.DB, contest model.Contest) ([]string, model.RetVal) {
	images, hasPodChallenges, ret := db.InitContestChallengeRepo(tx).ListContestImages(contest.ID)
	if !ret.OK {
		return nil, ret
	}
	imageSet := make(map[string]struct{}, len(images))
	for _, image := range images {
		imageSet[image] = struct{}{}
	}
	addImage := func(image string) {
		if image == "" {
			return
		}
		if _, ok := imageSet[image]; ok {
			return
		}
		imageSet[image] = struct{}{}
		images = append(images, image)
	}

	if hasPodChallenges {
		addImage(config.Env.K8S.TCPDumpImage)
		if config.Env.K8S.Frp.On {
			addImage(config.Env.K8S.Frp.FrpcImage)
			addImage(config.Env.K8S.Frp.NginxImage)
		}
	}
	sort.Strings(images)
	return images, model.SuccessRetVal()
}
