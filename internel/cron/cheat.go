package cron

import (
	"CBCTF/internel/model"
	db "CBCTF/internel/repo"
	"CBCTF/internel/utils"
	"github.com/robfig/cron/v3"
	"time"
)

// CheckCheat 检查作弊事件
func CheckCheat(c *cron.Cron) {
	function := executionTime("CheckCheats", func() {
		repo := db.InitContestRepo(db.DB)
		contests, _, ok, _ := repo.GetAll(-1, -1, false)
		if !ok {
			return
		}
		for _, contest := range contests {
			if !contest.IsRunning() {
				continue
			}
			go checkFlag(contest)
		}
	})
	function()
	c.Schedule(cron.Every(30*time.Minute), cron.FuncJob(function))
}

// checkFlag 检查是否提交他人 flag
func checkFlag(contest model.Contest) (bool, string) {
	submissionRepo := db.InitSubmissionRepo(db.DB)
	cheatRepo := db.InitCheatRepo(db.DB)
	submissions, _, ok, msg := submissionRepo.GetAllByKeyID("contest_id", contest.ID, -1, -1, false)
	if !ok {
		return false, msg
	}
	flagSubmission := map[string][]model.Submission{}
	flagTeamID := map[string][]uint{}
	for _, submission := range submissions {
		if !utils.In(submission.TeamID, flagTeamID[submission.Value]) {
			flagSubmission[submission.Value] = append(flagSubmission[submission.Value], submission)
			flagTeamID[submission.Value] = append(flagTeamID[submission.Value], submission.TeamID)
		}
	}
	for _, flag := range flagSubmission {
		if len(flag) < 2 {
			continue
		}
		var cheats []db.CreateCheatOptions
		for _, submission := range flag {
			cheats = append(cheats, db.CreateCheatOptions{
				ID:        utils.UUID(),
				UserID:    submission.UserID,
				TeamID:    submission.TeamID,
				ContestID: contest.ID,
				Reason:    model.SameFlag,
				Type:      model.Suspect,
			})
		}
		for _, cheat := range cheats {
			others := make([]string, 0)
			for _, c := range cheats {
				if c.ID == cheat.ID {
					continue
				}
				others = append(others, c.ID)
			}
			cheat.Cheats = others
			cheatRepo.Create(cheat)
		}
	}
	return true, "Success"
}
