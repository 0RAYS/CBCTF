package cron

import (
	"CBCTF/internal/db"
	"CBCTF/internal/model"
	"CBCTF/internal/service"
	"math"
	"time"
)

// updateTeamRankingTask 全量更新 model.Team 的分数和排名
func updateTeamRankingTask() model.RetVal {
	job, ret := db.InitCronJobRepo(db.DB).GetByUniqueField("name", model.UpdateTeamRankingCronJob)
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

	contestFlagIDSet := make(map[uint]struct{})
	for _, contestFlag := range solvedContestFlags {
		contestFlagIDSet[contestFlag.ID] = struct{}{}
	}

	contestFlagIDL := make([]uint, 0, len(contestFlagIDSet))
	for contestFlagID := range contestFlagIDSet {
		contestFlagIDL = append(contestFlagIDL, contestFlagID)
	}

	bloodRankMap, ret := db.InitSubmissionRepo(db.DB).GetBloodRankMap(contestFlagIDL...)
	if !ret.OK {
		return ret
	}

	userSolvedCount := make(map[uint]int64)
	userScore := make(map[uint]float64)
	for _, user := range users {
		userSolvedCount[user.ID] = 0
		userScore[user.ID] = 0
	}
	for _, contestFlag := range solvedContestFlags {
		userSolvedCount[contestFlag.UserID]++
		score := contestFlag.CurrentScore
		switch bloodRankMap[contestFlag.ID][contestFlag.TeamID] {
		case 1:
			score += contestFlag.Score * model.FirstBloodRate
		case 2:
			score += contestFlag.Score * model.SecondBloodRate
		case 3:
			score += contestFlag.Score * model.ThirdBloodRate
		}
		userScore[contestFlag.UserID] += score
	}

	for _, user := range users {
		solved := userSolvedCount[user.ID]
		score := math.Trunc(userScore[user.ID]*100) / 100
		userRepo.Update(user.ID, db.UpdateUserOptions{
			Score:  &score,
			Solved: &solved,
		})
	}
	service.UpdateUserRanking(db.DB, -1, -1)
	return model.SuccessRetVal()
}
