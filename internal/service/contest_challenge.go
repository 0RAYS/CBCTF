package service

import (
	"CBCTF/internal/config"
	"CBCTF/internal/db"
	"CBCTF/internal/dto"
	"CBCTF/internal/model"
	"CBCTF/internal/view"
	"os"
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

func buildContestChallengeFileName(tx *gorm.DB, challenge model.Challenge, teamID uint) string {
	path := challenge.AttachmentPath(teamID)
	record, _ := db.InitFileRepo(tx).Get(db.GetOptions{
		Conditions: map[string]any{
			"model":    model.ModelName(challenge),
			"model_id": challenge.ID,
			"type":     model.ChallengeFileType,
		},
	})
	filename := model.AttachmentFileName
	if string(record.Path) == path && record.Filename != "" {
		filename = record.Filename
	}
	if _, err := os.Stat(path); err != nil {
		return ""
	}
	return filename
}

func BuildContestChallengeRuntimeView(tx *gorm.DB, team model.Team, contestChallenge model.ContestChallenge) view.ContestChallengeView {
	return view.ContestChallengeView{
		ContestChallenge: contestChallenge,
		Attempts:         CountAttempts(tx, team, contestChallenge),
		Init:             CheckIfGenerated(tx, team, contestChallenge.ContestFlags),
		Solved:           CheckIfSolved(tx, team, contestChallenge.ContestFlags),
		Remote:           GetVictimStatus(tx, team.ID, contestChallenge.Challenge),
		FileName: func() string {
			if _, err := os.Stat(contestChallenge.Challenge.AttachmentPath(team.ID)); err != nil {
				return ""
			}
			return contestChallenge.Challenge.AttachmentPath(team.ID)
		}(),
	}
}

func ListContestChallengeViews(tx *gorm.DB, contest model.Contest, team model.Team, form dto.GetContestChallengesForm) ([]view.ContestChallengeView, int64, model.RetVal) {
	var (
		contestChallenges []model.ContestChallenge
		count             int64
		ret               model.RetVal
	)
	repo := db.InitContestChallengeRepo(tx)
	if form.Unsolved {
		ids, unsolvedCount, listRet := repo.ListUnsolvedID(team.ID, contest.ID, form.Category, form.Limit, form.Offset)
		if !listRet.OK {
			return nil, 0, listRet
		}
		count = unsolvedCount
		contestChallenges, _, ret = repo.List(-1, -1, db.GetOptions{
			Conditions: map[string]any{"id": ids},
			Preloads:   map[string]db.GetOptions{"Challenge": {}, "ContestFlags": {}},
		})
	} else {
		options := db.GetOptions{
			Conditions: map[string]any{"contest_id": contest.ID, "hidden": false},
			Preloads:   map[string]db.GetOptions{"Challenge": {}, "ContestFlags": {}},
		}
		if form.Category != "" {
			options.Conditions["category"] = form.Category
		}
		contestChallenges, count, ret = repo.List(form.Limit, form.Offset, options)
	}
	if !ret.OK {
		return nil, 0, ret
	}
	views := make([]view.ContestChallengeView, 0, len(contestChallenges))
	for _, contestChallenge := range contestChallenges {
		views = append(views, BuildContestChallengeRuntimeView(tx, team, contestChallenge))
	}
	return views, count, model.SuccessRetVal()
}

func ListAdminContestChallenges(tx *gorm.DB, contest model.Contest, form dto.GetAllContestChallengesForm) ([]model.ContestChallenge, int64, model.RetVal) {
	options := db.GetOptions{
		Conditions: map[string]any{"contest_id": contest.ID},
		Preloads:   map[string]db.GetOptions{"Challenge": {}, "ContestFlags": {}},
		Search:     make(map[string]string),
	}
	if form.Name != "" {
		options.Search["name"] = form.Name
	}
	if form.Category != "" {
		options.Conditions["category"] = form.Category
	}
	if form.Type != "" {
		options.Conditions["type"] = form.Type
	}
	return db.InitContestChallengeRepo(tx).List(form.Limit, form.Offset, options)
}

func ListContestChallengeCategories(tx *gorm.DB, contest model.Contest, form dto.GetCategoriesForm) ([]string, model.RetVal) {
	return db.InitContestChallengeRepo(tx).ListCategories(contest.ID, form.Type)
}

func GetContestChallengeStatus(tx *gorm.DB, team model.Team, challenge model.Challenge, contestChallenge model.ContestChallenge) (view.ContestChallengeStatusView, model.RetVal) {
	contestFlags, _, ret := db.InitContestFlagRepo(tx).List(-1, -1, db.GetOptions{
		Conditions: map[string]any{"contest_challenge_id": contestChallenge.ID},
	})
	if !ret.OK {
		return view.ContestChallengeStatusView{}, ret
	}
	return view.ContestChallengeStatusView{
		Attempts: CountAttempts(tx, team, contestChallenge),
		Init:     CheckIfGenerated(tx, team, contestFlags),
		Solved:   CheckIfSolved(tx, team, contestFlags),
		Remote:   GetVictimStatus(tx, team.ID, challenge),
		FileName: buildContestChallengeFileName(tx, challenge, team.ID),
	}, model.SuccessRetVal()
}

func UpdateContestChallenge(tx *gorm.DB, contestChallenge model.ContestChallenge, form dto.UpdateContestChallengeForm) model.RetVal {
	return db.InitContestChallengeRepo(tx).Update(contestChallenge.ID, db.UpdateContestChallengeOptions{
		Name:        form.Name,
		Description: form.Description,
		Hidden:      form.Hidden,
		Attempt:     form.Attempt,
		Hints:       form.Hints,
		Tags:        form.Tags,
	})
}

func DeleteContestChallenge(tx *gorm.DB, contestChallenge model.ContestChallenge) model.RetVal {
	return db.WithTransactionDB(tx, func(tx2 *gorm.DB) model.RetVal {
		return db.InitContestChallengeRepo(tx2).Delete(contestChallenge.ID)
	})
}

func ListContestFlagSolvers(tx *gorm.DB, contestFlag model.ContestFlag) ([]view.ContestFlagSolverView, model.RetVal) {
	rows, ret := db.InitSubmissionRepo(tx).ListFlagSolvers(contestFlag.ID)
	if !ret.OK {
		return nil, ret
	}
	solvers := make([]view.ContestFlagSolverView, 0, len(rows))
	for _, row := range rows {
		solvers = append(solvers, view.ContestFlagSolverView{
			UserID:   row.UserID,
			UserName: row.UserName,
			TeamID:   row.TeamID,
			TeamName: row.TeamName,
			Score:    row.Score,
			SolvedAt: row.SolvedAt,
		})
	}
	return solvers, model.SuccessRetVal()
}
