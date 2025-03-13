package cron

import (
	"CBCTF/internal/db"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"CBCTF/internal/utils"
	"github.com/robfig/cron/v3"
	"time"
)

func CheckCheat(c *cron.Cron) {
	function := func() {
		contests, _, ok, _ := db.GetContests(db.DB, -1, -1, false, true, true)
		if !ok {
			return
		}
		for _, contest := range contests {
			if !contest.IsRunning() {
				continue
			}
			go checkFlag(contest)
		}
	}
	function()
	c.Schedule(cron.Every(30*time.Minute), cron.FuncJob(function))
}

func checkFlag(contest model.Contest) (bool, string) {
	var (
		submissions []model.Submission
		challenge   model.Challenge
		ok          bool
		msg         string
	)
	submissions, _, ok, msg = db.GetSubmissions(db.DB, -1, -1, "contest_id", contest.ID)
	if !ok {
		return false, msg
	}
	teamsSubmission := map[string][]model.Submission{}
	var teamsID []uint
	for _, submission := range submissions {
		if !utils.In(submission.TeamID, teamsID) {
			teamsSubmission[submission.Value] = append(teamsSubmission[submission.Value], submission)
			teamsID = append(teamsID, submission.TeamID)
		}
	}
	challenges := map[string]model.Challenge{}
	for _, flag := range teamsSubmission {
		if len(flag) < 2 {
			continue
		}
		challengeID := flag[0].ChallengeID
		challenge, ok = challenges[challengeID]
		if !ok {
			challenge, _, _ = db.GetChallengeByID(db.DB, challengeID)
			challenges[challengeID] = challenge
		}
		if challenge.Type != model.Static {
			var cheats []model.Cheat
			for _, submission := range flag {
				cheat := model.InitCheat(submission.UserID, submission.TeamID, submission.ContestID, model.SameFlag, model.Cheater)
				cheats = append(cheats, cheat)
			}
			tx := db.DB.Begin()
			for _, cheat := range cheats {
				for _, c := range cheats {
					if c.ID == cheat.ID {
						continue
					}
					cheat.Associations = append(cheat.Associations, c.ID)
				}
				if _, ok, msg = db.RecordCheat(tx, cheat); !ok {
					log.Logger.Warningf("Failed to record cheat %s in contest %d", cheat.ID, contest.ID)
				}
			}
			tx.Commit()
		}
	}
	return ok, "Success"
}
