package cron

import (
	"CBCTF/internal/db"
	"CBCTF/internal/model"
	"CBCTF/internal/service"
	"math"
	"slices"
	"time"
)

// updateTeamRankingTask 全量更新 model.Team 的分数和排名
func updateTeamRankingTask() model.RetVal {
	job, ret := db.InitCronJobRepo(db.DB).GetByUniqueKey("name", model.UpdateTeamRankingCronJob)
	if !ret.OK {
		return ret
	}
	repo := db.InitContestRepo(db.DB)
	contests, _, ret := repo.List(-1, -1, db.GetOptions{Conditions: map[string]any{"hidden": false}})
	if !ret.OK {
		return ret
	}
	for _, contest := range contests {
		if time.Now().Sub(contest.Start.Add(contest.Duration)) > job.Schedule*2 {
			continue
		}
		service.UpdateTeamRanking(db.DB, contest, -1, -1)
	}
	return model.SuccessRetVal()
}

// updateUserRankingTask 全量更新 model.User 的分数和排名
func updateUserRankingTask() model.RetVal {
	userRepo := db.InitUserRepo(db.DB)
	users, _, ret := userRepo.List(-1, -1, db.GetOptions{
		Conditions: map[string]any{"banned": false},
	})
	if !ret.OK {
		return ret
	}
	userIDs := make([]uint, len(users))
	for i, user := range users {
		userIDs[i] = user.ID
	}

	solvedContestFlags, ret := db.InitContestFlagRepo(db.DB).GetUserSolvedContestFlags(userIDs...)
	if !ret.OK {
		return ret
	}

	userSolvedContestFlagsMap := make(map[uint][]db.UserSolvedContestFlag)
	contestIDSet := make(map[uint]bool)
	for _, contestFlag := range solvedContestFlags {
		userSolvedContestFlagsMap[contestFlag.UserID] = append(userSolvedContestFlagsMap[contestFlag.UserID], contestFlag)
		contestIDSet[contestFlag.ContestID] = true
	}

	contestIDL := make([]uint, 0, len(contestIDSet))
	for contestID := range contestIDSet {
		contestIDL = append(contestIDL, contestID)
	}
	contests, _, ret := db.InitContestRepo(db.DB).List(-1, -1, db.GetOptions{
		Conditions: map[string]any{"id": contestIDL},
	})
	if !ret.OK {
		return ret
	}
	blood := make(map[uint]bool)
	for _, contest := range contests {
		blood[contest.ID] = contest.Blood
	}
	submissionRepo := db.InitSubmissionRepo(db.DB)
	for _, user := range users {
		userSolvedContestFlags := userSolvedContestFlagsMap[user.ID]
		var solved int64 = 0
		var score float64 = 0
		for _, contestFlag := range userSolvedContestFlags {
			solved++
			var rate float64
			if blood[contestFlag.ContestID] {
				bloodTeam, _ := submissionRepo.GetBloodTeamID(contestFlag.ID)
				switch slices.IndexFunc(bloodTeam, func(i uint) bool {
					return i == contestFlag.TeamID
				}) {
				case 0:
					rate = model.FirstBloodRate
				case 1:
					rate = model.SecondBloodRate
				case 2:
					rate = model.ThirdBloodRate
				}
			}
			score += contestFlag.CurrentScore + contestFlag.Score*rate
		}
		score = math.Trunc(score*100) / 100
		userRepo.Update(user.ID, db.UpdateUserOptions{
			Score:  &score,
			Solved: &solved,
		})
	}
	service.UpdateUserRanking(db.DB, -1, -1)
	return model.SuccessRetVal()
}
