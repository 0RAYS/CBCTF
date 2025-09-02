package cron

import (
	"CBCTF/internal/db"
	"CBCTF/internal/service"
	"sync"
	"time"

	"github.com/robfig/cron/v3"
)

// updateFlagScore 依据数据库, 更新 model.Flag 的分数和解题人数
// 正常情况下该定时任务无意义, 每次有新解出时即更新 current_score 和 solvers
// 当 submissions 且 model.Submission.Solved == true 时的数据减少 (例如: 用户注销 / 队伍解散 引发的数据删除), 该函数才有意义
func updateFlagScore(c *cron.Cron) {
	function := exec("UpdateFlagScore", func() {
		contestRepo := db.InitContestRepo(db.DB)
		contests, _, ok, _ := contestRepo.List(-1, -1)
		if !ok {
			return
		}
		for _, contest := range contests {
			if time.Now().Sub(contest.Start.Add(contest.Duration)) > 90*time.Minute {
				continue
			}
			contestChallengeRepo := db.InitContestChallengeRepo(db.DB)
			contestChallengeL, _, ok, _ := contestChallengeRepo.List(-1, -1, db.GetOptions{
				Conditions: map[string]any{"contest_id": contest.ID},
				Selects:    []string{"id"},
				Preloads: map[string]db.GetOptions{
					"ContestFlags": {
						Selects: []string{"id", "contest_challenge_id", "solvers", "score", "current_score", "decay", "min_score", "score_type"},
					},
				},
			})
			if !ok {
				return
			}
			for _, contestChallenge := range contestChallengeL {
				for _, contestFlag := range contestChallenge.ContestFlags {
					mu, _ := service.SolvedMutex.LoadOrStore(contestFlag.ID, &sync.Mutex{})
					mu.(*sync.Mutex).Lock()
					contestFlagRepo := db.InitContestFlagRepo(db.DB)
					solvers, currentScore, ok, _ := service.CalcContestFlagState(db.DB, contestFlag)
					if !ok {
						mu.(*sync.Mutex).Unlock()
						continue
					}
					if solvers != contestFlag.Solvers || currentScore != contestFlag.CurrentScore {
						if ok, _ = contestFlagRepo.Update(contestFlag.ID, db.UpdateContestFlagOptions{
							CurrentScore: &currentScore,
							Solvers:      &solvers,
						}); !ok {
							mu.(*sync.Mutex).Unlock()
							continue
						}
					}
					mu.(*sync.Mutex).Unlock()
				}
			}
		}
	})
	function()
	c.Schedule(cron.Every(time.Hour), cron.FuncJob(function))
}
