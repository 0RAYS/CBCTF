package cron

import (
	"context"
	"fmt"
	"math"
	"time"

	"CBCTF/internal/db"
	"CBCTF/internal/model"
	"CBCTF/internal/service"
)

// updateTeamRankingTask 全量更新 model.Team 的分数和排名
func updateTeamRankingTask() model.RetVal {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	root := db.CronDB.WithContext(ctx)
	job, ret := db.InitCronJobRepo(root).GetByUniqueField("name", model.UpdateTeamRankingCronJob)
	if !ret.OK {
		return ret
	}
	repo := db.InitContestRepo(root)
	contests, _, ret := repo.List(-1, -1, db.GetOptions{Conditions: map[string]any{"hidden": false}})
	if !ret.OK {
		return ret
	}
	batch := model.NewBatch(len(contests))
	for _, contest := range contests {
		if ctx.Err() != nil {
			return batch.Result(ctx)
		}
		key := fmt.Sprint(contest.ID)
		if !contest.EndedWithin(time.Now(), job.Schedule*2) {
			batch.Skip(key, "outside_refresh_window")
			continue
		}
		_, _, ret := service.UpdateTeamRanking(root, contest, -1, -1)
		if ret.OK {
			batch.Success(key, "published")
		} else {
			batch.Fail(key, "refresh_ranking", ret)
			if model.BatchDependencyFailed(ret) {
				return batch.Result(ctx)
			}
		}
	}
	return batch.Result(ctx)
}

// updateUserRankingTask 全量更新 model.User 的分数和排名
func updateUserRankingTask() model.RetVal {
	userRepo := db.InitUserRepo(db.CronDB)
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

	solvedFlags, ret := db.InitChallengeFlagRepo(db.CronDB).GetUserSolvedChallengeFlags(userIDs...)
	if !ret.OK {
		return ret
	}

	challengeFlagIDSet := make(map[uint]struct{})
	for _, flag := range solvedFlags {
		challengeFlagIDSet[flag.ID] = struct{}{}
	}

	challengeFlagIDL := make([]uint, 0, len(challengeFlagIDSet))
	for challengeFlagID := range challengeFlagIDSet {
		challengeFlagIDL = append(challengeFlagIDL, challengeFlagID)
	}

	bloodRankMap, ret := db.InitSubmissionRepo(db.CronDB).GetBloodRankMap(challengeFlagIDL...)
	if !ret.OK {
		return ret
	}

	userSolvedCount := make(map[uint]int64)
	userScore := make(map[uint]float64)
	for _, user := range users {
		userSolvedCount[user.ID] = 0
		userScore[user.ID] = 0
	}
	for _, flag := range solvedFlags {
		userSolvedCount[flag.UserID]++
		score := flag.CurrentScore
		switch bloodRankMap[flag.ID][flag.TeamID] {
		case 1:
			score += flag.Score * model.FirstBloodRate
		case 2:
			score += flag.Score * model.SecondBloodRate
		case 3:
			score += flag.Score * model.ThirdBloodRate
		}
		userScore[flag.UserID] += score
	}

	ret = db.WithTransactionDB(db.CronDB, func(tx *db.Tx) model.RetVal {
		for _, user := range users {
			if ret := db.InitUserRepo(tx).Update(user.ID, db.UpdateUserOptions{Score: new(math.Trunc(userScore[user.ID]*100) / 100), Solved: new(userSolvedCount[user.ID])}); !ret.OK {
				return ret
			}
		}
		return model.SuccessRetVal()
	})
	if !ret.OK {
		return ret
	}
	_, _, ret = service.UpdateUserRanking(db.CronDB, -1, -1)
	return ret
}
